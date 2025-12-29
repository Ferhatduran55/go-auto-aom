<template>
  <div class="stock-list">
    <!-- Top Bar: Search and Filters -->
    <div class="filter-bar">
      <div class="search-box">
        <svg class="search-icon" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="11" cy="11" r="8"></circle>
          <path d="m21 21-4.35-4.35"></path>
        </svg>
        <input 
          type="text" 
          v-model="searchQuery" 
          placeholder="Ürün adı, OEM veya marka ara..."
          class="form-input search-input"
          @input="debouncedSearch"
        />
      </div>
      
      <div class="filter-controls">
        <AutocompleteSelect
          v-model="categoryFilter"
          :items="categoryOptions"
          placeholder="Tüm Kategoriler"
          @update:modelValue="handleFilterChange"
          class="filter-select"
        />
        
        <label class="critical-filter">
          <input type="checkbox" v-model="onlyCritical" @change="handleFilterChange" />
          <span class="critical-label">Sadece Kritik Stok</span>
        </label>
        
        <button @click="resetFilters" class="btn btn-secondary btn-sm" title="Filtreleri Temizle">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M3 6h18"></path>
            <path d="M7 12h10"></path>
            <path d="M10 18h4"></path>
          </svg>
        </button>
      </div>
      
      <div class="action-buttons">
        <!-- Bulk Operations -->
        <template v-if="selectedProducts.length > 0">
          <div class="bulk-indicator">
            {{ selectedProducts.length }} ürün seçildi
          </div>
          <button @click="$emit('bulk-stock-in', selectedProducts)" class="btn btn-success">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 5v14"></path>
              <path d="M5 12h14"></path>
            </svg>
            Toplu Stok Girişi
          </button>
          <button @click="$emit('bulk-stock-out', selectedProducts)" class="btn btn-warning">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M5 12h14"></path>
            </svg>
            Toplu Stok Çıkışı
          </button>
          <button @click="$emit('bulk-edit', selectedProducts)" class="btn btn-primary">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M17 3a2.85 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z"></path>
            </svg>
            Toplu Düzenle
          </button>
          <button @click="clearSelection" class="btn btn-secondary btn-sm">
            Seçimi Temizle
          </button>
        </template>
        
        <template v-else>
          <button @click="exportCSV" class="btn btn-secondary" title="CSV Dışa Aktar">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
              <polyline points="7 10 12 15 17 10"></polyline>
              <line x1="12" y1="15" x2="12" y2="3"></line>
            </svg>
            CSV
          </button>
          
          <!-- Column Visibility -->
          <div class="dropdown" ref="columnDropdownRef">
            <button @click="showColumnDropdown = !showColumnDropdown" class="btn btn-secondary" title="Sütunları Göster/Gizle">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="3" width="7" height="7"></rect>
                <rect x="14" y="3" width="7" height="7"></rect>
                <rect x="14" y="14" width="7" height="7"></rect>
                <rect x="3" y="14" width="7" height="7"></rect>
              </svg>
              Sütunlar
            </button>
            <div v-if="showColumnDropdown" class="dropdown-menu">
              <label v-for="col in allColumns" :key="col.key" class="dropdown-item">
                <input type="checkbox" v-model="visibleColumns" :value="col.key" />
                {{ col.label }}
              </label>
            </div>
          </div>
          
          <button @click="$emit('stokGiris')" class="btn btn-success">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 5v14"></path>
              <path d="M5 12h14"></path>
            </svg>
            Stok Girişi
          </button>
          <button @click="$emit('stokCikis')" class="btn btn-warning">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M5 12h14"></path>
            </svg>
            Stok Çıkışı
          </button>
          <button @click="$emit('yeniUrun')" class="btn btn-primary">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 5v14"></path>
              <path d="M5 12h14"></path>
            </svg>
            Yeni Ürün
          </button>
        </template>
      </div>
    </div>
    
    <!-- Products Table -->
    <div class="table-container">
      <table class="stock-table">
        <thead>
          <tr>
            <th class="checkbox-col">
              <input 
                type="checkbox" 
                :checked="allSelected" 
                :indeterminate="someSelected && !allSelected"
                @change="toggleSelectAll"
              />
            </th>
            <th 
              v-for="col in displayedColumns" 
              :key="col.key"
              :class="[col.align === 'center' ? 'text-center' : '', col.sortable ? 'sortable' : '']"
              @click="col.sortable && toggleSort(col.key)"
            >
              {{ col.label }}
              <span v-if="col.sortable" class="sort-icon">
                <template v-if="sortField === col.key">
                  {{ sortDir === 'asc' ? '▲' : '▼' }}
                </template>
                <template v-else>
                  <span class="sort-inactive">⇅</span>
                </template>
              </span>
            </th>
            <th class="text-center">İşlemler</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="loading">
            <td :colspan="displayedColumns.length + 2" class="loading-row">
              <div class="loading-spinner"></div>
              Yükleniyor...
            </td>
          </tr>
          <tr v-else-if="products.length === 0">
            <td :colspan="displayedColumns.length + 2" class="empty-row">
              <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
                <path d="m7.5 4.27 9 5.15"></path>
                <path d="M21 8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16Z"></path>
                <path d="m3.3 7 8.7 5 8.7-5"></path>
                <path d="M12 22V12"></path>
              </svg>
              <p>Henüz ürün bulunmuyor</p>
            </td>
          </tr>
          <tr 
            v-else
            v-for="product in products" 
            :key="product.id"
            :class="{ 'critical-row': isCriticalStock(product), 'selected-row': isSelected(product) }"
          >
            <td class="checkbox-col">
              <input 
                type="checkbox" 
                :checked="isSelected(product)" 
                @change="toggleSelect(product)"
              />
            </td>
            <td v-if="isColumnVisible('name')" class="product-name">
              <strong>{{ product.name }}</strong>
            </td>
            <td v-if="isColumnVisible('oem_number')" class="oem-no">{{ product.oem_number || '-' }}</td>
            <td v-if="isColumnVisible('brand')">{{ product.brand || '-' }}</td>
            <td v-if="isColumnVisible('category')">
              <span v-if="product.category" class="badge badge-category">{{ product.category }}</span>
              <span v-else>-</span>
            </td>
            <td v-if="isColumnVisible('unit')" class="text-center">{{ product.unit || 'adet' }}</td>
            <td v-if="isColumnVisible('stock_quantity')" class="text-center stock-quantity">
              <span :class="{ 'stock-critical': isCriticalStock(product), 'stock-normal': !isCriticalStock(product) }">
                {{ formatQuantity(product.stock_quantity, product.unit) }}
              </span>
            </td>
            <td v-if="isColumnVisible('critical_stock')" class="text-center">{{ product.critical_stock || 3 }}</td>
            <td v-if="isColumnVisible('status')" class="text-center">
              <span v-if="isCriticalStock(product)" class="badge badge-danger">
                <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"></path>
                  <path d="M12 9v4"></path>
                  <path d="M12 17h.01"></path>
                </svg>
                Kritik
              </span>
              <span v-else class="badge badge-success">Normal</span>
            </td>
            <td class="text-center actions-cell">
              <button @click="$emit('stock-in', product)" class="btn-icon btn-icon-success" title="Stok Girişi">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M12 5v14"></path>
                  <path d="M5 12h14"></path>
                </svg>
              </button>
              <button @click="$emit('stock-out', product)" class="btn-icon btn-icon-warning" title="Stok Çıkışı" :disabled="product.stock_quantity <= 0">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M5 12h14"></path>
                </svg>
              </button>
              <button @click="$emit('movements', product)" class="btn-icon btn-icon-info" title="Hareketler">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M12 20h9"></path>
                  <path d="M16.5 3.5a2.12 2.12 0 0 1 3 3L7 19l-4 1 1-4Z"></path>
                </svg>
              </button>
              <button @click="$emit('edit-product', product)" class="btn-icon" title="Düzenle">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M17 3a2.85 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z"></path>
                </svg>
              </button>
              <button @click="$emit('delete-product', product)" class="btn-icon btn-icon-danger" title="Sil">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M3 6h18"></path>
                  <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"></path>
                  <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"></path>
                </svg>
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    
    <!-- Footer with Pagination -->
    <div class="table-footer">
      <div class="stats">
        <span>Toplam: <strong>{{ totalProducts }}</strong> ürün</span>
        <span v-if="criticalCount > 0" class="critical-stat">
          <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"></path>
            <path d="M12 9v4"></path>
            <path d="M12 17h.01"></path>
          </svg>
          {{ criticalCount }} ürün kritik stokta
        </span>
      </div>
      
      <!-- Pagination Controls -->
      <div class="pagination">
        <select v-model.number="pageSize" @change="handlePageSizeChange" class="page-size-select">
          <option :value="10">10 / sayfa</option>
          <option :value="25">25 / sayfa</option>
          <option :value="50">50 / sayfa</option>
          <option :value="100">100 / sayfa</option>
        </select>
        
        <div class="page-info">
          Sayfa {{ currentPage }} / {{ totalPages }}
        </div>
        
        <div class="page-buttons">
          <button 
            class="page-btn" 
            :disabled="currentPage <= 1" 
            @click="goToPage(1)"
            title="İlk Sayfa"
          >
            ⟪
          </button>
          <button 
            class="page-btn" 
            :disabled="currentPage <= 1" 
            @click="goToPage(currentPage - 1)"
            title="Önceki Sayfa"
          >
            ◀
          </button>
          <button 
            class="page-btn" 
            :disabled="currentPage >= totalPages" 
            @click="goToPage(currentPage + 1)"
            title="Sonraki Sayfa"
          >
            ▶
          </button>
          <button 
            class="page-btn" 
            :disabled="currentPage >= totalPages" 
            @click="goToPage(totalPages)"
            title="Son Sayfa"
          >
            ⟫
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted, defineExpose } from 'vue'
import { useStock } from '../composables/useStock'
import { useSettings } from '../composables/useSettings'
import AutocompleteSelect from './AutocompleteSelect.vue'

