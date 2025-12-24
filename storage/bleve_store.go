package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/google/uuid"
)

// ============================================
// Bleve - Gömülü Elasticsearch Alternatifi
// Full-text arama, indexing - sunucu gerektirmez
// ============================================

const (
	IndexName = "auto_parts_orders"
)

// Product - Ürün Kataloğu
type Product struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`       // Ürün adı (örn: "Megane 4 Joy Far")
	OEMNumber string    `json:"oem_number"` // OEM numarası
	Brand     string    `json:"brand"`      // Marka (örn: "Renault")
	Category  string    `json:"category"`   // Kategori (örn: "Aydınlatma")
	UsedCount int       `json:"used_count"` // Kaç kez kullanıldı
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Customer - Müşteri/Tedarikçi
type Customer struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Phone       string    `json:"phone"`
	Address     string    `json:"address"`
	Notes       string    `json:"notes"`
	OrderCount  int       `json:"order_count"`
	TotalAmount float64   `json:"total_amount"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// OrderItem - Sipariş içindeki ürün
type OrderItem struct {
	ID          string  `json:"id"`
	ProductName string  `json:"product_name"`
	OEMNumber   string  `json:"oem_number"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	PartStatus  string  `json:"part_status"` // "original" veya "used"
	TotalPrice  float64 `json:"total_price"`
}

// Order - Sipariş
type Order struct {
	ID           string      `json:"id"`
	Title        string      `json:"title"`         // Sipariş başlığı (örn: "Aralık İhale 1")
	CustomerID   string      `json:"customer_id"`   // Müşteri ID
	CustomerName string      `json:"customer_name"` // Müşteri/Tedarikçi adı (denormalize)
	Items        []OrderItem `json:"items"`
	GrandTotal   float64     `json:"grand_total"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

// BleveStore - Bleve tabanlı depolama
type BleveStore struct {
	index    bleve.Index
	dataPath string
}

// NewBleveStore - Yeni Bleve store oluşturur
func NewBleveStore() (*BleveStore, error) {
	// Kullanıcı veri dizinini al
	dataPath, err := getDataPath()
	if err != nil {
		return nil, fmt.Errorf("veri dizini alınamadı: %w", err)
	}

	indexPath := filepath.Join(dataPath, IndexName+".bleve")

	var index bleve.Index

	// Index varsa aç, yoksa oluştur
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		// Yeni index oluştur
		mapping := buildIndexMapping()
		index, err = bleve.New(indexPath, mapping)
		if err != nil {
			return nil, fmt.Errorf("index oluşturulamadı: %w", err)
		}
	} else {
		// Mevcut index'i aç
		index, err = bleve.Open(indexPath)
		if err != nil {
			return nil, fmt.Errorf("index açılamadı: %w", err)
		}
	}

	return &BleveStore{
		index:    index,
		dataPath: dataPath,
	}, nil
}

// buildIndexMapping - Elasticsearch benzeri index mapping
func buildIndexMapping() mapping.IndexMapping {
	// Ürün mapping
	itemMapping := bleve.NewDocumentMapping()
	itemMapping.AddFieldMappingsAt("product_name", bleve.NewTextFieldMapping())
	itemMapping.AddFieldMappingsAt("oem_number", bleve.NewKeywordFieldMapping())
	itemMapping.AddFieldMappingsAt("part_status", bleve.NewKeywordFieldMapping())

	// Sipariş mapping
	orderMapping := bleve.NewDocumentMapping()
	orderMapping.AddSubDocumentMapping("items", itemMapping)
	orderMapping.AddFieldMappingsAt("customer_name", bleve.NewTextFieldMapping())

	// Ana mapping
	indexMapping := bleve.NewIndexMapping()
	indexMapping.DefaultMapping = orderMapping
	indexMapping.DefaultAnalyzer = "standard"

	return indexMapping
}

// getDataPath - Uygulama veri dizinini döner
func getDataPath() (string, error) {
	// Windows'ta: %APPDATA%\OtoParçaSiparis
	// Linux/Mac'te: ~/.otoparça-siparis
	var dataDir string

	if os.Getenv("APPDATA") != "" {
		dataDir = filepath.Join(os.Getenv("APPDATA"), "OtoParcaSiparis")
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dataDir = filepath.Join(homeDir, ".otoparca-siparis")
	}

	// Dizin yoksa oluştur
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return "", err
	}

	return dataDir, nil
}

