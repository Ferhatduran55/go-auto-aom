package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"go-auto-aom/storage"

	"github.com/jchv/go-webview2"
)

// ============================================
// Oto Parça Sipariş Sistemi - Masaüstü Uygulaması
// Bleve (gömülü Elasticsearch) + WebView2 GUI
// Tek EXE - Bağımlılık yok (Windows 10/11 WebView2 Runtime içerir)
// ============================================

//go:embed web/*
var webFS embed.FS

var store *storage.BleveStore

func main() {
	// Bleve store başlat
	var err error
	store, err = storage.NewBleveStore()
	if err != nil {
		fmt.Printf("Veritabanı hatası: %v\n", err)
		return
	}
	defer store.Close()

	// Boş port bul
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Printf("Port bulunamadı: %v\n", err)
		return
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	// HTTP sunucu başlat
	go startServer(port)
	time.Sleep(100 * time.Millisecond)

	// WebView2 oluştur
	w := webview2.NewWithOptions(webview2.WebViewOptions{
		Debug:     false,
		AutoFocus: true,
		WindowOptions: webview2.WindowOptions{
			Title:  "Oto Parça Sipariş Sistemi",
			Width:  1400,
			Height: 900,
			IconId: 1, // versioninfo.json'daki icon resource ID
			Center: true,
		},
	})
	if w == nil {
		// WebView2 yüklü değilse kullanıcıya bilgi ver
		showErrorDialog("WebView2 Runtime Gerekli",
			"Bu uygulamayı çalıştırmak için Microsoft WebView2 Runtime gereklidir.\n\n"+
				"Tamam'a tıklayarak indirme sayfasını açabilirsiniz.\n\n"+
				"İndirme sayfası: https://go.microsoft.com/fwlink/p/?LinkId=2124703")

		// WebView2 indirme sayfasını aç
		exec.Command("cmd", "/c", "start", "https://go.microsoft.com/fwlink/p/?LinkId=2124703").Start()
		return
	}
	defer w.Destroy()

	// JavaScript'ten Go fonksiyonlarını çağır
	w.Bind("saveOrderToBleve", saveOrderToBleve)
	w.Bind("loadOrdersFromBleve", loadOrdersFromBleve)
	w.Bind("loadOrderById", loadOrderById)
	w.Bind("deleteOrderFromBleve", deleteOrderFromBleve)
	w.Bind("searchOrders", searchOrders)
	w.Bind("searchOrdersAdvanced", searchOrdersAdvanced)
	w.Bind("searchCustomers", searchCustomers)
	w.Bind("getCustomerOrders", getCustomerOrders)
	w.Bind("listAllCustomers", listAllCustomers)
	w.Bind("searchProducts", searchProducts)
	w.Bind("listAllProducts", listAllProducts)
	w.Bind("saveProduct", saveProduct)
	w.Bind("updateProduct", updateProduct)
	w.Bind("deleteProduct", deleteProduct)
	w.Bind("updateCustomer", updateCustomer)
	w.Bind("deleteCustomer", deleteCustomer)

	w.Navigate(fmt.Sprintf("http://127.0.0.1:%d/", port))
	w.Run()
}

// HTTP sunucu - web dosyalarını sunar
func startServer(port int) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}

		content, err := webFS.ReadFile("web" + path)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		// Content type belirle
		switch {
		case strings.HasSuffix(path, ".html"):
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		case strings.HasSuffix(path, ".css"):
			w.Header().Set("Content-Type", "text/css")
		case strings.HasSuffix(path, ".js"):
			w.Header().Set("Content-Type", "application/javascript")
		case strings.HasSuffix(path, ".png"):
			w.Header().Set("Content-Type", "image/png")
		case strings.HasSuffix(path, ".ico"):
			w.Header().Set("Content-Type", "image/x-icon")
		}

		w.Write(content)
	})

	http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), nil)
}

// ============================================
// JavaScript'ten çağrılan Go fonksiyonları
// ============================================