const props = defineProps({
  filterCritical: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits([
  'stock-in', 'stock-out', 'new-product', 'edit-product', 'delete-product', 
  'movements', 'stokGiris', 'stokCikis', 'yeniUrun',
  'bulk-stock-in', 'bulk-stock-out', 'bulk-edit'
])

const { 
  products,
  categories,
  loading,
  totalProducts,
  currentPage,
  itemsPerPage,
  totalPages,
  sortField,
  sortDir,
  isCriticalStock,
  formatQuantity,
  loadProductsPaginated,
  goToPage: goToPageFn,
  setPageSize: setPageSizeFn,
  setSort
} = useStock()

const { settings } = useSettings()

// Local state
const searchQuery = ref('')
const categoryFilter = ref('')
const onlyCritical = ref(false)
const selectedProducts = ref([])
const showColumnDropdown = ref(false)
const columnDropdownRef = ref(null)
const pageSize = ref(settings.value.itemsPerPage || 25)

// Category options for autocomplete
const categoryOptions = computed(() => 
  categories.value.map(c => ({ label: c, value: c }))
)

// Column definitions
const allColumns = [
  { key: 'name', label: 'Ürün Adı', sortable: true },
  { key: 'oem_number', label: 'OEM No', sortable: true },
  { key: 'brand', label: 'Marka', sortable: true },
  { key: 'category', label: 'Kategori', sortable: true },
  { key: 'unit', label: 'Birim', sortable: false, align: 'center' },
  { key: 'stock_quantity', label: 'Stok', sortable: true, align: 'center' },
  { key: 'critical_stock', label: 'Kritik', sortable: true, align: 'center' },
  { key: 'status', label: 'Durum', sortable: false, align: 'center' }
]

// Convert settings to array format if needed
function getVisibleColumnsArray() {
  const cols = settings.value.stockListColumns
  if (Array.isArray(cols)) {
    return cols
  }
  if (cols && typeof cols === 'object') {
    return Object.keys(cols).filter(k => cols[k])
  }
  return allColumns.map(c => c.key)
}

const visibleColumns = ref(getVisibleColumnsArray())

// Computed
const displayedColumns = computed(() => 
  allColumns.filter(col => visibleColumns.value.includes(col.key))
)

const allSelected = computed(() => 
  products.value.length > 0 && selectedProducts.value.length === products.value.length
)

const someSelected = computed(() => 
  selectedProducts.value.length > 0
)

const criticalCount = computed(() => 
  products.value.filter(p => isCriticalStock(p)).length
)

// Methods
function isColumnVisible(key) {
  return visibleColumns.value.includes(key)
}

function isSelected(product) {
  return selectedProducts.value.some(p => p.id === product.id)
}

function toggleSelect(product) {
  const index = selectedProducts.value.findIndex(p => p.id === product.id)
  if (index >= 0) {
    selectedProducts.value.splice(index, 1)
  } else {
    selectedProducts.value.push(product)
  }
}

function toggleSelectAll() {
  if (allSelected.value) {
    selectedProducts.value = []
  } else {
    selectedProducts.value = [...products.value]
  }
}

function clearSelection() {
  selectedProducts.value = []
}

async function toggleSort(field) {
  if (sortField.value === field) {
    await setSort(field, sortDir.value === 'asc' ? 'desc' : 'asc')
  } else {
    await setSort(field, 'asc')
  }
}

async function goToPage(page) {
  await goToPageFn(page)
}

async function handlePageSizeChange() {
  await setPageSizeFn(pageSize.value)
}

function handleFilterChange() {
  currentPage.value = 1
  loadProducts()
}

let searchTimeout = null
function debouncedSearch() {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    currentPage.value = 1
    loadProducts()
  }, 300)
}

