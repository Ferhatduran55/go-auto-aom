<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-container">
      <div class="modal-header">
        <h2>
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M5 12h14"></path>
          </svg>
          Stok Çıkışı
        </h2>
        <button @click="$emit('close')" class="close-btn">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M18 6 6 18"></path>
            <path d="m6 6 12 12"></path>
          </svg>
        </button>
      </div>
      
      <!-- Tab Selection -->
      <div class="tab-bar">
        <button 
          :class="{ active: activeTab === 'single' }" 
          @click="activeTab = 'single'"
        >
          Tekli Çıkış
        </button>
        <button 
          :class="{ active: activeTab === 'bulk' }" 
          @click="activeTab = 'bulk'"
        >
          Toplu Çıkış
        </button>
      </div>
      
      <div class="modal-body">
        <!-- Single Exit Form -->
        <div v-if="activeTab === 'single'" class="form-section">
          <div class="form-group">
            <label>Ürün *</label>
            <div class="autocomplete-wrapper" ref="autocompleteRef">
              <div class="autocomplete-input-container">
                <svg class="search-icon" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="11" cy="11" r="8"></circle>
                  <path d="m21 21-4.35-4.35"></path>
                </svg>
                <input 
                  type="text" 
                  v-model="searchQuery"
                  @input="onSearchInput"
                  @focus="showDropdown = true"
                  @keydown.down.prevent="navigateDown"
                  @keydown.up.prevent="navigateUp"
                  @keydown.enter.prevent="selectHighlighted"
                  @keydown.escape="showDropdown = false"
                  class="form-input autocomplete-input"
                  :placeholder="product ? '' : 'Ürün adı veya OEM ile ara...'"
                  :disabled="product !== null"
                />
                <button 
                  v-if="selectedProductDetail && !product" 
                  @click="clearSelection" 
                  class="clear-btn"
                  type="button"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M18 6 6 18"></path>
                    <path d="m6 6 12 12"></path>
                  </svg>
                </button>
              </div>
              
              <!-- Selected Product Display -->
              <div v-if="selectedProductDetail" class="selected-product" :class="{ 'low-stock': selectedProductDetail.stock_quantity <= 0 }">
                <div class="product-info">
                  <span class="product-name">{{ selectedProductDetail.name }}</span>
                  <span class="product-meta">
                    {{ selectedProductDetail.oem_number ? `OEM: ${selectedProductDetail.oem_number}` : '' }}
                    {{ selectedProductDetail.brand ? `• ${selectedProductDetail.brand}` : '' }}
                  </span>
                </div>
                <span class="stock-badge" :class="{ 'stock-low': selectedProductDetail.stock_quantity <= (selectedProductDetail.critical_stock || 3) }">
                  Stok: {{ formatQuantity(selectedProductDetail.stock_quantity, selectedProductDetail.unit) }} {{ selectedProductDetail.unit }}
                </span>
              </div>
              
              <!-- Dropdown -->
              <div v-if="showDropdown && filteredProducts.length > 0 && !selectedProductDetail" class="autocomplete-dropdown">
                <div 
                  v-for="(p, index) in filteredProducts" 
                  :key="p.id"
                  @click="selectProduct(p)"
                  @mouseenter="highlightedIndex = index"
                  :class="['dropdown-item', { highlighted: index === highlightedIndex, disabled: p.stock_quantity <= 0 }]"
                >
                  <div class="item-main">
                    <span class="item-name">{{ p.name }}</span>
                    <span class="item-stock" :class="{ 'stock-empty': p.stock_quantity <= 0, 'stock-low': p.stock_quantity <= (p.critical_stock || 3) }">
                      {{ formatQuantity(p.stock_quantity, p.unit) }} {{ p.unit }}
                    </span>
                  </div>
                  <div class="item-meta">
                    <span v-if="p.oem_number">OEM: {{ p.oem_number }}</span>
                    <span v-if="p.brand">{{ p.brand }}</span>
                    <span v-if="p.stock_quantity <= 0" class="item-warning">Stok yok!</span>
                  </div>
                </div>
              </div>
              
              <!-- No Results -->
              <div v-if="showDropdown && searchQuery.length >= 2 && filteredProducts.length === 0" class="autocomplete-dropdown">
                <div class="no-results">Sonuç bulunamadı</div>
              </div>
            </div>
          </div>
          
          <div class="form-row">
            <div class="form-group">
              <label>Miktar *</label>
              <input 
                type="number" 
                v-model.number="singleForm.amount" 
                :step="selectedProductDetail?.unit === 'litre' ? '0.1' : '1'"
                :min="selectedProductDetail?.unit === 'litre' ? '0.1' : '1'"
                :max="selectedProductDetail?.stock_quantity || 0"
                class="form-input"
                :class="{ 'input-error': amountError }"
                placeholder="Çıkış miktarı"
              />
              <span v-if="selectedProductDetail" class="input-hint">
                Maksimum: {{ formatQuantity(selectedProductDetail.stock_quantity, selectedProductDetail.unit) }} {{ selectedProductDetail.unit }}
              </span>
            </div>
            <div class="form-group unit-display">
              <label>Birim</label>
              <div class="unit-badge">{{ selectedProductDetail?.unit || 'adet' }}</div>
            </div>
          </div>
          
          <div class="form-group">
            <label>Açıklama / Plaka / Araç Bilgisi</label>
            <textarea 
              v-model="singleForm.note" 
              class="form-input"
              rows="3"
              placeholder="Örn: 34 ABC 123 - 2015 Opel Astra bakım"
            ></textarea>
          </div>
          
          <!-- Preview -->
          <div v-if="selectedProductDetail && singleForm.amount > 0" class="preview-box" :class="{ 'preview-warning': amountError }">
            <div class="preview-row">
              <span>Mevcut Stok:</span>
              <strong>{{ formatQuantity(selectedProductDetail.stock_quantity, selectedProductDetail.unit) }} {{ selectedProductDetail.unit }}</strong>
            </div>
            <div class="preview-row">
              <span>Çıkış:</span>
              <strong class="text-danger">-{{ formatQuantity(singleForm.amount, selectedProductDetail.unit) }} {{ selectedProductDetail.unit }}</strong>
            </div>
            <div class="preview-row total">
              <span>Kalan Stok:</span>
              <strong :class="{ 'text-danger': newStock < 0, 'text-warning': newStock < (selectedProductDetail.critical_stock || 3) }">
                {{ formatQuantity(newStock, selectedProductDetail.unit) }} {{ selectedProductDetail.unit }}
              </strong>
            </div>
            
            <!-- Warnings -->
            <div v-if="amountError" class="warning-box danger">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"></circle>
                <path d="m15 9-6 6"></path>
                <path d="m9 9 6 6"></path>
              </svg>
              Yetersiz stok! Mevcut stoktan fazla çıkış yapamazsınız.
            </div>
            <div v-else-if="newStock < (selectedProductDetail?.critical_stock || 3)" class="warning-box warning">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"></path>
                <path d="M12 9v4"></path>
                <path d="M12 17h.01"></path>
              </svg>
              Bu çıkıştan sonra stok kritik seviyenin altına düşecek!
            </div>
          </div>
        </div>
        
        <!-- Bulk Exit Form -->
        <div v-else class="form-section">
          <div class="bulk-header">
            <span>{{ bulkExits.filter(e => e.productId).length }} ürün seçildi</span>
            <button @click="addBulkRow" class="btn btn-sm btn-secondary">
              <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M12 5v14"></path>
                <path d="M5 12h14"></path>
              </svg>
              Satır Ekle
            </button>
          </div>
          
          <div class="form-group">
            <label>Ortak Açıklama (Tüm çıkışlar için)</label>
            <input 
              type="text" 
              v-model="sharedNote" 
              class="form-input"
              placeholder="Örn: 34 ABC 123 - Periyodik bakım"
            />
          </div>
          
          <div class="bulk-list">
            <div v-for="(exit, index) in bulkExits" :key="index" class="bulk-row">
              <div class="bulk-product-select">
                <AutocompleteSelect
                  v-model="exit.productId"
                  :items="productOptionsWithStock"
                  placeholder="Ürün seçin..."
                  meta-key="meta"
                />
              </div>
              
              <input 
                type="number" 
                v-model.number="exit.amount" 
                :step="getProductUnit(exit.productId) === 'litre' ? '0.1' : '1'"
                :min="0.1"
                :max="getProductStock(exit.productId)"
                class="form-input amount-col"
                placeholder="Miktar"
              />
              
              <button @click="removeBulkRow(index)" class="btn-icon btn-icon-danger">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M3 6h18"></path>
                  <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"></path>
                </svg>
              </button>
            </div>
          </div>
          
          <div v-if="bulkExits.length === 0" class="empty-state">
            <p>Toplu çıkış için satır ekleyin</p>
          </div>
        </div>
        
        <!-- Error Message -->
        <div v-if="error" class="error-message">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"></circle>
            <path d="m15 9-6 6"></path>
            <path d="m9 9 6 6"></path>
          </svg>
          {{ error }}
        </div>
      </div>
      
      <div class="modal-footer">
        <button @click="$emit('close')" class="btn btn-secondary">İptal</button>
        <button 
          @click="save" 
          class="btn btn-warning"
          :disabled="saving || !formValid || amountError"
        >
          <span v-if="saving" class="loading-spinner-sm"></span>
          {{ activeTab === 'single' ? 'Stok Çıkışı Yap' : 'Toplu Çıkış Yap' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useStock } from '../composables/useStock'
import AutocompleteSelect from './AutocompleteSelect.vue'

const props = defineProps({
  product: {
    type: Object,
    default: null
  },
  bulkProducts: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['close', 'success'])

const { 
  products, 
  handleStockOut, 
  bulkStockOut,
  formatQuantity 
} = useStock()

const activeTab = ref('single')
const saving = ref(false)
const error = ref('')
const sharedNote = ref('')
const autocompleteRef = ref(null)

// Autocomplete state
const searchQuery = ref('')
const showDropdown = ref(false)
const highlightedIndex = ref(0)

// Single exit form
const singleForm = ref({
  productId: '',
  amount: 1,
  note: ''
})

// Bulk exits list
const bulkExits = ref([])

// Product options with stock for autocomplete
const productOptionsWithStock = computed(() => 
  products.value
    .filter(p => p.stock_quantity > 0)
    .map(p => ({
      label: p.name,
      value: p.id,
      meta: `${formatQuantity(p.stock_quantity, p.unit)} ${p.unit}`
    }))
)

// Filtered products for autocomplete (only with stock > 0)
const filteredProducts = computed(() => {
  const availableProducts = products.value.filter(p => p.stock_quantity > 0)
  
  if (!searchQuery.value || searchQuery.value.length < 1) {
    return availableProducts.slice(0, 20)
  }
  
  const query = searchQuery.value.toLowerCase()
  return availableProducts.filter(p => 
    p.name?.toLowerCase().includes(query) ||
    p.oem_number?.toLowerCase().includes(query) ||
    p.brand?.toLowerCase().includes(query)
  ).slice(0, 20)
})

// Selected product detail
const selectedProductDetail = computed(() => {
  if (!singleForm.value.productId) return null
  return products.value.find(p => p.id === singleForm.value.productId)
})

// New stock after exit
const newStock = computed(() => {
  if (!selectedProductDetail.value) return 0
  return selectedProductDetail.value.stock_quantity - singleForm.value.amount
})

// Amount error check
const amountError = computed(() => {
  if (!selectedProductDetail.value) return false
  return singleForm.value.amount > selectedProductDetail.value.stock_quantity
})

// Form validity
const formValid = computed(() => {
  if (activeTab.value === 'single') {
    return singleForm.value.productId && singleForm.value.amount > 0
  } else {
    return bulkExits.value.some(e => e.productId && e.amount > 0)
  }
})

// Autocomplete methods
function onSearchInput() {
  showDropdown.value = true
  highlightedIndex.value = 0
}

function selectProduct(product) {
  if (product.stock_quantity <= 0) return
  
  singleForm.value.productId = product.id
  searchQuery.value = product.name
  showDropdown.value = false
}

function clearSelection() {
  singleForm.value.productId = ''
  searchQuery.value = ''
  showDropdown.value = false
}

function navigateDown() {
  if (highlightedIndex.value < filteredProducts.value.length - 1) {
    highlightedIndex.value++
  }
}

function navigateUp() {
  if (highlightedIndex.value > 0) {
    highlightedIndex.value--
  }
}

function selectHighlighted() {
  if (filteredProducts.value[highlightedIndex.value]) {
    selectProduct(filteredProducts.value[highlightedIndex.value])
  }
}

// Close dropdown on outside click
function handleClickOutside(e) {
  if (autocompleteRef.value && !autocompleteRef.value.contains(e.target)) {
    showDropdown.value = false
  }
}

// Get product unit
const getProductUnit = (productId) => {
  const prod = products.value.find(p => p.id === productId)
  return prod?.unit || 'adet'
}

// Get product stock
const getProductStock = (productId) => {
  const prod = products.value.find(p => p.id === productId)
  return prod?.stock_quantity || 0
}

// Add bulk row
const addBulkRow = () => {
  bulkExits.value.push({
    productId: '',
    amount: 1
  })
}

// Remove bulk row
const removeBulkRow = (index) => {
  bulkExits.value.splice(index, 1)
}

// Save
const save = async () => {
  error.value = ''
  saving.value = true
  
  try {
    if (activeTab.value === 'single') {
      const result = await handleStockOut(
        singleForm.value.productId,
        singleForm.value.amount,
        singleForm.value.note
      )
      
      if (result.error) {
        error.value = result.error
      } else {
        emit('success', 'Stok çıkışı başarıyla tamamlandı')
        emit('close')
      }
    } else {
      // Filter valid exits
      const validExits = bulkExits.value
        .filter(e => e.productId && e.amount > 0)
        .map(e => ({
          product_id: e.productId,
          amount: e.amount,
          note: sharedNote.value
        }))
      
      if (validExits.length === 0) {
        error.value = 'En az bir geçerli çıkış yapmalısınız'
        return
      }
      
      const result = await bulkStockOut(validExits)
      
      if (result.error) {
        error.value = result.error
      } else {
        const message = `${result.successful} üründen stok çıkışı yapıldı`
        emit('success', message)
        emit('close')
      }
    }
  } finally {
    saving.value = false
  }
}

// Update form when product prop changes
watch(() => props.product, (newProduct) => {
  if (newProduct) {
    singleForm.value.productId = newProduct.id
    searchQuery.value = newProduct.name
  }
}, { immediate: true })

// Handle bulk products from selection
watch(() => props.bulkProducts, (newProducts) => {
  if (newProducts && newProducts.length > 0) {
    activeTab.value = 'bulk'
    bulkExits.value = newProducts.map(p => ({
      productId: p.id,
      amount: 1
    }))
  }
}, { immediate: true })

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  
  // Add initial row for bulk if empty
  if (bulkExits.value.length === 0) {
    addBulkRow()
  }
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  backdrop-filter: blur(2px);
}

.modal-container {
  background: var(--bg-card);
  border-radius: 0.75rem;
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.25rem 1.5rem;
  border-bottom: 1px solid var(--border-color);
}

.modal-header h2 {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin: 0;
  font-size: 1.25rem;
  color: var(--text-primary);
}

.modal-header h2 svg {
  color: #eab308;
}

.close-btn {
  display: flex;
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

.close-btn:hover {
  background: var(--bg-secondary);
  color: var(--text-primary);
}

.tab-bar {
  display: flex;
  border-bottom: 1px solid var(--border-color);
}

.tab-bar button {
  flex: 1;
  padding: 0.75rem 1rem;
  border: none;
  background: transparent;
  color: var(--text-muted);
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s;
  border-bottom: 2px solid transparent;
  margin-bottom: -1px;
}

.tab-bar button:hover {
  color: var(--text-primary);
  background: var(--bg-secondary);
}

.tab-bar button.active {
  color: #eab308;
  border-bottom-color: #eab308;
}

.modal-body {
  padding: 1.5rem;
  overflow-y: auto;
  flex: 1;
}

.form-section {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-group label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-primary);
}

.form-row {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 1rem;
}

/* Autocomplete Styles */
.autocomplete-wrapper {
  position: relative;
}

.autocomplete-input-container {
  position: relative;
  display: flex;
  align-items: center;
}

.search-icon {
  position: absolute;
  left: 0.75rem;
  color: var(--text-muted);
  pointer-events: none;
  z-index: 1;
}

.autocomplete-input {
  padding-left: 2.5rem !important;
  padding-right: 2.5rem !important;
}

.clear-btn {
  position: absolute;
  right: 0.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 1.5rem;
  height: 1.5rem;
  border: none;
  border-radius: 50%;
  background: var(--bg-secondary);
  color: var(--text-muted);
  cursor: pointer;
  transition: all 0.15s;
}

.clear-btn:hover {
  background: var(--danger-color);
  color: white;
}

.selected-product {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.75rem;
  margin-top: 0.5rem;
  background: var(--bg-secondary);
  border: 2px solid var(--accent-color);
  border-radius: 0.5rem;
}

.selected-product.low-stock {
  border-color: #ef4444;
}

.product-info {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.product-name {
  font-weight: 600;
  color: var(--text-primary);
}

.product-meta {
  font-size: 0.75rem;
  color: var(--text-muted);
}

.stock-badge {
  padding: 0.25rem 0.75rem;
  background: var(--accent-color);
  color: white;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 500;
}

.stock-badge.stock-low {
  background: #eab308;
}

.autocomplete-dropdown {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  margin-top: 0.25rem;
  background: var(--bg-card);
  border: 2px solid var(--border-color);
  border-radius: 0.75rem;
  max-height: 300px;
  overflow-y: auto;
  z-index: 100;
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.2);
}

.dropdown-item {
  padding: 0.75rem 1rem;
  cursor: pointer;
  border-bottom: 1px solid var(--border-color);
  transition: all 0.15s;
}

.dropdown-item:last-child {
  border-bottom: none;
}

.dropdown-item:hover,
.dropdown-item.highlighted {
  background: var(--accent-color);
  color: white;
}

.dropdown-item.disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.dropdown-item.highlighted .item-meta,
.dropdown-item:hover .item-meta {
  color: rgba(255, 255, 255, 0.8);
}

.item-main {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.item-name {
  font-weight: 600;
}

.item-stock {
  font-size: 0.75rem;
  padding: 0.125rem 0.5rem;
  background: rgba(34, 197, 94, 0.2);
  border-radius: 9999px;
  color: #22c55e;
}

.item-stock.stock-low {
  background: rgba(234, 179, 8, 0.2);
  color: #eab308;
}

.item-stock.stock-empty {
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
}

.item-meta {
  display: flex;
  gap: 0.75rem;
  font-size: 0.75rem;
  color: var(--text-muted);
  margin-top: 0.25rem;
}

.item-warning {
  color: #ef4444;
  font-weight: 500;
}

.no-results {
  padding: 1rem;
  text-align: center;
  color: var(--text-muted);
}

.input-hint {
  font-size: 0.75rem;
  color: var(--text-muted);
}

.input-error {
  border-color: #ef4444 !important;
}

.unit-display {
  min-width: 80px;
}

.unit-badge {
  padding: 0.5rem 1rem;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 0.375rem;
  text-align: center;
  font-weight: 500;
  color: var(--text-primary);
}

.preview-box {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 0.5rem;
  padding: 1rem;
  margin-top: 0.5rem;
}

.preview-box.preview-warning {
  border-color: #ef4444;
  background: rgba(239, 68, 68, 0.05);
}

.preview-row {
  display: flex;
  justify-content: space-between;
  padding: 0.5rem 0;
  color: var(--text-muted);
}

.preview-row.total {
  border-top: 1px solid var(--border-color);
  margin-top: 0.5rem;
  padding-top: 1rem;
  font-size: 1.1rem;
}

.text-danger {
  color: #ef4444;
}

.text-warning {
  color: #eab308;
}

.warning-box {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  margin-top: 0.75rem;
}

.warning-box.danger {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.warning-box.warning {
  background: rgba(234, 179, 8, 0.1);
  color: #eab308;
  border: 1px solid rgba(234, 179, 8, 0.3);
}

/* Bulk Exit */
.bulk-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid var(--border-color);
}

.bulk-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  max-height: 250px;
  overflow-y: auto;
  padding-right: 0.5rem;
}

.bulk-row {
  display: grid;
  grid-template-columns: 1fr 80px 32px;
  gap: 0.5rem;
  align-items: center;
}

.bulk-product-select {
  min-width: 0;
}

.amount-col {
  width: 80px;
}

.empty-state {
  text-align: center;
  padding: 2rem;
  color: var(--text-muted);
}

.error-message {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1rem;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 0.375rem;
  color: #ef4444;
  font-size: 0.875rem;
  margin-top: 1rem;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  padding: 1.25rem 1.5rem;
  border-top: 1px solid var(--border-color);
}

.btn-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  border: none;
  border-radius: 0.375rem;
  background: transparent;
  color: var(--text-muted);
  cursor: pointer;
}

.btn-icon-danger:hover {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
}

.loading-spinner-sm {
  width: 1rem;
  height: 1rem;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
  margin-right: 0.5rem;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