// Close - Index'i kapat
func (s *BleveStore) Close() error {
	return s.index.Close()
}

// SaveOrder - Siparişi kaydet (Elasticsearch'teki Index işlemi)
func (s *BleveStore) SaveOrder(order *Order) error {
	if order.ID == "" {
		order.ID = uuid.New().String()
	}

	order.UpdatedAt = time.Now()
	if order.CreatedAt.IsZero() {
		order.CreatedAt = time.Now()
	}

	// Toplamları hesapla
	order.CalculateGrandTotal()

	// JSON'a çevir ve kaydet
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("JSON dönüştürme hatası: %w", err)
	}

	// Bleve'e indexle
	if err := s.index.Index(order.ID, order); err != nil {
		return fmt.Errorf("indexleme hatası: %w", err)
	}

	// Ayrıca JSON dosyası olarak da sakla (veri kaybını önlemek için)
	orderFile := filepath.Join(s.dataPath, "orders", order.ID+".json")
	if err := os.MkdirAll(filepath.Dir(orderFile), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(orderFile, data, 0644); err != nil {
		return fmt.Errorf("dosya yazma hatası: %w", err)
	}

	return nil
}

// GetOrder - Siparişi getir
func (s *BleveStore) GetOrder(id string) (*Order, error) {
	orderFile := filepath.Join(s.dataPath, "orders", id+".json")
	data, err := os.ReadFile(orderFile)
	if err != nil {
		return nil, fmt.Errorf("sipariş bulunamadı: %w", err)
	}

	var order Order
	if err := json.Unmarshal(data, &order); err != nil {
		return nil, fmt.Errorf("JSON çözümleme hatası: %w", err)
	}

	return &order, nil
}

// DeleteOrder - Siparişi sil
func (s *BleveStore) DeleteOrder(id string) error {
	// Bleve'den sil
	if err := s.index.Delete(id); err != nil {
		return fmt.Errorf("index silme hatası: %w", err)
	}

	// Dosyadan sil
	orderFile := filepath.Join(s.dataPath, "orders", id+".json")
	if err := os.Remove(orderFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("dosya silme hatası: %w", err)
	}

	return nil
}

// ListOrders - Tüm siparişleri listele
func (s *BleveStore) ListOrders() ([]*Order, error) {
	ordersDir := filepath.Join(s.dataPath, "orders")

	// Dizin yoksa boş liste döndür
	if _, err := os.Stat(ordersDir); os.IsNotExist(err) {
		return []*Order{}, nil
	}

	entries, err := os.ReadDir(ordersDir)
	if err != nil {
		return nil, fmt.Errorf("dizin okunamadı: %w", err)
	}

	var orders []*Order
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		id := entry.Name()[:len(entry.Name())-5] // .json uzantısını kaldır
		order, err := s.GetOrder(id)
		if err != nil {
			continue
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// ListOrdersByDateRange - Tarih aralığına göre siparişleri listele
func (s *BleveStore) ListOrdersByDateRange(startDate, endDate time.Time) ([]*Order, error) {
	allOrders, err := s.ListOrders()
	if err != nil {
		return nil, err
	}

	var filteredOrders []*Order
	for _, order := range allOrders {
		// Tarihi yerel zaman diliminde gün başlangıcına normalize et
		loc := time.Local
		orderDate := time.Date(order.CreatedAt.Year(), order.CreatedAt.Month(), order.CreatedAt.Day(), 0, 0, 0, 0, loc)
		start := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, loc)
		end := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, loc)

		if (orderDate.Equal(start) || orderDate.After(start)) && (orderDate.Equal(end) || orderDate.Before(end)) {
			filteredOrders = append(filteredOrders, order)
		}
	}

	return filteredOrders, nil
}

// ListTodayOrders - Bugünkü siparişleri listele
func (s *BleveStore) ListTodayOrders() ([]*Order, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	return s.ListOrdersByDateRange(today, today)
}

