package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"go-auto-aom/storage"

	"github.com/jchv/go-webview2"
)

// =============================================================================
// AutoManagement - Desktop Application
// Client-based relational management application with integrated
// order, stock, and reporting system.
//
// Tech Stack:
//   - Backend: Go + Bleve (embedded full-text search engine)
//   - Frontend: Vue.js 3 + TailwindCSS + Vite
//   - GUI: WebView2 (Windows native)
//   - Distribution: Single EXE, no external dependencies
//
// =============================================================================

// Application metadata - Keep in sync with versioninfo.json
var (
	AppName        = "AutoManagement"
	AppTitle       = "AutoManagement - Oto Yönetim Sistemi"
	AppVersion     = "26.2.1"
	BleveVersion   = ""
	WebViewVersion = ""
)

// Window configuration
const (
	WindowWidth  = 1400
	WindowHeight = 900
	WindowIconID = 1 // Resource ID from versioninfo.json
)

//go:embed web/*
var webFS embed.FS

var store *storage.BleveStore

func main() {
	// Initialize Bleve store
	var err error
	store, err = storage.NewBleveStore()
	if err != nil {
		LogError("System", fmt.Sprintf("Database error: %v", err))
		fmt.Printf("Database error: %v\n", err)
		return
	}
	defer store.Close()
	LogSystemWithCode(fmt.Sprintf("%s v%s starting...", AppName, AppVersion), "app.starting", map[string]interface{}{"version": AppVersion, "app": AppName})
	LogSystemWithCode("Database connected", "db.connected", nil)

	// Find available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		LogError("System", fmt.Sprintf("Port not found: %v", err))
		fmt.Printf("Failed to find available port: %v\n", err)
		return
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()
	LogSystemWithCode(fmt.Sprintf("HTTP server starting on port %d", port), "http.starting", map[string]interface{}{"port": port})

	// Start HTTP server
	go startServer(port)
	time.Sleep(100 * time.Millisecond)

	// Check if developer mode is enabled from settings
	debugMode := storage.IsDeveloperMode()
	if debugMode {
		LogSystemWithCode("Developer mode enabled", "dev.mode.enabled", nil)
	}

	// Create WebView2 window
	w := webview2.NewWithOptions(webview2.WebViewOptions{
		Debug:     debugMode,
		AutoFocus: true,
		WindowOptions: webview2.WindowOptions{
			Title:  AppTitle,
			Width:  WindowWidth,
			Height: WindowHeight,
			IconId: WindowIconID,
			Center: true,
		},
	})
	if w == nil {
		// WebView2 Runtime is not installed
		showErrorDialog("WebView2 Runtime Required",
			"Microsoft WebView2 Runtime is required to run this application.\n\n"+
				"Click OK to open the download page.\n\n"+
				"Download: https://go.microsoft.com/fwlink/p/?LinkId=2124703")

		// Open WebView2 download page
		exec.Command("cmd", "/c", "start", "https://go.microsoft.com/fwlink/p/?LinkId=2124703").Start()
		return
	}

	// Try to detect the WebView2 runtime/browser version at runtime.
	// If detection succeeds, override the build-embedded WebViewVersion value.
	if detected := detectWebViewRuntimeVersion(w); detected != "" {
		WebViewVersion = detected
		LogSystem(fmt.Sprintf("Detected WebView2 runtime version: %s", WebViewVersion))
	}

	defer w.Destroy()

	// Bind Go functions to JavaScript
	bindOrderFunctions(w)
	bindCustomerFunctions(w)
	bindProductFunctions(w)
	bindStockFunctions(w)
	bindSettingsFunctions(w)
	bindNoteFunctions(w)
	bindWhatsAppOrderFunctions(w)
	bindUpdateFunctions(w)
	bindLogFunctions(w)
	bindSystemFunctions(w)

	w.Navigate(fmt.Sprintf("http://127.0.0.1:%d/", port))
	w.Run()
}

// =============================================================================
// Function Bindings
// =============================================================================

// bindOrderFunctions binds order-related functions to WebView
func bindOrderFunctions(w webview2.WebView) {
	w.Bind("saveOrderToBleve", saveOrderToBleve)
	w.Bind("loadOrdersFromBleve", loadOrdersFromBleve)
	w.Bind("loadOrderById", loadOrderById)
	w.Bind("deleteOrderFromBleve", deleteOrderFromBleve)
	w.Bind("searchOrders", searchOrders)
	w.Bind("searchOrdersAdvanced", searchOrdersAdvanced)
}