// saveOrderToBleve - Siparişi Bleve'e kaydet (yeni veya güncelleme)
func saveOrderToBleve(orderJSON string) string {
	var orderData struct {
		ID            string              `json:"id"`
		Title         string              `json:"title"`
		CustomerID    string              `json:"customer_id"`
		CustomerName  string              `json:"customer_name"`
		CustomerPhone string              `json:"customer_phone"`
		Items         []storage.OrderItem `json:"items"`
	}

	if err := json.Unmarshal([]byte(orderJSON), &orderData); err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}

	// Müşteri varsa al veya oluştur (ID bazlı)
	var customerID string
	if orderData.CustomerName != "" {
		customer, err := store.GetOrCreateCustomer(orderData.CustomerName, orderData.CustomerPhone, orderData.CustomerID)
		if err == nil && customer != nil {
			customerID = customer.ID
		}
	}

	var order *storage.Order
	var err error

	// ID varsa güncelleme, yoksa yeni kayıt
	if orderData.ID != "" {
		order, err = store.GetOrder(orderData.ID)
		if err != nil {
			order = storage.NewOrder()
		}
	} else {
		order = storage.NewOrder()
	}

	order.Title = orderData.Title
	order.CustomerID = customerID
	order.CustomerName = orderData.CustomerName
	order.Items = orderData.Items
	order.CalculateGrandTotal()

	// ID varsa güncelle, yoksa yeni kaydet
	if orderData.ID != "" && err == nil {
		order.ID = orderData.ID
		if err := store.UpdateOrder(order); err != nil {
			return fmt.Sprintf(`{"error": "%s"}`, err.Error())
		}
	} else {
		if err := store.SaveOrder(order); err != nil {
			return fmt.Sprintf(`{"error": "%s"}`, err.Error())
		}
	}

	// Müşteri istatistiklerini güncelle
	if customerID != "" {
		store.UpdateCustomerStats(customerID)
	}

	return fmt.Sprintf(`{"success": true, "id": "%s", "customer_id": "%s"}`, order.ID, customerID)
}