function resetFilters() {
  searchQuery.value = ''
  categoryFilter.value = ''
  onlyCritical.value = false
  currentPage.value = 1
  loadProducts()
}

function loadProducts() {
  loadProductsPaginated({
    search: searchQuery.value,
    category: categoryFilter.value,
    onlyCritical: onlyCritical.value
  })
}

function exportCSV() {
  const headers = displayedColumns.value.map(c => c.label)
  const rows = products.value.map(p => 
    displayedColumns.value.map(col => {
      if (col.key === 'status') {
        return isCriticalStock(p) ? 'Kritik' : 'Normal'
      }
      return p[col.key] || ''
    })
  )
  
  const csvContent = [
    headers.join(','),
    ...rows.map(row => row.map(cell => `"${cell}"`).join(','))
  ].join('\n')
  
  const blob = new Blob(['\ufeff' + csvContent], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.setAttribute('href', url)
  link.setAttribute('download', `stok_listesi_${new Date().toISOString().split('T')[0]}.csv`)
  link.style.visibility = 'hidden'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

// Close dropdown on outside click
function handleClickOutside(e) {
  if (columnDropdownRef.value && !columnDropdownRef.value.contains(e.target)) {
    showColumnDropdown.value = false
  }
}

// Watch for column visibility changes and save to settings
watch(visibleColumns, (newVal) => {
  settings.value.stockListColumns = newVal
  localStorage.setItem('appSettings', JSON.stringify(settings.value))
}, { deep: true })

// Watch filterCritical prop
watch(() => props.filterCritical, (newVal) => {
  if (newVal) {
    onlyCritical.value = true
    handleFilterChange()
  }
}, { immediate: true })

// Expose method to set critical filter from outside
function setCriticalFilter(value) {
  onlyCritical.value = value
  handleFilterChange()
}

defineExpose({
  setCriticalFilter
})

onMounted(() => {
  loadProducts()
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  clearTimeout(searchTimeout)
})
</script>

<style scoped>
.stock-list {
  display: flex;
  flex-direction: column;
  height: 100%;
  gap: 1rem;
}

.filter-bar {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  background: var(--bg-card);
  border-radius: 0.5rem;
  border: 1px solid var(--border-color);
  flex-wrap: wrap;
}

.search-box {
  position: relative;
  flex: 1;
  min-width: 250px;
}

.search-icon {
  position: absolute;
  left: 0.75rem;
  top: 50%;
  transform: translateY(-50%);
  color: var(--text-muted);
}

.search-input {
  padding-left: 2.5rem !important;
  width: 100%;
}

.filter-controls {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.filter-select {
  min-width: 150px;
}

.critical-filter {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  white-space: nowrap;
}

.critical-filter input {
  width: 1rem;
  height: 1rem;
  cursor: pointer;
}

.critical-label {
  font-size: 0.875rem;
  color: var(--text-primary);
}

.action-buttons {
  display: flex;
  gap: 0.5rem;
  margin-left: auto;
  align-items: center;
  flex-wrap: wrap;
}

.bulk-indicator {
  background: var(--accent-color);
  color: white;
  padding: 0.5rem 0.75rem;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  font-weight: 500;
}

/* Dropdown */
.dropdown {
  position: relative;
}

.dropdown-menu {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 0.25rem;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 0.5rem;
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.2);
  padding: 0.5rem;
  z-index: 100;
  min-width: 180px;
}

.dropdown-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem;
  cursor: pointer;
  border-radius: 0.25rem;
  font-size: 0.875rem;
  color: var(--text-primary);
}

.dropdown-item:hover {
  background: var(--bg-tertiary);
}

.dropdown-item input {
  cursor: pointer;
}

.table-container {
  flex: 1;
  overflow: auto;
  background: var(--bg-card);
  border-radius: 0.5rem;
  border: 1px solid var(--border-color);
}

.stock-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.875rem;
}

