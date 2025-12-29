import { ref, computed } from 'vue'
import { api } from '@/api'

// Reaktif state
const products = ref([])
const currentOrderId = ref(null)
const currentCustomerId = ref(null)
const orderTitle = ref('')
const customerName = ref('')
const customerPhone = ref('')
const allProducts = ref([])
const allCustomers = ref([])
const editingProductId = ref(null)
const customerFilterActive = ref(false) // Müşteri filtresi aktif mi?
const advancedSearchFilter = ref(null) // Gelişmiş arama filtreleri

// Settings
const settings = ref({
  autoDeductStock: false // Automatically deduct stock when order is saved
})

// Load settings (localStorage)
function loadSettings() {
  try {
    const saved = localStorage.getItem('orderSettings')
    if (saved) {
      settings.value = { ...settings.value, ...JSON.parse(saved) }
    }
  } catch (e) {
    console.error('Settings could not be loaded:', e)
  }
}

// Save settings
function saveSettings() {
  localStorage.setItem('orderSettings', JSON.stringify(settings.value))
}

// Computed
const grandTotal = computed(() => {
  return products.value.reduce((sum, p) => sum + p.total_price, 0)
})

const isEditing = computed(() => currentOrderId.value !== null)

// Data loading
async function loadData() {
  try {
    loadSettings() // Load settings too
    
    const [customers, prods] = await Promise.all([
      api.listCustomers(),
      api.listProducts()
    ])
    allCustomers.value = customers
    allProducts.value = prods
  } catch (e) {
    console.error('Data loading error:', e)
  }
}

// OEM çakışma kontrolü
function checkOEMConflict(name, oem) {
  if (!oem || oem === '-') return null
  
  // Bu kombinasyon zaten kayıtlı mı?
  const exactMatch = allProducts.value.find(p => 
    p.oem_number?.toLowerCase() === oem.toLowerCase() && 
    p.name.toLowerCase() === name.toLowerCase()
  )
  
  if (exactMatch) return null // Çakışma yok
  
  // Aynı OEM farklı isimle kayıtlı mı?
  const existingProduct = allProducts.value.find(p => 
    p.oem_number?.toLowerCase() === oem.toLowerCase() && 
    p.name.toLowerCase() !== name.toLowerCase()
  )
  
  return existingProduct || null
}

// Ürün ekleme
function addProduct(productData) {
  const { name, oem, quantity, unitPrice, partStatus } = productData
  
  products.value.push({
    id: Date.now().toString(),
    product_name: name,
    oem_number: oem || '-',
    quantity: quantity,
    unit_price: unitPrice,
    part_status: partStatus,
    total_price: quantity * unitPrice
  })

  // Ürünü kataloga ekle
  api.saveProduct({ name, oem_number: oem })
  loadData()
}

// Ürün güncelleme
function updateProduct(id, productData) {
  const index = products.value.findIndex(p => p.id === id)
  if (index === -1) return

  const { name, oem, quantity, unitPrice, partStatus } = productData
  products.value[index] = {
    ...products.value[index],
    product_name: name,
    oem_number: oem || '-',
    quantity: quantity,
    unit_price: unitPrice,
    part_status: partStatus,
    total_price: quantity * unitPrice
  }
  editingProductId.value = null
}

// Ürün silme
function deleteProduct(id) {
  products.value = products.value.filter(p => p.id !== id)
}

// Sipariş kaydetme
async function saveOrder() {
  if (products.value.length === 0) {
    return { error: 'Liste boş' }
  }

  const orderData = {
    id: currentOrderId.value || '',
    title: orderTitle.value,
    customer_id: currentCustomerId.value || '',
    customer_name: customerName.value,
    customer_phone: customerPhone.value,
    items: products.value
  }

  const result = await api.saveOrder(orderData)
  
  if (!result.error) {
    currentOrderId.value = result.id
    if (result.customer_id) {
      currentCustomerId.value = result.customer_id
    }
    
    // Auto deduct stock if enabled
    if (settings.value.autoDeductStock) {
      await deductStock()
    }
    
    await loadData()
  }
  
  return result
}

// Deduct stock for order items
async function deductStock() {
  // Prepare bulk stock out
  const bulkOut = []
  
  for (const item of products.value) {
    // Find product in catalog (by OEM or name)
    const catalogProduct = allProducts.value.find(p => 
      p.oem_number?.toLowerCase() === item.oem_number?.toLowerCase() ||
      p.name.toLowerCase() === item.product_name.toLowerCase()
    )
    
    if (catalogProduct && catalogProduct.id) {
      bulkOut.push({
        product_id: catalogProduct.id,
        amount: item.quantity
      })
    }
  }
  
  if (bulkOut.length > 0) {
    const note = `Order: ${orderTitle.value || 'Unnamed'} - ${customerName.value || 'No customer'}`
    try {
      const entries = bulkOut.map(c => ({ product_id: c.product_id, amount: c.amount, note: note }))
      await api.bulkStockOut(entries)
    } catch (e) {
      console.error('Stock deduction error:', e)
    }
  }
}

// Sipariş yükleme
async function loadOrder(orderId) {
  const order = await api.loadOrderById(orderId)
  
  if (order.error) {
    return { error: order.error }
  }

  currentOrderId.value = order.id
  currentCustomerId.value = order.customer_id || null
  orderTitle.value = order.title || ''
  customerName.value = order.customer_name || ''
  
  // Müşterinin telefon numarasını al
  if (order.customer_id) {
    const customer = allCustomers.value.find(c => c.id === order.customer_id)
    customerPhone.value = customer?.phone || ''
  } else {
    customerPhone.value = ''
  }
  
  products.value = (order.items || []).map((item, i) => ({
    id: item.id || (Date.now() + i).toString(),
    product_name: item.product_name,
    oem_number: item.oem_number || '-',
    quantity: item.quantity,
    unit_price: item.unit_price,
    part_status: item.part_status,
    total_price: item.total_price
  }))

  return order
}

// Yeni sipariş
function resetOrder() {
  products.value = []
  currentOrderId.value = null
  currentCustomerId.value = null
  orderTitle.value = ''
  customerName.value = ''
  customerPhone.value = ''
  editingProductId.value = null
  customerFilterActive.value = false // Filtreyi kapat
}

// Listeyi temizle
function clearProducts() {
  products.value = []
}

export function useOrder() {
  return {
    // State
    products,
    currentOrderId,
    currentCustomerId,
    orderTitle,
    customerName,
    customerPhone,
    allProducts,
    allCustomers,
    editingProductId,
    customerFilterActive,
    advancedSearchFilter,
    settings,
    
    // Computed
    grandTotal,
    isEditing,
    
    // Methods
    loadData,
    checkOEMConflict,
    addProduct,
    updateProduct,
    deleteProduct,
    saveOrder,
    loadOrder,
    resetOrder,
    clearProducts,
    deductStock,
    saveSettings,
  }
}