// detectWebViewRuntimeVersion attempts to discover the installed WebView2 runtime/browser version
// by trying a few known methods on the WebView object via reflection. This is best-effort and will not
// panic on failures; if nothing is found, it returns an empty string.
func detectWebViewRuntimeVersion(w webview2.WebView) string {
	defer func() {
		if r := recover(); r != nil {
			// ignore panics from reflection/calls
		}
	}()

	v := reflect.ValueOf(w)
	if !v.IsValid() {
		return ""
	}

	// Common method names to try on the WebView itself
	methodNames := []string{"GetBrowserVersionString", "BrowserVersionString", "GetVersion"}
	for _, name := range methodNames {
		m := v.MethodByName(name)
		if m.IsValid() && m.Type().NumIn() == 0 {
			res := m.Call(nil)
			if len(res) >= 1 {
				if s, ok := res[0].Interface().(string); ok && s != "" {
					return s
				}
			}
		}
	}

	// Try to get a CoreWebView2-like object and call methods on it
	coreM := v.MethodByName("CoreWebView2")
	if coreM.IsValid() && coreM.Type().NumIn() == 0 {
		r := coreM.Call(nil)
		if len(r) >= 1 {
			core := r[0]
			for _, name := range []string{"GetBrowserVersionString", "BrowserVersionString", "GetVersion"} {
				mc := core.MethodByName(name)
				if mc.IsValid() && mc.Type().NumIn() == 0 {
					res := mc.Call(nil)
					if len(res) >= 1 {
						if s, ok := res[0].Interface().(string); ok && s != "" {
							return s
						}
					}
				}
			}
		}
	}

	return ""
}

// bindCustomerFunctions binds customer-related functions to WebView
func bindCustomerFunctions(w webview2.WebView) {
	w.Bind("searchCustomers", searchCustomers)
	w.Bind("getCustomerOrders", getCustomerOrders)
	w.Bind("listAllCustomers", listAllCustomers)
	w.Bind("updateCustomer", updateCustomer)
	w.Bind("deleteCustomer", deleteCustomer)
}

// bindProductFunctions binds product-related functions to WebView
func bindProductFunctions(w webview2.WebView) {
	w.Bind("searchProducts", searchProducts)
	w.Bind("listAllProducts", listAllProducts)
	w.Bind("listProductsPaginated", listProductsPaginated)
	w.Bind("saveProduct", saveProduct)
	w.Bind("updateProduct", updateProduct)
	w.Bind("deleteProduct", deleteProduct)
	w.Bind("createProductFull", createProductFull)
	w.Bind("getCategories", getCategories)
	w.Bind("getBrands", getBrands)
	w.Bind("getUnits", getUnits)
}

// bindStockFunctions binds stock management functions to WebView
func bindStockFunctions(w webview2.WebView) {
	w.Bind("stockIn", stockIn)
	w.Bind("stockOut", stockOut)
	w.Bind("bulkStockIn", bulkStockIn)
	w.Bind("bulkStockOut", bulkStockOut)
	w.Bind("getStockMovements", getStockMovements)
	w.Bind("getCriticalStockProducts", getCriticalStockProducts)
	w.Bind("getStockReport", getStockReport)
}

// bindSettingsFunctions binds settings functions to WebView
func bindSettingsFunctions(w webview2.WebView) {
	w.Bind("setDeveloperMode", setDeveloperMode)
	w.Bind("getDeveloperMode", getDeveloperMode)
}

// bindNoteFunctions binds note-related functions to WebView
func bindNoteFunctions(w webview2.WebView) {
	w.Bind("saveNote", saveNote)
	w.Bind("loadNotes", loadNotes)
	w.Bind("loadNoteById", loadNoteById)
	w.Bind("deleteNote", deleteNote)
	w.Bind("searchNotes", searchNotes)
	w.Bind("autoFormatText", autoFormatText)
}