// loadOrdersFromBleve - Siparişleri yükle (filtre ile)
func loadOrdersFromBleve(filterJSON string) string {
	var filter struct {
		Type      string `json:"type"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	// Eski format desteği (sadece string)
	if !strings.HasPrefix(filterJSON, "{") {
		filter.Type = filterJSON
	} else {
		json.Unmarshal([]byte(filterJSON), &filter)
	}

	var orders []*storage.Order
	var err error

	switch filter.Type {
	case "today":
		orders, err = store.ListTodayOrders()
	case "range":
		startDate, _ := time.Parse("2006-01-02", filter.StartDate)
		endDate, _ := time.Parse("2006-01-02", filter.EndDate)
		if startDate.IsZero() {
			startDate = time.Now().AddDate(0, -1, 0) // Son 1 ay
		}
		if endDate.IsZero() {
			endDate = time.Now()
		}
		orders, err = store.ListOrdersByDateRange(startDate, endDate)
	case "all":
		orders, err = store.ListOrders()
	default:
		orders, err = store.ListTodayOrders()
	}

	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}

	data, _ := json.Marshal(orders)
	return string(data)
}

// loadOrderById - ID'ye göre sipariş yükle
func loadOrderById(id string) string {
	order, err := store.GetOrder(id)
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}

	data, _ := json.Marshal(order)
	return string(data)
}

// deleteOrderFromBleve - Siparişi sil
func deleteOrderFromBleve(id string) string {
	if err := store.DeleteOrder(id); err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}

	return `{"success": true}`
}

// searchOrders - Siparişlerde arama yap
func searchOrders(searchTerm string) string {
	orders, err := store.SearchOrders(searchTerm)
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}

	data, _ := json.Marshal(orders)
	return string(data)
}

// searchOrdersAdvanced - Gelişmiş sipariş arama
func searchOrdersAdvanced(filterJSON string) string {
	var filter struct {
		ProductName  string  `json:"product_name"`
		OEMNumber    string  `json:"oem_number"`
		CustomerName string  `json:"customer_name"`
		MinQuantity  int     `json:"min_quantity"`
		MaxQuantity  int     `json:"max_quantity"`
		MinTotal     float64 `json:"min_total"`
		MaxTotal     float64 `json:"max_total"`
		MinUnitPrice float64 `json:"min_unit_price"`
		MaxUnitPrice float64 `json:"max_unit_price"`
		DateFilter   string  `json:"date_filter"`
		StartDate    string  `json:"start_date"`
		EndDate      string  `json:"end_date"`
	}

	if err := json.Unmarshal([]byte(filterJSON), &filter); err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}

	orders, err := store.SearchOrdersAdvanced(
		filter.ProductName,
		filter.OEMNumber,
		filter.CustomerName,
		filter.MinQuantity,
		filter.MaxQuantity,
		filter.MinTotal,
		filter.MaxTotal,
		filter.MinUnitPrice,
		filter.MaxUnitPrice,
		filter.DateFilter,
		filter.StartDate,
		filter.EndDate,
	)
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}

	data, _ := json.Marshal(orders)
	return string(data)
}

// searchCustomers - Müşteri ara (autocomplete için)
func searchCustomers(searchTerm string) string {
	customers, err := store.SearchCustomers(searchTerm)
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}

	data, _ := json.Marshal(customers)
	return string(data)
}

// listAllCustomers - Tüm müşterileri listele
func listAllCustomers() string {
	customers, err := store.ListCustomers()
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}

	data, _ := json.Marshal(customers)
	return string(data)
}

// getCustomerOrders - Müşterinin siparişlerini getir
func getCustomerOrders(customerID string) string {
	orders, err := store.GetCustomerOrders(customerID)
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}

	data, _ := json.Marshal(orders)
	return string(data)
}

// searchProducts - Ürün ara (autocomplete için)
func searchProducts(searchTerm string) string {
	products, err := store.SearchProducts(searchTerm)
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}

	data, _ := json.Marshal(products)
	return string(data)
}

// listAllProducts - Tüm ürünleri listele
func listAllProducts() string {
	products, err := store.ListProducts()
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}

	data, _ := json.Marshal(products)
	return string(data)
}

// saveProduct - Ürün kaydet veya güncelle
func saveProduct(productJSON string) string {
	var productData struct {
		Name      string `json:"name"`
		OEMNumber string `json:"oem_number"`
	}

	if err := json.Unmarshal([]byte(productJSON), &productData); err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}

	product, err := store.GetOrCreateProduct(productData.Name, productData.OEMNumber)
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}

	if product == nil {
		return `{"success": false}`
	}

	return fmt.Sprintf(`{"success": true, "id": "%s"}`, product.ID)
}

// updateProduct - Ürünü güncelle
func updateProduct(productJSON string) string {
	var productData struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		OEMNumber string `json:"oem_number"`
	}

	if err := json.Unmarshal([]byte(productJSON), &productData); err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}

	if err := store.UpdateProduct(productData.ID, productData.Name, productData.OEMNumber); err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}

	return `{"success": true}`
}

// deleteProduct - Ürünü sil
func deleteProduct(id string) string {
	if err := store.DeleteProduct(id); err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}

	return `{"success": true}`
}

// updateCustomer - Müşteriyi güncelle
func updateCustomer(customerJSON string) string {
	var customerData struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}

	if err := json.Unmarshal([]byte(customerJSON), &customerData); err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}

	if err := store.UpdateCustomer(customerData.ID, customerData.Name, customerData.Phone); err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}

	return `{"success": true}`
}

// deleteCustomer - Müşteriyi sil
func deleteCustomer(id string) string {
	if err := store.DeleteCustomer(id); err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}

	return `{"success": true}`
}

// showErrorDialog - Windows hata dialog'u gösterir
func showErrorDialog(title, message string) {
	user32 := syscall.NewLazyDLL("user32.dll")
	messageBox := user32.NewProc("MessageBoxW")

	titlePtr, _ := syscall.UTF16PtrFromString(title)
	messagePtr, _ := syscall.UTF16PtrFromString(message)

	// MB_OK | MB_ICONERROR = 0x10
	messageBox.Call(0, uintptr(unsafe.Pointer(messagePtr)), uintptr(unsafe.Pointer(titlePtr)), 0x10)
}
