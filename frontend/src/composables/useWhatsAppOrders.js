import { ref, computed } from 'vue'
import { api } from '@/api'

// ============================================
// WhatsApp Sipariş Yönetimi Composable
// ============================================

// Durum sabitleri
export const ORDER_STATUSES = {
  tedarik_surecinde: { label: 'Tedarik Sürecinde', color: 'primary', emoji: '🔵' },
  hazirlaniyor: { label: 'Hazırlanıyor', color: 'warning', emoji: '🟡' },
  hazirlandi: { label: 'Hazırlandı', color: 'info', emoji: '🟢' },
  kargoya_verildi: { label: 'Kargoya Verildi', color: 'purple', emoji: '🟣' },
  tamamlandi: { label: 'Tamamlandı', color: 'success', emoji: '✅' },
  iade: { label: 'İade', color: 'danger', emoji: '🔴' },
  hukumsuz: { label: 'Hükümsüz', color: 'gray', emoji: '⚫' }
}

// Ödeme yöntemleri
export const PAYMENT_METHODS = {
  havale: { label: 'Havale', emoji: '🏦' },
  nakit: { label: 'Nakit', emoji: '💵' },
  kredi_karti: { label: 'Kredi Kartı', emoji: '💳' }
}

// Parça türleri
export const PART_TYPES = {
  cikma: { label: 'Çıkma', color: 'primary', emoji: '🔵' },
  sifir: { label: 'Sıfır', color: 'success', emoji: '🟢' }
}

// Reaktif state
const orders = ref([])
const currentOrder = ref(null)
const loading = ref(false)
const searchTerm = ref('')
const statusFilter = ref('all')

// Computed
const filteredOrders = computed(() => {
  let result = orders.value

  // Durum filtrelemesi
  if (statusFilter.value && statusFilter.value !== 'all') {
    result = result.filter(order => order.current_status === statusFilter.value)
  }

  // Arama filtrelemesi (client-side için yedek)
  if (searchTerm.value) {
    const term = searchTerm.value.toLowerCase()
    result = result.filter(order => 
      order.customer_name?.toLowerCase().includes(term) ||
      order.customer_phone?.toLowerCase().includes(term) ||
      order.items?.some(item => item.product_name?.toLowerCase().includes(term))
    )
  }

  return result
})

const hasUnsavedChanges = computed(() => {
  if (!currentOrder.value) return false
  return currentOrder.value.id === '' || currentOrder.value._modified
})

// Actions
async function loadOrders(search = '') {
  loading.value = true
  try {
    orders.value = await api.loadWhatsAppOrders(search)
  } catch (e) {
    console.error('WhatsApp siparişleri yüklenemedi:', e)
    orders.value = []
  } finally {
    loading.value = false
  }
}

async function loadOrderById(id) {
  try {
    const order = await api.loadWhatsAppOrderById(id)
    if (!order.error) {
      currentOrder.value = { ...order, _modified: false }
    }
    return order
  } catch (e) {
    console.error('Sipariş yüklenemedi:', e)
    return { error: e.message }
  }
}

async function saveOrder() {
  if (!currentOrder.value) return { error: 'Sipariş yok' }
  
  try {
    const orderData = {
      id: currentOrder.value.id,
      date: currentOrder.value.date,
      customer_name: currentOrder.value.customer_name || '',
      customer_phone: currentOrder.value.customer_phone || '',
      payment_method: currentOrder.value.payment_method || '',
      payment_note: currentOrder.value.payment_note || '',
      items: currentOrder.value.items || []
    }

    const result = await api.saveWhatsAppOrder(orderData)
    
    if (result.success) {
      currentOrder.value.id = result.id
      currentOrder.value._modified = false
      await loadOrders()
    }
    
    return result
  } catch (e) {
    console.error('Sipariş kaydedilemedi:', e)
    return { error: e.message }
  }
}

async function deleteOrder(id) {
  try {
    const result = await api.deleteWhatsAppOrder(id)
    if (result.success) {
      if (currentOrder.value?.id === id) {
        currentOrder.value = null
      }
      await loadOrders()
    }
    return result
  } catch (e) {
    console.error('Sipariş silinemedi:', e)
    return { error: e.message }
  }
}