// UpdateOrder - Mevcut siparişi güncelle (yeni kayıt oluşturmaz)
func (s *BleveStore) UpdateOrder(order *Order) error {
	// Mevcut siparişi kontrol et
	existingOrder, err := s.GetOrder(order.ID)
	if err != nil {
		return fmt.Errorf("sipariş bulunamadı: %w", err)
	}

	// CreatedAt'ı koru, UpdatedAt'ı güncelle
	order.CreatedAt = existingOrder.CreatedAt
	order.UpdatedAt = time.Now()

	// Toplamları hesapla
	order.CalculateGrandTotal()

	// JSON'a çevir
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("JSON dönüştürme hatası: %w", err)
	}

	// Bleve'e güncelle (aynı ID ile tekrar indexle)
	if err := s.index.Index(order.ID, order); err != nil {
		return fmt.Errorf("indexleme hatası: %w", err)
	}

	// JSON dosyasını güncelle
	orderFile := filepath.Join(s.dataPath, "orders", order.ID+".json")
	if err := os.WriteFile(orderFile, data, 0644); err != nil {
		return fmt.Errorf("dosya yazma hatası: %w", err)
	}

	return nil
}

// SearchOrders - Ürün adı veya OEM numarasına göre ara (Elasticsearch query)
func (s *BleveStore) SearchOrders(searchTerm string) ([]*Order, error) {
	// Bleve query oluştur - Elasticsearch query_string benzeri
	query := bleve.NewQueryStringQuery(searchTerm)

	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Size = 100

	searchResult, err := s.index.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("arama hatası: %w", err)
	}

	var orders []*Order
	for _, hit := range searchResult.Hits {
		order, err := s.GetOrder(hit.ID)
		if err != nil {
			continue
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// SearchOrdersAdvanced - Gelişmiş sipariş arama
func (s *BleveStore) SearchOrdersAdvanced(productName, oemNumber, customerName string, minQty, maxQty int, minTotal, maxTotal, minUnitPrice, maxUnitPrice float64, dateFilter, startDate, endDate string) ([]*Order, error) {
	allOrders, err := s.ListOrders()
	if err != nil {
		return nil, err
	}

	// Tarih filtreleme için zaman aralığını hesapla
	now := time.Now()
	var filterStartTime, filterEndTime time.Time

	switch dateFilter {
	case "today":
		filterStartTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		filterEndTime = filterStartTime.Add(24 * time.Hour)
	case "range":
		if startDate != "" {
			filterStartTime, _ = time.Parse("2006-01-02", startDate)
		}
		if endDate != "" {
			filterEndTime, _ = time.Parse("2006-01-02", endDate)
			filterEndTime = filterEndTime.Add(24 * time.Hour) // Gün sonuna kadar
		}
	case "all":
		// Tarih filtresi yok
	default:
		// Varsayılan olarak tarih filtresi yok
	}

	var filtered []*Order
	for _, order := range allOrders {
		// Tarih filtresi
		if dateFilter == "today" || dateFilter == "range" {
			if !filterStartTime.IsZero() && order.CreatedAt.Before(filterStartTime) {
				continue
			}
			if !filterEndTime.IsZero() && order.CreatedAt.After(filterEndTime) {
				continue
			}
		}

		// Müşteri adı filtresi
		if customerName != "" {
			if !strings.Contains(strings.ToLower(order.CustomerName), strings.ToLower(customerName)) {
				continue
			}
		}

		// Toplam tutar filtresi
		if minTotal > 0 && order.GrandTotal < minTotal {
			continue
		}
		if maxTotal > 0 && order.GrandTotal > maxTotal {
			continue
		}

		// Ürün/OEM/Adet/Birim Fiyat filtresi - en az bir ürün eşleşmeli
		matchProduct := productName == "" && oemNumber == "" && minQty == 0 && maxQty == 0 && minUnitPrice == 0 && maxUnitPrice == 0
		for _, item := range order.Items {
			productMatch := productName == "" || strings.Contains(strings.ToLower(item.ProductName), strings.ToLower(productName))
			oemMatch := oemNumber == "" || strings.Contains(strings.ToLower(item.OEMNumber), strings.ToLower(oemNumber))
			minQtyMatch := minQty == 0 || item.Quantity >= minQty
			maxQtyMatch := maxQty == 0 || item.Quantity <= maxQty
			minUnitPriceMatch := minUnitPrice == 0 || item.UnitPrice >= minUnitPrice
			maxUnitPriceMatch := maxUnitPrice == 0 || item.UnitPrice <= maxUnitPrice

			if productMatch && oemMatch && minQtyMatch && maxQtyMatch && minUnitPriceMatch && maxUnitPriceMatch {
				matchProduct = true
				break
			}
		}

		if matchProduct {
			filtered = append(filtered, order)
		}
	}

	// Tarihe göre sırala (yeniden eskiye)
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
	})

	return filtered, nil
}