// bindWhatsAppOrderFunctions binds WhatsApp order functions to WebView
func bindWhatsAppOrderFunctions(w webview2.WebView) {
	w.Bind("saveWhatsAppOrder", saveWhatsAppOrder)
	w.Bind("loadWhatsAppOrders", loadWhatsAppOrders)
	w.Bind("loadWhatsAppOrderById", loadWhatsAppOrderById)
	w.Bind("deleteWhatsAppOrder", deleteWhatsAppOrder)
	w.Bind("searchWhatsAppOrders", searchWhatsAppOrders)
	w.Bind("addWhatsAppOrderStatus", addWhatsAppOrderStatus)
	w.Bind("filterWhatsAppOrdersByStatus", filterWhatsAppOrdersByStatus)
}

// bindUpdateFunctions binds update-related functions to WebView
func bindUpdateFunctions(w webview2.WebView) {
	w.Bind("checkForUpdates", checkForUpdates)
	w.Bind("getAppVersion", getAppVersion)
}

// bindLogFunctions binds log-related functions to WebView
func bindLogFunctions(w webview2.WebView) {
	w.Bind("getBackendLogs", getLogsJSON)
	w.Bind("getBackendLogsAfterId", getLogsAfterIdJSON)
	w.Bind("clearBackendLogs", clearLogs)
	w.Bind("setLogBufferSize", setLogBufferSizeJSON)
	w.Bind("getLogBufferSize", getLogBufferSize)
}

// Uygulama başlangıç zamanı
var appStartTime = time.Now()

// bindSystemFunctions binds system-related functions to WebView
func bindSystemFunctions(w webview2.WebView) {
	w.Bind("getSystemInfo", getSystemInfo)
	w.Bind("exportAllData", exportAllData)
	w.Bind("copySystemInfo", copySystemInfo)
	w.Bind("clearCache", clearCache)
}

// getSystemInfo returns comprehensive system information
func getSystemInfo() string {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	// CPU model bilgisini al
	cpuModel := getCPUModel()

	info := map[string]interface{}{
		"app_name":        AppName,
		"app_version":     AppVersion,
		"go_version":      strings.TrimPrefix(runtime.Version(), "go"),
		"bleve_version":   BleveVersion,
		"webview_version": WebViewVersion,
		"os":              runtime.GOOS,
		"arch":            runtime.GOARCH,
		"cpu_model":       cpuModel,
		"uptime":          formatDuration(time.Since(appStartTime)),
		"uptime_seconds":  int64(time.Since(appStartTime).Seconds()),

	}

	data, _ := json.Marshal(info)
	return string(data)
}

// formatDuration formats duration as HH:MM:SS (zero-padded)
func formatDuration(d time.Duration) string {
	t := int64(d.Seconds())
	h := t / 3600
	m := (t % 3600) / 60
	s := t % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

// getCPUModel Windows'ta CPU model adını alır
func getCPUModel() string {
	// 1. Try PowerShell (most reliable on modern Windows)
	// -ExecutionPolicy Bypass is used to avoid script restriction issues
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-NoProfile", "-Command", "Get-CimInstance -ClassName Win32_Processor | Select-Object -ExpandProperty Name")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	output, err := cmd.Output()
	if err == nil {
		model := strings.TrimSpace(string(output))
		if model != "" {
			return model
		}
	}

	// 2. Fallback to wmic (older Windows)
	cmd = exec.Command("wmic", "cpu", "get", "name")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	output, err = cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && line != "Name" {
				return line
			}
		}
	}

	// 3. Last resort: Environment variable
	if val := os.Getenv("PROCESSOR_IDENTIFIER"); val != "" {
		return val
	}

	return "Bilinmiyor"
}

// exportAllData exports all data as JSON
func exportAllData() string {
	LogSystemWithCode("Exporting data...", "export.started", nil)

	// Get all data
	orders, _ := store.ListOrders()
	customers, _ := store.ListCustomers()
	products, _ := store.ListProducts()

	exportData := map[string]interface{}{
		"export_date": time.Now().Format(time.RFC3339),
		"app_version": fmt.Sprintf("v%s", AppVersion),
		"orders":      orders,
		"customers":   customers,
		"products":    products,
	}

	data, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		LogErrorWithCode("System", fmt.Sprintf("Export error: %v", err), "export.error", map[string]interface{}{"error": err.Error()})
		return fmt.Sprintf(`{"success": false, "error": "%s"}`, err.Error())
	}

	LogSystemWithCode(fmt.Sprintf("Export completed: %d orders, %d customers, %d products", len(orders), len(customers), len(products)), "export.completed", map[string]interface{}{"orders": len(orders), "customers": len(customers), "products": len(products)})
	return string(data)
}

