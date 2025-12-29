// Go/WebView2 ile iletişim için API wrapper
// window.* fonksiyonları Go tarafından bind edilir

export const api = {
  // Sipariş işlemleri
  async saveOrder(orderData) {
    if (typeof saveOrderToBleve !== 'undefined') {
      const result = await saveOrderToBleve(JSON.stringify(orderData))
      return JSON.parse(result)
    }
    return { error: 'API not available' }
  },

  async loadOrders(filter) {
    if (typeof loadOrdersFromBleve !== 'undefined') {
      const result = await loadOrdersFromBleve(JSON.stringify(filter))
      return JSON.parse(result) || []
    }
    return []
  },

  async loadOrderById(id) {
    if (typeof loadOrderById !== 'undefined') {
      const result = await loadOrderById(id)
      return JSON.parse(result)
    }
    return { error: 'API not available' }
  },

  async deleteOrder(id) {
    if (typeof deleteOrderFromBleve !== 'undefined') {
      const result = await deleteOrderFromBleve(id)
      return JSON.parse(result)
    }
    return { error: 'API not available' }
  },

  async searchOrders(term) {
    if (typeof searchOrders !== 'undefined') {
      const result = await searchOrders(term)
      return JSON.parse(result) || []
    }
    return []
  },

  async searchOrdersAdvanced(filter) {
    if (typeof searchOrdersAdvanced !== 'undefined') {
      const result = await searchOrdersAdvanced(JSON.stringify(filter))
      return JSON.parse(result) || []
    }
    return []
  },

  // Müşteri işlemleri
  async listCustomers() {
    if (typeof listAllCustomers !== 'undefined') {
      const result = await listAllCustomers()
      return JSON.parse(result) || []
    }
    return []
  },

  async getCustomerOrders(customerId) {
    if (typeof getCustomerOrders !== 'undefined') {
      const result = await getCustomerOrders(customerId)
      return JSON.parse(result) || []
    }
    return []
  },

  // Ürün işlemleri
  async listProducts() {
    if (typeof listAllProducts !== 'undefined') {
      const result = await listAllProducts()
      return JSON.parse(result) || []
    }
    return []
  },

  // Sayfalı ürün listesi
  async listProductsPaginated(options = {}) {
    const params = {
      page: options.page || 1,
      page_size: options.pageSize || 25,
      search: options.search || '',
      category: options.category || '',
      only_critical: options.onlyCritical || false,
      sort_field: options.sortField || '',
      sort_dir: options.sortDir || 'asc'
    }
    if (typeof listProductsPaginated !== 'undefined') {
      const result = await listProductsPaginated(JSON.stringify(params))
      return JSON.parse(result) || { products: [], total: 0, page: 1, page_size: 25 }
    }
    return { products: [], total: 0, page: 1, page_size: 25 }
  },

  async saveProduct(productData) {
    if (typeof saveProduct !== 'undefined') {
      const result = await saveProduct(JSON.stringify(productData))
      return JSON.parse(result)
    }
    return { error: 'API not available' }
  },

  async updateProduct(productData) {
    if (typeof updateProduct !== 'undefined') {
      const result = await updateProduct(JSON.stringify(productData))
      return JSON.parse(result)
    }
    return { error: 'API not available' }
  },

  async deleteProduct(id) {
    if (typeof deleteProduct !== 'undefined') {
      const result = await deleteProduct(id)
      return JSON.parse(result)
    }
    return { error: 'API not available' }
  },

  async updateCustomer(customerData) {
    if (typeof updateCustomer !== 'undefined') {
      const result = await updateCustomer(JSON.stringify(customerData))
      return JSON.parse(result)
    }
    return { error: 'API not available' }
  },

  async deleteCustomer(id) {
    if (typeof deleteCustomer !== 'undefined') {
      const result = await deleteCustomer(id)
      return JSON.parse(result)
    }
    return { error: 'API not available' }
  },

  // ============================================
  // Stok Yönetimi İşlemleri
  // ============================================

  // Stok girişi yap
  async stockIn(productId, amount, note) {
    if (typeof stockIn !== 'undefined') {
      const result = await stockIn(JSON.stringify({
        product_id: productId,
        amount: amount,
        note: note
      }))
      return JSON.parse(result)
    }
    return { error: 'API not available' }
  },


  // Stok çıkışı yap
  async stockOut(productId, amount, note) {
    if (typeof stockOut !== 'undefined') {
      const result = await stockOut(JSON.stringify({
        product_id: productId,
        amount: amount,
        note: note
      }))
      return JSON.parse(result)
    }
    return { error: 'API not available' }
  },

  // Toplu stok girişi
  async bulkStockIn(entries) {
    if (typeof bulkStockIn !== 'undefined') {
      const result = await bulkStockIn(JSON.stringify(entries))
      return JSON.parse(result)
    }
    return { error: 'API not available' }
  },

  // Toplu stok çıkışı
  async bulkStockOut(entries) {
    if (typeof bulkStockOut !== 'undefined') {
      const result = await bulkStockOut(JSON.stringify(entries))
      return JSON.parse(result)
    }
    return { error: 'API not available' }
  },

  // Get stock movements
  async getStockMovements(productId = '', start = '', end = '') {
    if (typeof getStockMovements !== 'undefined') {
      const result = await getStockMovements(JSON.stringify({
        product_id: productId,
        start: start,
        end: end
      }))
      return JSON.parse(result) || []
    }
    return []
  },

  // Kritik stok altındaki ürünleri getir
  async getCriticalStockProducts() {
    if (typeof getCriticalStockProducts !== 'undefined') {
      const result = await getCriticalStockProducts()
      return JSON.parse(result) || []
    }
    return []
  },

  // Generate stock report
  async getStockReport(period, date) {
    if (typeof getStockReport !== 'undefined') {
      const result = await getStockReport(JSON.stringify({
        period: period,
        date: date
      }))
      return JSON.parse(result)
    }
    return { error: 'API not available' }
  },

  // Kategorileri getir
  async getCategories() {
    if (typeof getCategories !== 'undefined') {
      const result = await getCategories()
      return JSON.parse(result) || []
    }
    return ['Yağ', 'Filtre', 'Sprey', 'Fren', 'Diğer']
  },

  // Markaları getir
  async getBrands() {
    if (typeof getBrands !== 'undefined') {
      const result = await getBrands()
      return JSON.parse(result) || []
    }
    return []
  },

  // Get units
  async getUnits() {
    if (typeof getUnits !== 'undefined') {
      const result = await getUnits()
      return JSON.parse(result) || []
    }
    return ['adet', 'litre', 'kutu', 'paket']
  },

  // Tüm alanlarla ürün oluştur
  async createProductFull(productData) {
    if (typeof createProductFull !== 'undefined') {
      const result = await createProductFull(JSON.stringify(productData))
      return JSON.parse(result)
    }
    return { error: 'API not available' }
  },

  // ============================================
  // Helper Functions
  // ============================================

  // Format quantity based on unit
  formatQuantity(quantity, unit) {
    if (unit === 'litre') {
      return parseFloat(quantity).toFixed(1)
    }
    return Math.floor(quantity).toString()
  },

  // Get input step value based on unit
  getInputStep(unit) {
    return unit === 'litre' ? '0.1' : '1'
  },

  // Normalize quantity based on unit
  normalizeQuantity(quantity, unit) {
    if (unit === 'litre') {
      return Math.round(parseFloat(quantity) * 10) / 10
    }
    return Math.floor(parseFloat(quantity))
  },

  // ============================================
  // Settings Functions
  // ============================================

  // Set developer mode (requires app restart)
  async setDeveloperMode(enabled) {
    if (typeof setDeveloperMode !== 'undefined') {
      const result = await setDeveloperMode(enabled)
      return JSON.parse(result)
    }
    return { error: 'API not available' }
  },

  // Get developer mode status
  async getDeveloperMode() {
    if (typeof getDeveloperMode !== 'undefined') {
      const result = await getDeveloperMode()
      return JSON.parse(result)
    }
    return { enabled: false }
  }
}