// CalculateTotalPrice - Ürün toplam fiyatını hesaplar
func (item *OrderItem) CalculateTotalPrice() {
	item.TotalPrice = float64(item.Quantity) * item.UnitPrice
}

// CalculateGrandTotal - Siparişin genel toplamını hesaplar
func (order *Order) CalculateGrandTotal() {
	var total float64
	for i := range order.Items {
		order.Items[i].CalculateTotalPrice()
		total += order.Items[i].TotalPrice
	}
	order.GrandTotal = total
}

// NewOrderItem - Yeni sipariş kalemi oluşturur
func NewOrderItem(productName, oemNumber string, quantity int, unitPrice float64, partStatus string) OrderItem {
	item := OrderItem{
		ID:          uuid.New().String(),
		ProductName: productName,
		OEMNumber:   oemNumber,
		Quantity:    quantity,
		UnitPrice:   unitPrice,
		PartStatus:  partStatus,
	}
	item.CalculateTotalPrice()
	return item
}

// NewOrder - Yeni sipariş oluşturur
func NewOrder() *Order {
	return &Order{
		ID:           uuid.New().String(),
		CustomerID:   "",
		CustomerName: "",
		Items:        []OrderItem{},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// NewOrderWithCustomer - Müşteri ile yeni sipariş oluşturur
func NewOrderWithCustomer(customerID, customerName string) *Order {
	return &Order{
		ID:           uuid.New().String(),
		CustomerID:   customerID,
		CustomerName: customerName,
		Items:        []OrderItem{},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// ============================================
// Müşteri İşlemleri
// ============================================

// SaveCustomer - Müşteriyi kaydet
func (s *BleveStore) SaveCustomer(customer *Customer) error {
	customer.UpdatedAt = time.Now()

	data, err := json.Marshal(customer)
	if err != nil {
		return fmt.Errorf("JSON dönüştürme hatası: %w", err)
	}

	// Bleve'e indexle
	if err := s.index.Index("customer_"+customer.ID, customer); err != nil {
		return fmt.Errorf("indexleme hatası: %w", err)
	}

	// JSON dosyası olarak sakla
	customerFile := filepath.Join(s.dataPath, "customers", customer.ID+".json")
	if err := os.MkdirAll(filepath.Dir(customerFile), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(customerFile, data, 0644); err != nil {
		return fmt.Errorf("dosya yazma hatası: %w", err)
	}

	return nil
}

// GetCustomer - Müşteriyi getir
func (s *BleveStore) GetCustomer(id string) (*Customer, error) {
	customerFile := filepath.Join(s.dataPath, "customers", id+".json")
	data, err := os.ReadFile(customerFile)
	if err != nil {
		return nil, fmt.Errorf("müşteri bulunamadı: %w", err)
	}

	var customer Customer
	if err := json.Unmarshal(data, &customer); err != nil {
		return nil, fmt.Errorf("JSON çözümleme hatası: %w", err)
	}

	return &customer, nil
}

// ListCustomers - Tüm müşterileri listele
func (s *BleveStore) ListCustomers() ([]*Customer, error) {
	customersDir := filepath.Join(s.dataPath, "customers")

	if _, err := os.Stat(customersDir); os.IsNotExist(err) {
		return []*Customer{}, nil
	}

	entries, err := os.ReadDir(customersDir)
	if err != nil {
		return nil, fmt.Errorf("dizin okunamadı: %w", err)
	}

	var customers []*Customer
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		id := entry.Name()[:len(entry.Name())-5]
		customer, err := s.GetCustomer(id)
		if err != nil {
			continue
		}
		customers = append(customers, customer)
	}

	// İsme göre sırala
	sort.Slice(customers, func(i, j int) bool {
		return strings.ToLower(customers[i].Name) < strings.ToLower(customers[j].Name)
	})

	return customers, nil
}

// SearchCustomers - Müşteri adına göre ara
func (s *BleveStore) SearchCustomers(searchTerm string) ([]*Customer, error) {
	customers, err := s.ListCustomers()
	if err != nil {
		return nil, err
	}

	if searchTerm == "" {
		return customers, nil
	}

	searchLower := strings.ToLower(searchTerm)
	var filtered []*Customer
	for _, c := range customers {
		if strings.Contains(strings.ToLower(c.Name), searchLower) {
			filtered = append(filtered, c)
		}
	}

	return filtered, nil
}

// GetOrCreateCustomer - Müşteri ID'si varsa getir, yoksa isimle yeni oluştur
// Aynı isimde farklı müşteriler olabilir - ID bazlı çalışır
func (s *BleveStore) GetOrCreateCustomer(name string, phone string, customerID string) (*Customer, error) {
	if name == "" {
		return nil, nil
	}

	// Eğer customerID verilmişse, o müşteriyi getir
	if customerID != "" {
		customer, err := s.GetCustomer(customerID)
		if err == nil && customer != nil {
			// İsim veya telefon değişmişse güncelle
			updated := false
			if customer.Name != strings.TrimSpace(name) {
				customer.Name = strings.TrimSpace(name)
				updated = true
			}
			if customer.Phone != strings.TrimSpace(phone) {
				customer.Phone = strings.TrimSpace(phone)
				updated = true
			}
			if updated {
				customer.UpdatedAt = time.Now()
				s.SaveCustomer(customer)
			}
			return customer, nil
		}
	}

	// Yeni müşteri oluştur (aynı isimde olsa bile)
	customer := &Customer{
		ID:        uuid.New().String(),
		Name:      strings.TrimSpace(name),
		Phone:     strings.TrimSpace(phone),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.SaveCustomer(customer); err != nil {
		return nil, err
	}

	return customer, nil
}

// GetCustomerOrders - Müşterinin siparişlerini getir
func (s *BleveStore) GetCustomerOrders(customerID string) ([]*Order, error) {
	allOrders, err := s.ListOrders()
	if err != nil {
		return nil, err
	}

	var customerOrders []*Order
	for _, order := range allOrders {
		if order.CustomerID == customerID {
			customerOrders = append(customerOrders, order)
		}
	}

	// Tarihe göre sırala (yeniden eskiye)
	sort.Slice(customerOrders, func(i, j int) bool {
		return customerOrders[i].CreatedAt.After(customerOrders[j].CreatedAt)
	})

	return customerOrders, nil
}

// UpdateCustomerStats - Müşteri istatistiklerini güncelle
func (s *BleveStore) UpdateCustomerStats(customerID string) error {
	if customerID == "" {
		return nil
	}

	customer, err := s.GetCustomer(customerID)
	if err != nil {
		return nil
	}

	orders, err := s.GetCustomerOrders(customerID)
	if err != nil {
		return err
	}

	customer.OrderCount = len(orders)
	customer.TotalAmount = 0
	for _, order := range orders {
		customer.TotalAmount += order.GrandTotal
	}

	return s.SaveCustomer(customer)
}

// NewCustomer - Yeni müşteri oluşturur
func NewCustomer(name string) *Customer {
	return &Customer{
		ID:        uuid.New().String(),
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// ============================================
// Ürün Kataloğu İşlemleri
// ============================================

// SaveProduct - Ürünü kaydet
func (s *BleveStore) SaveProduct(product *Product) error {
	product.UpdatedAt = time.Now()

	data, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("JSON dönüştürme hatası: %w", err)
	}

	// Bleve'e indexle
	if err := s.index.Index("product_"+product.ID, product); err != nil {
		return fmt.Errorf("indexleme hatası: %w", err)
	}

	// JSON dosyası olarak sakla
	productFile := filepath.Join(s.dataPath, "products", product.ID+".json")
	if err := os.MkdirAll(filepath.Dir(productFile), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(productFile, data, 0644); err != nil {
		return fmt.Errorf("dosya yazma hatası: %w", err)
	}

	return nil
}

// GetProduct - Ürünü getir
func (s *BleveStore) GetProduct(id string) (*Product, error) {
	productFile := filepath.Join(s.dataPath, "products", id+".json")
	data, err := os.ReadFile(productFile)
	if err != nil {
		return nil, fmt.Errorf("ürün bulunamadı: %w", err)
	}

	var product Product
	if err := json.Unmarshal(data, &product); err != nil {
		return nil, fmt.Errorf("JSON çözümleme hatası: %w", err)
	}

	return &product, nil
}

// ListProducts - Tüm ürünleri listele
func (s *BleveStore) ListProducts() ([]*Product, error) {
	productsDir := filepath.Join(s.dataPath, "products")

	if _, err := os.Stat(productsDir); os.IsNotExist(err) {
		return []*Product{}, nil
	}

	entries, err := os.ReadDir(productsDir)
	if err != nil {
		return nil, fmt.Errorf("dizin okunamadı: %w", err)
	}

	var products []*Product
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		id := entry.Name()[:len(entry.Name())-5]
		product, err := s.GetProduct(id)
		if err != nil {
			continue
		}
		products = append(products, product)
	}

	// Kullanım sayısına göre sırala (çok kullanılan önce)
	sort.Slice(products, func(i, j int) bool {
		return products[i].UsedCount > products[j].UsedCount
	})

	return products, nil
}

// SearchProducts - Ürün adı veya OEM'e göre ara
func (s *BleveStore) SearchProducts(searchTerm string) ([]*Product, error) {
	products, err := s.ListProducts()
	if err != nil {
		return nil, err
	}

	if searchTerm == "" {
		return products, nil
	}

	searchLower := strings.ToLower(searchTerm)
	var filtered []*Product
	for _, p := range products {
		if strings.Contains(strings.ToLower(p.Name), searchLower) ||
			strings.Contains(strings.ToLower(p.OEMNumber), searchLower) {
			filtered = append(filtered, p)
		}
	}

	return filtered, nil
}

// GetOrCreateProduct - Ürün varsa getir, yoksa oluştur
func (s *BleveStore) GetOrCreateProduct(name, oemNumber string) (*Product, error) {
	if name == "" {
		return nil, nil
	}

	products, err := s.ListProducts()
	if err != nil {
		return nil, err
	}

	// Aynı isim ve OEM ile ürün var mı kontrol et
	nameLower := strings.ToLower(strings.TrimSpace(name))
	oemLower := strings.ToLower(strings.TrimSpace(oemNumber))
	for _, p := range products {
		if strings.ToLower(p.Name) == nameLower && strings.ToLower(p.OEMNumber) == oemLower {
			// Kullanım sayısını artır
			p.UsedCount++
			s.SaveProduct(p)
			return p, nil
		}
	}

	// Yoksa yeni oluştur
	product := &Product{
		ID:        uuid.New().String(),
		Name:      strings.TrimSpace(name),
		OEMNumber: strings.TrimSpace(oemNumber),
		UsedCount: 1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.SaveProduct(product); err != nil {
		return nil, err
	}

	return product, nil
}

// IncrementProductUsage - Ürün kullanım sayısını artır
func (s *BleveStore) IncrementProductUsage(productID string) error {
	product, err := s.GetProduct(productID)
	if err != nil {
		return err
	}

	product.UsedCount++
	return s.SaveProduct(product)
}

// UpdateProduct - Ürün bilgilerini güncelle
func (s *BleveStore) UpdateProduct(id, name, oemNumber string) error {
	product, err := s.GetProduct(id)
	if err != nil {
		return err
	}

	product.Name = strings.TrimSpace(name)
	product.OEMNumber = strings.TrimSpace(oemNumber)
	product.UpdatedAt = time.Now()

	return s.SaveProduct(product)
}

// DeleteProduct - Ürünü sil
func (s *BleveStore) DeleteProduct(id string) error {
	// Index'ten sil
	if err := s.index.Delete(id); err != nil {
		return err
	}

	// Dosyayı sil
	productFile := filepath.Join(s.dataPath, "products", id+".json")
	return os.Remove(productFile)
}

// UpdateCustomer - Müşteri bilgilerini güncelle
func (s *BleveStore) UpdateCustomer(id, name, phone string) error {
	customer, err := s.GetCustomer(id)
	if err != nil {
		return err
	}

	customer.Name = strings.TrimSpace(name)
	customer.Phone = strings.TrimSpace(phone)
	customer.UpdatedAt = time.Now()

	return s.SaveCustomer(customer)
}

// DeleteCustomer - Müşteriyi sil
func (s *BleveStore) DeleteCustomer(id string) error {
	// Index'ten sil
	if err := s.index.Delete(id); err != nil {
		return err
	}

	// Dosyayı sil
	customerFile := filepath.Join(s.dataPath, "customers", id+".json")
	return os.Remove(customerFile)
}
