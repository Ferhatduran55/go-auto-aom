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
  },

  // ============================================
  // Not İşlemleri
  // ============================================

  // Not kaydet
  async saveNote(noteData) {
    if (typeof saveNote !== 'undefined') {
      const result = await saveNote(JSON.stringify(noteData))
      return JSON.parse(result)
    }
    return { error: 'API not available' }
  },

  // Tüm notları yükle (opsiyonel arama terimi ile)
  async loadNotes(searchTerm = '') {
    if (typeof loadNotes !== 'undefined') {
      const result = await loadNotes(searchTerm)
      return JSON.parse(result) || []
    }
    return []
  },

  // Tek not yükle
  async loadNoteById(id) {
    if (typeof loadNoteById !== 'undefined') {
      const result = await loadNoteById(id)
      return JSON.parse(result)
    }
    return { error: 'API not available' }
  },

  // Not sil
  async deleteNote(id) {
    if (typeof deleteNote !== 'undefined') {
      const result = await deleteNote(id)
      return JSON.parse(result)
    }
    return { error: 'API not available' }
  },

  // Notlarda ara
  async searchNotes(searchTerm) {
    if (typeof searchNotes !== 'undefined') {
      const result = await searchNotes(searchTerm)
      return JSON.parse(result) || []
    }
    return []
  },

  // Metni otomatik formatla
  async autoFormatText(rawText) {
    if (typeof autoFormatText !== 'undefined') {
      const result = await autoFormatText(rawText)
      return JSON.parse(result)
    }
    return { formatted: rawText }
  },

  // ============================================
  // WhatsApp Sipariş İşlemleri
  // ============================================

  // WhatsApp siparişi kaydet
  async saveWhatsAppOrder(orderData) {
    if (typeof saveWhatsAppOrder !== 'undefined') {
      const result = await saveWhatsAppOrder(JSON.stringify(orderData))
      return JSON.parse(result)
    }
    return { error: 'API not available' }
  },

  // Tüm WhatsApp siparişlerini yükle (opsiyonel arama terimi ile)
  async loadWhatsAppOrders(searchTerm = '') {
    if (typeof loadWhatsAppOrders !== 'undefined') {
      const result = await loadWhatsAppOrders(searchTerm)
      return JSON.parse(result) || []
    }
    return []
  },

  // Tek WhatsApp siparişi yükle
  async loadWhatsAppOrderById(id) {
    if (typeof loadWhatsAppOrderById !== 'undefined') {
      const result = await loadWhatsAppOrderById(id)
      return JSON.parse(result)
    }
    return { error: 'API not available' }
  },

  // WhatsApp siparişi sil
  async deleteWhatsAppOrder(id) {
    if (typeof deleteWhatsAppOrder !== 'undefined') {
      const result = await deleteWhatsAppOrder(id)
      return JSON.parse(result)
    }
    return { error: 'API not available' }
  },

  // WhatsApp siparişlerinde ara
  async searchWhatsAppOrders(searchTerm) {
    if (typeof searchWhatsAppOrders !== 'undefined') {
      const result = await searchWhatsAppOrders(searchTerm)
      return JSON.parse(result) || []
    }
    return []
  },

  // WhatsApp siparişine durum ekle
  async addWhatsAppOrderStatus(orderId, status, note = '') {
    if (typeof addWhatsAppOrderStatus !== 'undefined') {
      const result = await addWhatsAppOrderStatus(JSON.stringify({
        order_id: orderId,
        status: status,
        note: note
      }))
      return JSON.parse(result)
    }
    return { error: 'API not available' }
  },

  // WhatsApp siparişlerini duruma göre filtrele
  async filterWhatsAppOrdersByStatus(status) {
    if (typeof filterWhatsAppOrdersByStatus !== 'undefined') {
      const result = await filterWhatsAppOrdersByStatus(status)
      return JSON.parse(result) || []
    }
    return []
  },

  // ==================== Güncelleme İşlemleri ====================

  async checkForUpdates() {
    if (typeof checkForUpdates !== 'undefined') {
      try {
        const result = await checkForUpdates()
        return JSON.parse(result)
      } catch (e) {
        console.error('[API] checkForUpdates error:', e)
        return { error: e.message }
      }
    }
    return { error: 'API not available' }
  },

  async getAppVersion() {
    if (typeof window.getAppVersion !== 'undefined') {
      try {
        return await window.getAppVersion()
      } catch (e) {
        console.error('[API] getAppVersion error:', e)
        return typeof __APP_VERSION__ !== 'undefined' ? __APP_VERSION__ : '0.0.0'
      }
    }
    return typeof __APP_VERSION__ !== 'undefined' ? __APP_VERSION__ : '0.0.0'
  },

  // ==================== Log İşlemleri ====================

  async getBackendLogs() {
    if (typeof getBackendLogs !== 'undefined') {
      try {
        const result = await getBackendLogs()
        return JSON.parse(result)
      } catch (e) {
        console.error('[API] getBackendLogs error:', e)
        return []
      }
    }
    return []
  },

  async getBackendLogsAfterId(afterId) {
    if (typeof getBackendLogsAfterId !== 'undefined') {
      try {
        const result = await getBackendLogsAfterId(afterId)
        return JSON.parse(result)
      } catch (e) {
        console.error('[API] getBackendLogsAfterId error:', e)
        return []
      }
    }
    return []
  },

  async clearBackendLogs() {
    if (typeof clearBackendLogs !== 'undefined') {
      try {
        const result = await clearBackendLogs()
        return JSON.parse(result)
      } catch (e) {
        console.error('[API] clearBackendLogs error:', e)
        return { success: false }
      }
    }
    return { success: false }
  },

  async setLogBufferSize(size) {
    if (typeof setLogBufferSize !== 'undefined') {
      try {
        const result = await setLogBufferSize(size)
        return JSON.parse(result)
      } catch (e) {
        console.error('[API] setLogBufferSize error:', e)
        return { success: false }
      }
    }
    return { success: false }
  },

  async getLogBufferSize() {
    if (typeof getLogBufferSize !== 'undefined') {
      try {
        const result = await getLogBufferSize()
        return JSON.parse(result)
      } catch (e) {
        console.error('[API] getLogBufferSize error:', e)
        return { size: 1000, count: 0 }
      }
    }
    return { size: 1000, count: 0 }
  },

  // ==================== Sistem İşlemleri ====================

  async getSystemInfo() {
    if (typeof getSystemInfo !== 'undefined') {
      try {
        const result = await getSystemInfo()
        return JSON.parse(result)
      } catch (e) {
        console.error('[API] getSystemInfo error:', e)
        return null
      }
    }
    return null
  },

  async exportAllData() {
    if (typeof exportAllData !== 'undefined') {
      try {
        const result = await exportAllData()
        return result // JSON string olarak döner
      } catch (e) {
        console.error('[API] exportAllData error:', e)
        return null
      }
    }
    return null
  },

  async copySystemInfo() {
    if (typeof copySystemInfo !== 'undefined') {
      try {
        return await copySystemInfo()
      } catch (e) {
        console.error('[API] copySystemInfo error:', e)
        return null
      }
    }
    return null
  },

  async clearCache() {
    if (typeof clearCache !== 'undefined') {
      try {
        const result = await clearCache()
        return JSON.parse(result)
      } catch (e) {
        console.error('[API] clearCache error:', e)
        return { success: false, error: e.message }
      }
    }
    return { success: false, error: 'API not available' }
  }
}

