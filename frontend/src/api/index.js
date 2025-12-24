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
}
