import { ref, computed } from 'vue'
import { api } from '../api'

// Reactive state
const products = ref([])
const stockMovements = ref([])
const criticalStockProducts = ref([])
const categories = ref(['Yağ', 'Filtre', 'Sprey', 'Fren', 'Diğer'])
const brands = ref([])
const units = ref(['adet', 'litre', 'kutu', 'paket'])
const loading = ref(false)
const activeFilter = ref({
  search: '',
  category: '',
  onlyCritical: false
})

// Pagination state
const currentPage = ref(1)
const itemsPerPage = ref(25)
const totalProducts = ref(0)
const sortField = ref('')
const sortDir = ref('asc')

// Stock report state
const stockReport = ref(null)

// Computed properties
const totalPages = computed(() => Math.ceil(totalProducts.value / itemsPerPage.value) || 1)

// Client-side filtered products (for backward compatibility)
const filteredProducts = computed(() => {
  let result = [...products.value]
  
  // Search filter
  if (activeFilter.value.search) {
    const q = activeFilter.value.search.toLowerCase()
    result = result.filter(p => 
      p.name.toLowerCase().includes(q) ||
      (p.oem_number && p.oem_number.toLowerCase().includes(q)) ||
      (p.brand && p.brand.toLowerCase().includes(q))
    )
  }
  
  // Category filter
  if (activeFilter.value.category) {
    result = result.filter(p => p.category === activeFilter.value.category)
  }
  
  // Only critical stock
  if (activeFilter.value.onlyCritical) {
    result = result.filter(p => {
      const threshold = p.critical_stock || 3
      return p.stock_quantity < threshold
    })
  }
  
  return result
})

const criticalStockCount = computed(() => criticalStockProducts.value.length)

// Check if product is in critical stock
const isCriticalStock = (product) => {
  const threshold = product.critical_stock || 3
  return product.stock_quantity < threshold
}

// Get step value based on unit
const getUnitStep = (unit) => {
  return unit === 'litre' ? 0.1 : 1
}

// Format quantity based on unit
const formatQuantity = (quantity, unit) => {
  if (unit === 'litre') {
    return parseFloat(quantity).toFixed(1)
  }
  return Math.floor(quantity).toString()
}

// Normalize quantity based on unit
const normalizeQuantity = (quantity, unit) => {
  if (unit === 'litre') {
    return Math.round(parseFloat(quantity) * 10) / 10
  }
  return Math.floor(parseFloat(quantity))
}

// API çağrıları
const loadProducts = async () => {
  loading.value = true
  try {
    products.value = await api.listProducts()
    await loadCriticalStockProducts()
  } catch (error) {
    console.error('Products could not be loaded:', error)
  } finally {
    loading.value = false
  }
}

// Server-side paginated loading
const loadProductsPaginated = async (options = {}) => {
  loading.value = true
  try {
    const result = await api.listProductsPaginated({
      page: options.page || currentPage.value,
      pageSize: options.pageSize || itemsPerPage.value,
      search: options.search ?? activeFilter.value.search,
      category: options.category ?? activeFilter.value.category,
      onlyCritical: options.onlyCritical ?? activeFilter.value.onlyCritical,
      sortField: options.sortField ?? sortField.value,
      sortDir: options.sortDir ?? sortDir.value
    })
    
    products.value = result.products || []
    totalProducts.value = result.total || 0
    currentPage.value = result.page || 1
    
    await loadCriticalStockProducts()
    return { success: true, total: result.total }
  } catch (error) {
    console.error('Products could not be loaded:', error)
    return { error: error.message }
  } finally {
    loading.value = false
  }
}

// Change page
const goToPage = async (page) => {
  if (page < 1 || page > totalPages.value) return
  currentPage.value = page
  await loadProductsPaginated({ page })
}

// Change page size
const setPageSize = async (size) => {
  itemsPerPage.value = size
  currentPage.value = 1 // Reset to first page
  await loadProductsPaginated({ page: 1, pageSize: size })
}

// Change sort
const setSort = async (field, dir = 'asc') => {
  sortField.value = field
  sortDir.value = dir
  currentPage.value = 1
  await loadProductsPaginated({ page: 1, sortField: field, sortDir: dir })
}

// Apply filter with pagination reset
const applyFilter = async (filter) => {
  activeFilter.value = { ...activeFilter.value, ...filter }
  currentPage.value = 1
  await loadProductsPaginated({ page: 1 })
}

const loadCriticalStockProducts = async () => {
  try {
    criticalStockProducts.value = await api.getCriticalStockProducts()
  } catch (error) {
    console.error('Critical stocks could not be loaded:', error)
  }
}

const loadCategories = async () => {
  try {
    const result = await api.getCategories()
    if (result && result.length > 0) {
      categories.value = result
    }
  } catch (error) {
    console.error('Categories could not be loaded:', error)
  }
}