// copySystemInfo returns system info in a copyable format for error reporting
func copySystemInfo() string {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	info := fmt.Sprintf(`=== AutoManagement System Info ===
Application: %s %s
Go: %s
Bleve: %s
WebView: %s
OS: %s/%s
CPU: %d cores
Goroutines: %d
Memory (Alloc): %d MB
Memory (Sys): %d MB
GC Cycles: %d
Uptime: %s
Date: %s
=====================================`,
		AppName, AppVersion,
		runtime.Version(),
		BleveVersion,
		WebViewVersion,
		runtime.GOOS, runtime.GOARCH,
		runtime.NumCPU(),
		runtime.NumGoroutine(),
		mem.Alloc/1024/1024,
		mem.Sys/1024/1024,
		mem.NumGC,
		formatDuration(time.Since(appStartTime)),
		time.Now().Format("2006-01-02 15:04:05"),
	)

	return info
}

// clearCache clears temporary files
func clearCache() string {
	LogSystemWithCode("Clearing cache...", "cache.clearing", nil)

	tempDir := filepath.Join(os.TempDir(), "AutoManagement")
	if err := os.RemoveAll(tempDir); err != nil {
		LogErrorWithCode("System", fmt.Sprintf("Cache clear error: %v", err), "cache.clear_error", map[string]interface{}{"error": err.Error()})
		return fmt.Sprintf(`{"success": false, "error": "%s"}`, err.Error())
	}

	LogSystemWithCode("Cache cleared", "cache.cleared", nil)
	return `{"success": true}`
}

// =============================================================================
// HTTP Server
// =============================================================================

// startServer starts the embedded HTTP server for serving web files
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

		// Set Content-Type header
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
		case strings.HasSuffix(path, ".svg"):
			w.Header().Set("Content-Type", "image/svg+xml")
		case strings.HasSuffix(path, ".json"):
			w.Header().Set("Content-Type", "application/json")
		case strings.HasSuffix(path, ".woff"), strings.HasSuffix(path, ".woff2"):
			w.Header().Set("Content-Type", "font/woff2")
		}

		w.Write(content)
	})

	http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), nil)
}

// =============================================================================
// Order Functions
// =============================================================================

// saveOrderToBleve saves or updates an order in the Bleve store
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
		return jsonError(err)
	}

	// Get or create customer if name is provided
	var customerID string
	if orderData.CustomerName != "" {
		customer, err := store.GetOrCreateCustomer(orderData.CustomerName, orderData.CustomerPhone, orderData.CustomerID)
		if err == nil && customer != nil {
			customerID = customer.ID
		}
	}

	var order *storage.Order
	var err error

	// Update existing or create new order
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

	// Save or update order
	if orderData.ID != "" && err == nil {
		order.ID = orderData.ID
		if err := store.UpdateOrder(order); err != nil {
			return jsonError(err)
		}
	} else {
		if err := store.SaveOrder(order); err != nil {
			return jsonError(err)
		}
	}

	// Update customer statistics
	if customerID != "" {
		store.UpdateCustomerStats(customerID)
	}

	return fmt.Sprintf(`{"success": true, "id": "%s", "customer_id": "%s"}`, order.ID, customerID)
}