.stock-table th {
  background: var(--bg-secondary);
  padding: 0.75rem 1rem;
  text-align: left;
  font-weight: 600;
  color: var(--text-primary);
  border-bottom: 2px solid var(--border-color);
  position: sticky;
  top: 0;
  z-index: 10;
  user-select: none;
}

.stock-table th.sortable {
  cursor: pointer;
}

.stock-table th.sortable:hover {
  background: var(--bg-tertiary);
}

.sort-icon {
  margin-left: 0.5rem;
  font-size: 0.75rem;
}

.sort-inactive {
  opacity: 0.3;
}

.stock-table td {
  padding: 0.75rem 1rem;
  border-bottom: 1px solid var(--border-color);
  color: var(--text-primary);
}

.stock-table tbody tr:hover {
  background: var(--bg-secondary);
}

.checkbox-col {
  width: 40px;
  text-align: center;
}

.checkbox-col input {
  width: 1rem;
  height: 1rem;
  cursor: pointer;
}

.critical-row {
  background: var(--danger-bg) !important;
}

.critical-row:hover {
  background: rgba(239, 68, 68, 0.15) !important;
}

.selected-row {
  background: rgba(99, 102, 241, 0.1) !important;
}

.product-name {
  max-width: 200px;
}

.oem-no {
  font-family: monospace;
  font-size: 0.8rem;
}