const loadBrands = async () => {
  try {
    const result = await api.getBrands()
    if (result && result.length > 0) {
      brands.value = result
    }
  } catch (error) {
    console.error('Brands could not be loaded:', error)
  }
}

const loadUnits = async () => {
  try {
    const result = await api.getUnits()
    if (result && result.length > 0) {
      units.value = result
    }
  } catch (error) {
    console.error('Units could not be loaded:', error)
  }
}

const loadStockMovements = async (productId = '', start = '', end = '') => {
  loading.value = true
  try {
    stockMovements.value = await api.getStockMovements(productId, start, end)
  } catch (error) {
    console.error('Stock movements could not be loaded:', error)
  } finally {
    loading.value = false
  }
}

// Stock operations
const handleStockIn = async (productId, amount, note) => {
  try {
    const result = await api.stockIn(productId, amount, note)
    if (result.success) {
      await loadProducts()
      return { success: true }
    }
    return { error: result.error || 'Stock in failed' }
  } catch (error) {
    return { error: error.message }
  }
}

const handleStockOut = async (productId, amount, note) => {
  try {
    const result = await api.stockOut(productId, amount, note)
    if (result.success) {
      await loadProducts()
      return { success: true }
    }
    return { error: result.error || 'Stock out failed' }
  } catch (error) {
    return { error: error.message }
  }
}

const bulkStockIn = async (entries) => {
  try {
    const result = await api.bulkStockIn(entries)
    if (result.success) {
      await loadProducts()
      return { 
        success: true, 
        basarili: result.basarili, 
        hatalar: result.hatalar 
      }
    }
    return { error: result.error || 'Bulk stock in failed' }
  } catch (error) {
    return { error: error.message }
  }
}

const bulkStockOut = async (entries) => {
  try {
    const result = await api.bulkStockOut(entries)
    if (result.success) {
      await loadProducts()
      return { 
        success: true, 
        basarili: result.basarili, 
        hatalar: result.hatalar 
      }
    }
    return { error: result.error || 'Bulk stock out failed' }
  } catch (error) {
    return { error: error.message }
  }
}

// Ürün işlemleri
const createProduct = async (productData) => {
  try {
    const result = await api.createProductFull(productData)
    if (result && result.id) {
      await loadProducts()
      await loadCategories()
      return { success: true, product: result }
    }
    return { error: result.error || 'Product creation failed' }
  } catch (error) {
    return { error: error.message }
  }
}

const updateProduct = async (productData) => {
  try {
    const result = await api.updateProduct(productData)
    if (result.success) {
      await loadProducts()
      await loadCategories()
      return { success: true }
    }
    return { error: result.error || 'Product update failed' }
  } catch (error) {
    return { error: error.message }
  }
}

const deleteProduct = async (productId) => {
  try {
    const result = await api.deleteProduct(productId)
    if (result.success) {
      await loadProducts()
      return { success: true }
    }
    return { error: result.error || 'Product deletion failed' }
  } catch (error) {
    return { error: error.message }
  }
}

// Rapor işlemleri
const getStockReport = async (period, date) => {
  loading.value = true
  try {
    const result = await api.getStockReport(period, date)
    if (result && !result.error) {
      stockReport.value = result
      return { success: true, report: result }
    }
    return { error: result.error || 'Report generation failed' }
  } catch (error) {
    return { error: error.message }
  } finally {
    loading.value = false
  }
}

// Reset filters
const resetFilters = () => {
  activeFilter.value = {
    search: '',
    category: '',
    onlyCritical: false
  }
}

// Initial load
const initialLoad = async () => {
  await Promise.all([
    loadProducts(),
    loadCategories(),
    loadUnits()
  ])
}

export function useStock() {
  return {
    // State
    products,
    stockMovements,
    criticalStockProducts,
    categories,
    brands,
    units,
    loading,
    activeFilter,
    stockReport,
    
    // Pagination state
    currentPage,
    itemsPerPage,
    totalProducts,
    totalPages,
    sortField,
    sortDir,
    
    // Computed
    filteredProducts,
    criticalStockCount,
    criticalCount: criticalStockCount, // Alias for CriticalStockBadge
    
    // Helpers
    isCriticalStock,
    getUnitStep,
    formatQuantity,
    normalizeQuantity,
    
    // API operations
    loadProducts,
    loadProductsPaginated,
    loadCriticalStockProducts,
    loadCategories,
    loadBrands,
    loadUnits,
    loadStockMovements,
    
    // Pagination operations
    goToPage,
    setPageSize,
    setSort,
    applyFilter,
    
    // Stock operations
    handleStockIn,
    handleStockOut,
    bulkStockIn,
    bulkStockOut,
    
    // Product operations
    createProduct,
    updateProduct,
    deleteProduct,
    
    // Report
    getStockReport,
    generateReport: getStockReport,  // Alias for ReportsView.vue
    
    // Other
    resetFilters,
    initialLoad
  }
}