// loadOrdersFromBleve loads orders with optional filtering
func loadOrdersFromBleve(filterJSON string) string {
	var filter struct {
		Type      string `json:"type"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	// Support legacy string format
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
			startDate = time.Now().AddDate(0, -1, 0) // Last month
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
		return jsonError(err)
	}

	return jsonMarshal(orders)
}

// loadOrderById loads a single order by ID
func loadOrderById(id string) string {
	order, err := store.GetOrder(id)
	if err != nil {
		return jsonError(err)
	}
	return jsonMarshal(order)
}

// deleteOrderFromBleve deletes an order by ID
func deleteOrderFromBleve(id string) string {
	if err := store.DeleteOrder(id); err != nil {
		return jsonError(err)
	}
	return jsonSuccess()
}

// searchOrders searches orders by term
func searchOrders(searchTerm string) string {
	orders, err := store.SearchOrders(searchTerm)
	if err != nil {
		return jsonError(err)
	}
	return jsonMarshal(orders)
}

// searchOrdersAdvanced performs advanced order search with multiple filters
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
		return jsonError(err)
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
		return jsonError(err)
	}

	return jsonMarshal(orders)
}

// =============================================================================
// Customer Functions
// =============================================================================

// searchCustomers searches customers for autocomplete
func searchCustomers(searchTerm string) string {
	customers, err := store.SearchCustomers(searchTerm)
	if err != nil {
		return jsonError(err)
	}
	return jsonMarshal(customers)
}

// listAllCustomers returns all customers
func listAllCustomers() string {
	customers, err := store.ListCustomers()
	if err != nil {
		return jsonError(err)
	}
	return jsonMarshal(customers)
}

// getCustomerOrders returns orders for a specific customer
func getCustomerOrders(customerID string) string {
	orders, err := store.GetCustomerOrders(customerID)
	if err != nil {
		return jsonError(err)
	}
	return jsonMarshal(orders)
}

// updateCustomer updates customer information
func updateCustomer(customerJSON string) string {
	var data struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}

	if err := json.Unmarshal([]byte(customerJSON), &data); err != nil {
		return jsonError(err)
	}

	if err := store.UpdateCustomer(data.ID, data.Name, data.Phone); err != nil {
		return jsonError(err)
	}

	return jsonSuccess()
}

// deleteCustomer removes a customer
func deleteCustomer(id string) string {
	if err := store.DeleteCustomer(id); err != nil {
		return jsonError(err)
	}
	return jsonSuccess()
}

// =============================================================================
// Product Functions
// =============================================================================

// searchProducts searches products for autocomplete
func searchProducts(searchTerm string) string {
	products, err := store.SearchProducts(searchTerm)
	if err != nil {
		return jsonError(err)
	}
	return jsonMarshal(products)
}

// listAllProducts returns all products
func listAllProducts() string {
	products, err := store.ListProducts()
	if err != nil {
		return jsonError(err)
	}
	return jsonMarshal(products)
}

// listProductsPaginated returns paginated and filtered product list
func listProductsPaginated(filterJSON string) string {
	var request struct {
		Page         int    `json:"page"`
		PageSize     int    `json:"page_size"`
		Search       string `json:"search"`
		Category     string `json:"category"`
		OnlyCritical bool   `json:"only_critical"`
		SortField    string `json:"sort_field"`
		SortDir      string `json:"sort_dir"`
	}

	if err := json.Unmarshal([]byte(filterJSON), &request); err != nil {
		return jsonError(err)
	}

	filter := storage.ProductFilter{
		Search:       request.Search,
		Category:     request.Category,
		OnlyCritical: request.OnlyCritical,
		SortField:    request.SortField,
		SortDir:      request.SortDir,
	}

	result, err := store.ListProductsPaginated(request.Page, request.PageSize, filter)
	if err != nil {
		return jsonError(err)
	}

	return jsonMarshal(result)
}

// saveProduct saves or updates a product
func saveProduct(productJSON string) string {
	var data struct {
		Name      string `json:"name"`
		OEMNumber string `json:"oem_number"`
	}

	if err := json.Unmarshal([]byte(productJSON), &data); err != nil {
		return jsonError(err)
	}

	product, err := store.GetOrCreateProduct(data.Name, data.OEMNumber)
	if err != nil {
		return jsonError(err)
	}

	if product == nil {
		return `{"success": false}`
	}

	return fmt.Sprintf(`{"success": true, "id": "%s"}`, product.ID)
}

// updateProduct updates product information
func updateProduct(productJSON string) string {
	var data struct {
		ID            string `json:"id"`
		Name          string `json:"name"`
		OEMNumber     string `json:"oem_number"`
		Brand         string `json:"brand"`
		Category      string `json:"category"`
		Unit          string `json:"unit"`
		CriticalStock int    `json:"critical_stock"`
	}

	if err := json.Unmarshal([]byte(productJSON), &data); err != nil {
		return jsonError(err)
	}

	if err := store.UpdateProduct(
		data.ID,
		data.Name,
		data.OEMNumber,
		data.Brand,
		data.Category,
		data.Unit,
		data.CriticalStock,
	); err != nil {
		return jsonError(err)
	}

	return jsonSuccess()
}

// deleteProduct removes a product
func deleteProduct(id string) string {
	if err := store.DeleteProduct(id); err != nil {
		return jsonError(err)
	}
	return jsonSuccess()
}

// createProductFull creates a product with all fields
func createProductFull(productJSON string) string {
	var data struct {
		Name          string  `json:"name"`
		OEMNumber     string  `json:"oem_number"`
		Brand         string  `json:"brand"`
		Category      string  `json:"category"`
		Unit          string  `json:"unit"`
		StockQuantity float64 `json:"stock_quantity"`
		CriticalStock int     `json:"critical_stock"`
	}

	if err := json.Unmarshal([]byte(productJSON), &data); err != nil {
		return jsonError(err)
	}

	product, err := store.CreateProductFull(
		data.Name,
		data.OEMNumber,
		data.Brand,
		data.Category,
		data.Unit,
		data.StockQuantity,
		data.CriticalStock,
	)

	if err != nil {
		return jsonError(err)
	}

	return jsonMarshal(product)
}

// getCategories returns all product categories
func getCategories() string {
	categories, err := store.GetCategories()
	if err != nil {
		return jsonError(err)
	}
	return jsonMarshal(categories)
}

// getBrands returns all product brands
func getBrands() string {
	brands, err := store.GetBrands()
	if err != nil {
		return jsonError(err)
	}
	return jsonMarshal(brands)
}

// getUnits returns available units
func getUnits() string {
	units := storage.GetUnits()
	return jsonMarshal(units)
}

// =============================================================================
// Stock Management Functions
// =============================================================================

// stockIn adds stock to a product
func stockIn(dataJSON string) string {
	var data struct {
		ProductID string  `json:"product_id"`
		Amount    float64 `json:"amount"`
		Note      string  `json:"note"`
	}

	if err := json.Unmarshal([]byte(dataJSON), &data); err != nil {
		return jsonError(err)
	}

	if err := store.StockIn(data.ProductID, data.Amount, data.Note); err != nil {
		return jsonError(err)
	}

	return jsonSuccess()
}

// stockOut removes stock from a product
func stockOut(dataJSON string) string {
	var data struct {
		ProductID string  `json:"product_id"`
		Amount    float64 `json:"amount"`
		Note      string  `json:"note"`
	}

	if err := json.Unmarshal([]byte(dataJSON), &data); err != nil {
		return jsonError(err)
	}

	if err := store.StockOut(data.ProductID, data.Amount, data.Note); err != nil {
		return jsonError(err)
	}

	return jsonSuccess()
}

// bulkStockIn adds stock to multiple products
func bulkStockIn(dataJSON string) string {
	var entries []storage.BulkStockInfo

	if err := json.Unmarshal([]byte(dataJSON), &entries); err != nil {
		return jsonError(err)
	}

	successful, errors := store.BulkStockIn(entries)

	result := map[string]interface{}{
		"success":    true,
		"successful": successful,
		"errors":     errors,
	}

	return jsonMarshal(result)
}

// bulkStockOut removes stock from multiple products
func bulkStockOut(dataJSON string) string {
	var entries []storage.BulkStockInfo

	if err := json.Unmarshal([]byte(dataJSON), &entries); err != nil {
		return jsonError(err)
	}

	successful, errors := store.BulkStockOut(entries)

	result := map[string]interface{}{
		"success":    true,
		"successful": successful,
		"errors":     errors,
	}

	return jsonMarshal(result)
}

// getStockMovements returns stock movement history
func getStockMovements(filterJSON string) string {
	var filter struct {
		ProductID string `json:"product_id"`
		StartStr  string `json:"start"`
		EndStr    string `json:"end"`
	}

	if err := json.Unmarshal([]byte(filterJSON), &filter); err != nil {
		return jsonError(err)
	}

	var start, end time.Time
	if filter.StartStr != "" {
		start, _ = time.Parse("2006-01-02", filter.StartStr)
	}
	if filter.EndStr != "" {
		end, _ = time.Parse("2006-01-02", filter.EndStr)
		end = end.Add(24 * time.Hour) // Include end of day
	}

	movements, err := store.GetStockMovements(filter.ProductID, start, end)
	if err != nil {
		return jsonError(err)
	}

	return jsonMarshal(movements)
}

// getCriticalStockProducts returns products below critical stock level
func getCriticalStockProducts() string {
	products, err := store.GetCriticalStockProducts()
	if err != nil {
		return jsonError(err)
	}
	return jsonMarshal(products)
}

// getStockReport generates a stock report for a given period
func getStockReport(filterJSON string) string {
	var filter struct {
		Period  string `json:"period"` // "daily" or "monthly"
		DateStr string `json:"date"`   // "2024-01-15" or "2024-01"
	}

	if err := json.Unmarshal([]byte(filterJSON), &filter); err != nil {
		return jsonError(err)
	}

	var date time.Time
	if filter.Period == "daily" {
		date, _ = time.Parse("2006-01-02", filter.DateStr)
	} else {
		date, _ = time.Parse("2006-01", filter.DateStr)
	}

	if date.IsZero() {
		date = time.Now()
	}

	report, err := store.GetStockReport(filter.Period, date)
	if err != nil {
		return jsonError(err)
	}

	return jsonMarshal(report)
}

// =============================================================================
// Helper Functions
// =============================================================================

// showErrorDialog displays a Windows error dialog
func showErrorDialog(title, message string) {
	user32 := syscall.NewLazyDLL("user32.dll")
	messageBox := user32.NewProc("MessageBoxW")

	titlePtr, _ := syscall.UTF16PtrFromString(title)
	messagePtr, _ := syscall.UTF16PtrFromString(message)

	// MB_OK | MB_ICONERROR = 0x10
	messageBox.Call(0, uintptr(unsafe.Pointer(messagePtr)), uintptr(unsafe.Pointer(titlePtr)), 0x10)
}

// jsonError returns a JSON error response
func jsonError(err error) string {
	return fmt.Sprintf(`{"error": "%s"}`, err.Error())
}

// jsonSuccess returns a JSON success response
func jsonSuccess() string {
	return `{"success": true}`
}

// jsonMarshal marshals data to JSON string
func jsonMarshal(data interface{}) string {
	bytes, _ := json.Marshal(data)
	return string(bytes)
}

// =============================================================================
// Settings Functions
// =============================================================================

// setDeveloperMode enables or disables developer mode
// Note: Changes take effect after application restart
func setDeveloperMode(enabled bool) string {
	if err := storage.UpdateDeveloperMode(enabled); err != nil {
		return jsonError(err)
	}
	return jsonMarshal(map[string]interface{}{
		"success":          true,
		"restart_required": true,
		"message":          "Developer mode changed. Restart the application for changes to take effect.",
	})
}

// getDeveloperMode returns current developer mode status
func getDeveloperMode() string {
	enabled := storage.IsDeveloperMode()
	return jsonMarshal(map[string]bool{"enabled": enabled})
}

// =============================================================================
// Note Functions
// =============================================================================

// saveNote saves or updates a note in the Bleve store
func saveNote(noteJSON string) string {
	var noteData struct {
		ID           string `json:"id"`
		Title        string `json:"title"`
		Content      string `json:"content"`
		CustomerName string `json:"customer_name"`
	}

	if err := json.Unmarshal([]byte(noteJSON), &noteData); err != nil {
		return jsonError(err)
	}

	note := &storage.Note{
		ID:           noteData.ID,
		Title:        noteData.Title,
		Content:      noteData.Content,
		CustomerName: noteData.CustomerName,
	}

	// Mevcut notu güncelle veya yeni oluştur
	if noteData.ID != "" {
		existingNote, err := store.GetNote(noteData.ID)
		if err == nil {
			note.CreatedAt = existingNote.CreatedAt
		}
	}

	if err := store.SaveNote(note); err != nil {
		return jsonError(err)
	}

	return fmt.Sprintf(`{"success": true, "id": "%s"}`, note.ID)
}

// loadNotes loads all notes, optionally filtered by search term
func loadNotes(searchTerm string) string {
	var notes []*storage.Note
	var err error

	if searchTerm != "" {
		notes, err = store.SearchNotes(searchTerm)
	} else {
		notes, err = store.ListNotes()
	}

	if err != nil {
		return jsonError(err)
	}

	return jsonMarshal(notes)
}

// loadNoteById loads a single note by ID
func loadNoteById(id string) string {
	note, err := store.GetNote(id)
	if err != nil {
		return jsonError(err)
	}
	return jsonMarshal(note)
}

// deleteNote deletes a note by ID
func deleteNote(id string) string {
	if err := store.DeleteNote(id); err != nil {
		return jsonError(err)
	}
	return jsonSuccess()
}

// searchNotes searches notes by term
func searchNotes(searchTerm string) string {
	notes, err := store.SearchNotes(searchTerm)
	if err != nil {
		return jsonError(err)
	}
	return jsonMarshal(notes)
}

// autoFormatText formats raw text with alignment
func autoFormatText(rawText string) string {
	formatted := storage.AutoFormatText(rawText)
	return jsonMarshal(map[string]string{"formatted": formatted})
}

// =============================================================================
// WhatsApp Order Functions
// =============================================================================

// saveWhatsAppOrder saves or updates a WhatsApp order
func saveWhatsAppOrder(orderJSON string) string {
	var orderData struct {
		ID            string                      `json:"id"`
		Date          string                      `json:"date"` // ISO string from frontend
		CustomerName  string                      `json:"customer_name"`
		CustomerPhone string                      `json:"customer_phone"`
		PaymentMethod string                      `json:"payment_method"`
		PaymentNote   string                      `json:"payment_note"`
		Items         []storage.WhatsAppOrderItem `json:"items"`
	}

	if err := json.Unmarshal([]byte(orderJSON), &orderData); err != nil {
		return jsonError(err)
	}

	// Parse date
	var orderDate time.Time
	if orderData.Date != "" {
		// Try parsing ISO format first, then date-only format
		var err error
		orderDate, err = time.Parse(time.RFC3339, orderData.Date)
		if err != nil {
			orderDate, err = time.Parse("2006-01-02", orderData.Date)
			if err != nil {
				orderDate = time.Now()
			}
		}
	} else {
		orderDate = time.Now()
	}

	order := &storage.WhatsAppOrder{
		ID:            orderData.ID,
		Date:          orderDate,
		CustomerName:  orderData.CustomerName,
		CustomerPhone: orderData.CustomerPhone,
		PaymentMethod: orderData.PaymentMethod,
		PaymentNote:   orderData.PaymentNote,
		Items:         orderData.Items,
	}

	// Mevcut siparişi güncelle veya yeni oluştur
	if orderData.ID != "" {
		existingOrder, err := store.GetWhatsAppOrder(orderData.ID)
		if err == nil {
			order.CreatedAt = existingOrder.CreatedAt
			order.StatusHistory = existingOrder.StatusHistory
			order.CurrentStatus = existingOrder.CurrentStatus
		}
	}

	if err := store.SaveWhatsAppOrder(order); err != nil {
		return jsonError(err)
	}

	return fmt.Sprintf(`{"success": true, "id": "%s"}`, order.ID)
}

// loadWhatsAppOrders loads all WhatsApp orders, optionally filtered by search term
func loadWhatsAppOrders(searchTerm string) string {
	var orders []*storage.WhatsAppOrder
	var err error

	if searchTerm != "" {
		orders, err = store.SearchWhatsAppOrders(searchTerm)
	} else {
		orders, err = store.ListWhatsAppOrders()
	}

	if err != nil {
		return jsonError(err)
	}

	return jsonMarshal(orders)
}

// loadWhatsAppOrderById loads a single WhatsApp order by ID
func loadWhatsAppOrderById(id string) string {
	order, err := store.GetWhatsAppOrder(id)
	if err != nil {
		return jsonError(err)
	}
	return jsonMarshal(order)
}

// deleteWhatsAppOrder deletes a WhatsApp order by ID
func deleteWhatsAppOrder(id string) string {
	if err := store.DeleteWhatsAppOrder(id); err != nil {
		return jsonError(err)
	}
	return jsonSuccess()
}

// searchWhatsAppOrders searches WhatsApp orders by term
func searchWhatsAppOrders(searchTerm string) string {
	orders, err := store.SearchWhatsAppOrders(searchTerm)
	if err != nil {
		return jsonError(err)
	}
	return jsonMarshal(orders)
}

// addWhatsAppOrderStatus adds a new status to a WhatsApp order
func addWhatsAppOrderStatus(statusJSON string) string {
	var data struct {
		OrderID string `json:"order_id"`
		Status  string `json:"status"`
		Note    string `json:"note"`
	}

	if err := json.Unmarshal([]byte(statusJSON), &data); err != nil {
		return jsonError(err)
	}

	if err := store.AddWhatsAppOrderStatus(data.OrderID, data.Status, data.Note); err != nil {
		return jsonError(err)
	}

	return jsonSuccess()
}

// filterWhatsAppOrdersByStatus filters WhatsApp orders by status
func filterWhatsAppOrdersByStatus(status string) string {
	orders, err := store.FilterWhatsAppOrdersByStatus(status)
	if err != nil {
		return jsonError(err)
	}
	return jsonMarshal(orders)
}