.text-center {
  text-align: center;
}

.stock-quantity {
  font-weight: 600;
  font-size: 1rem;
}

.stock-critical {
  color: var(--danger-color);
}

.stock-normal {
  color: var(--success-color);
}

.badge {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.25rem 0.5rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 500;
}

.badge-category {
  background: var(--accent-color);
  color: white;
}

.badge-success {
  background: var(--success-bg);
  color: var(--success-color);
}

.badge-danger {
  background: var(--danger-bg);
  color: var(--danger-color);
}

.actions-cell {
  white-space: nowrap;
}

.btn-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  border: none;
  border-radius: 0.375rem;
  background: transparent;
  color: var(--text-muted);
  cursor: pointer;
  transition: all 0.15s;
}

.btn-icon:hover {
  background: var(--bg-secondary);
  color: var(--text-primary);
}

.btn-icon:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.btn-icon-success:hover {
  background: rgba(34, 197, 94, 0.1);
  color: var(--success-color);
}

.btn-icon-warning:hover {
  background: rgba(234, 179, 8, 0.1);
  color: var(--warning-color);
}

.btn-icon-info:hover {
  background: rgba(59, 130, 246, 0.1);
  color: #3b82f6;
}

.btn-icon-danger:hover {
  background: rgba(239, 68, 68, 0.1);
  color: var(--danger-color);
}

.loading-row, .empty-row {
  text-align: center;
  padding: 3rem !important;
  color: var(--text-muted);
}

.empty-row {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.loading-spinner {
  width: 2rem;
  height: 2rem;
  border: 3px solid var(--border-color);
  border-top-color: var(--accent-color);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 1rem;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.table-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem 1rem;
  background: var(--bg-card);
  border-radius: 0.5rem;
  border: 1px solid var(--border-color);
  font-size: 0.875rem;
  color: var(--text-muted);
  flex-wrap: wrap;
  gap: 1rem;
}

.stats {
  display: flex;
  gap: 1.5rem;
}

.critical-stat {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: var(--danger-color);
}

/* Pagination */
.pagination {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.page-size-select {
  padding: 0.375rem 0.75rem;
  border: 1px solid var(--border-color);
  border-radius: 0.375rem;
  background: var(--bg-secondary);
  color: var(--text-primary);
  font-size: 0.875rem;
  cursor: pointer;
}

.page-info {
  font-size: 0.875rem;
  color: var(--text-primary);
}

.page-buttons {
  display: flex;
  gap: 0.25rem;
}

.page-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  border: 1px solid var(--border-color);
  border-radius: 0.375rem;
  background: var(--bg-secondary);
  color: var(--text-primary);
  cursor: pointer;
  transition: all 0.15s;
  font-size: 0.875rem;
}

.page-btn:hover:not(:disabled) {
  background: var(--accent-color);
  border-color: var(--accent-color);
  color: white;
}

.page-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

/* Responsive */
@media (max-width: 1024px) {
  .filter-bar {
    flex-direction: column;
    align-items: stretch;
  }
  
  .search-box {
    min-width: 100%;
  }
  
  .filter-controls {
    flex-wrap: wrap;
  }
  
  .action-buttons {
    margin-left: 0;
    justify-content: flex-end;
  }
  
  .table-footer {
    flex-direction: column;
    align-items: stretch;
  }
  
  .pagination {
    justify-content: center;
  }
}
</style>