async function addStatus(orderId, status, note = '') {
  try {
    const result = await api.addWhatsAppOrderStatus(orderId, status, note)
    if (result.success) {
      // Siparişi yeniden yükle
      await loadOrderById(orderId)
      await loadOrders()
    }
    return result
  } catch (e) {
    console.error('Durum eklenemedi:', e)
    return { error: e.message }
  }
}

async function filterByStatus(status) {
  statusFilter.value = status
  if (status && status !== 'all') {
    try {
      orders.value = await api.filterWhatsAppOrdersByStatus(status)
    } catch (e) {
      console.error('Filtreleme hatası:', e)
    }
  } else {
    await loadOrders()
  }
}

function newOrder() {
  const today = new Date().toISOString().split('T')[0]
  currentOrder.value = {
    id: '',
    date: today,
    customer_name: '',
    customer_phone: '',
    payment_method: 'havale',
    payment_note: '',
    items: [],
    status_history: [],
    current_status: '',
    grand_total: 0,
    _modified: false
  }
}

function selectOrder(order) {
  currentOrder.value = { ...order, _modified: false }
}

function clearSelection() {
  currentOrder.value = null
}

// Item yönetimi
function addItem() {
  if (!currentOrder.value) return
  
  currentOrder.value.items.push({
    id: '',
    product_name: '',
    quantity: 1,
    price: 0,
    type: 'cikma',
    total: 0
  })
  currentOrder.value._modified = true
  recalculateTotal()
}

function removeItem(index) {
  if (!currentOrder.value) return
  
  currentOrder.value.items.splice(index, 1)
  currentOrder.value._modified = true
  recalculateTotal()
}

function updateItem(index, field, value) {
  if (!currentOrder.value || !currentOrder.value.items[index]) return
  
  currentOrder.value.items[index][field] = value
  
  // Toplam hesapla
  const item = currentOrder.value.items[index]
  item.total = (item.quantity || 0) * (item.price || 0)
  
  currentOrder.value._modified = true
  recalculateTotal()
}

function recalculateTotal() {
  if (!currentOrder.value) return
  
  let grandTotal = 0
  for (const item of currentOrder.value.items) {
    item.total = (item.quantity || 0) * (item.price || 0)
    grandTotal += item.total
  }
  currentOrder.value.grand_total = grandTotal
}

function markAsModified() {
  if (currentOrder.value) {
    currentOrder.value._modified = true
  }
}

// Helper functions
function formatDate(dateString) {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('tr-TR', { 
    day: '2-digit', 
    month: '2-digit', 
    year: 'numeric' 
  })
}

function formatDateTime(dateString) {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('tr-TR', { 
    day: '2-digit', 
    month: '2-digit', 
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

function formatCurrency(amount) {
  return new Intl.NumberFormat('tr-TR', {
    style: 'currency',
    currency: 'TRY',
    minimumFractionDigits: 0,
    maximumFractionDigits: 2
  }).format(amount || 0)
}

function getStatusInfo(status) {
  return ORDER_STATUSES[status] || { label: status, color: 'gray', emoji: '⚪' }
}

function getPaymentMethodInfo(method) {
  return PAYMENT_METHODS[method] || { label: method, emoji: '💰' }
}

function getPartTypeInfo(type) {
  return PART_TYPES[type] || { label: type, color: 'gray', emoji: '⚪' }
}

// Export composable
export function useWhatsAppOrders() {
  return {
    // State
    orders,
    currentOrder,
    loading,
    searchTerm,
    statusFilter,
    
    // Computed
    filteredOrders,
    hasUnsavedChanges,
    
    // Actions
    loadOrders,
    loadOrderById,
    saveOrder,
    deleteOrder,
    addStatus,
    filterByStatus,
    newOrder,
    selectOrder,
    clearSelection,
    
    // Item management
    addItem,
    removeItem,
    updateItem,
    recalculateTotal,
    markAsModified,
    
    // Helpers
    formatDate,
    formatDateTime,
    formatCurrency,
    getStatusInfo,
    getPaymentMethodInfo,
    getPartTypeInfo,
    
    // Constants
    ORDER_STATUSES,
    PAYMENT_METHODS,
    PART_TYPES
  }
}
