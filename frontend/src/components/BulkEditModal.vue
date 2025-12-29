<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-container">
      <div class="modal-header">
        <h2>
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M17 3a2.85 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z"></path>
          </svg>
          Toplu Ürün Düzenleme
        </h2>
        <button @click="$emit('close')" class="close-btn">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M18 6 6 18"></path>
            <path d="m6 6 12 12"></path>
          </svg>
        </button>
      </div>
      
      <div class="modal-body">
        <div class="info-banner">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"></circle>
            <path d="M12 16v-4"></path>
            <path d="M12 8h.01"></path>
          </svg>
          <span><strong>{{ products.length }}</strong> ürün için ortak değerler düzenlenecek. Boş bırakılan alanlar değiştirilmeyecek.</span>
        </div>
        
        <!-- Selected Products List -->
        <div class="selected-products">
          <div class="selected-header" @click="showProductList = !showProductList">
            <span>Seçili Ürünler</span>
            <svg :class="{ 'rotate': showProductList }" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="m6 9 6 6 6-6"/>
            </svg>
          </div>
          <div v-if="showProductList" class="product-list">
            <div v-for="p in products" :key="p.id" class="product-item">
              <span class="product-name">{{ p.name }}</span>
              <span class="product-meta">{{ p.oem_number || '-' }}</span>
            </div>
          </div>
        </div>
        
        <div class="form-section">
          <!-- Category -->
          <div class="form-group">
            <label>
              <input type="checkbox" v-model="updateFields.category" />
              Kategori
            </label>
            <AutocompleteSelect
              v-if="updateFields.category"
              v-model="form.category"
              :items="categoryOptions"
              placeholder="Kategori seçin..."
              :allow-create="true"
              @create="handleNewCategory"
            />
          </div>
          
          <!-- Brand -->
          <div class="form-group">
            <label>
              <input type="checkbox" v-model="updateFields.brand" />
              Marka
            </label>
            <AutocompleteSelect
              v-if="updateFields.brand"
              v-model="form.brand"
              :items="brandOptions"
              placeholder="Marka seçin..."
              :allow-create="true"
              @create="handleNewBrand"
            />
          </div>
          
          <!-- Unit -->
          <div class="form-group">
            <label>
              <input type="checkbox" v-model="updateFields.unit" />
              Birim
            </label>
            <AutocompleteSelect
              v-if="updateFields.unit"
              v-model="form.unit"
              :items="unitOptions"
              placeholder="Birim seçin..."
              :show-search-icon="false"
            />
          </div>
          
          <!-- Critical Stock -->
          <div class="form-group">
            <label>
              <input type="checkbox" v-model="updateFields.critical_stock" />
              Kritik Stok Seviyesi
            </label>
            <input 
              v-if="updateFields.critical_stock"
              type="number" 
              v-model.number="form.critical_stock" 
              class="form-input"
              min="1"
              placeholder="Kritik stok seviyesi"
            />
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
        
        <!-- Preview -->
        <div v-if="hasUpdates" class="preview-section">
          <h4>Yapılacak Değişiklikler</h4>
          <ul class="changes-list">
            <li v-if="updateFields.category && form.category">
              <strong>Kategori:</strong> {{ form.category }}
            </li>
            <li v-if="updateFields.brand && form.brand">
              <strong>Marka:</strong> {{ form.brand }}
            </li>
            <li v-if="updateFields.unit && form.unit">
              <strong>Birim:</strong> {{ unitLabel(form.unit) }}
            </li>
            <li v-if="updateFields.critical_stock && form.critical_stock">
              <strong>Kritik Stok:</strong> {{ form.critical_stock }}
            </li>
          </ul>
        </div>
      </div>
      
      <div class="modal-footer">
        <button @click="$emit('close')" class="btn btn-secondary">İptal</button>
        <button 
          @click="save" 
          class="btn btn-primary"
          :disabled="saving || !hasUpdates"
        >
          <span v-if="saving" class="loading-spinner-sm"></span>
          {{ products.length }} Ürünü Güncelle
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useStock } from '../composables/useStock'
import AutocompleteSelect from './AutocompleteSelect.vue'

const props = defineProps({
  products: {
    type: Array,
    required: true
  }
})

const emit = defineEmits(['close', 'success'])

const { 
  categories, 
  brands,
  units,
  updateProduct,
  loadCategories
} = useStock()

const saving = ref(false)
const error = ref('')
const showProductList = ref(false)

// Which fields to update
const updateFields = ref({
  category: false,
  brand: false,
  unit: false,
  critical_stock: false
})

// Form values
const form = ref({
  category: '',
  brand: '',
  unit: 'adet',
  critical_stock: 3
})

// Get unique brands from products
const existingBrands = computed(() => {
  const brands = new Set()
  props.products.forEach(p => {
    if (p.brand) brands.add(p.brand)
  })
  return Array.from(brands)
})

// Options for autocomplete
const categoryOptions = computed(() => 
  categories.value.map(c => ({ label: c, value: c }))
)

const brandOptions = computed(() => 
  existingBrands.value.map(b => ({ label: b, value: b }))
)

const unitOptions = [
  { label: 'Adet', value: 'adet' },
  { label: 'Litre', value: 'litre' },
  { label: 'Kutu', value: 'kutu' },
  { label: 'Paket', value: 'paket' }
]

// Check if there are any updates to make
const hasUpdates = computed(() => {
  return (updateFields.value.category && form.value.category) ||
         (updateFields.value.brand && form.value.brand) ||
         (updateFields.value.unit && form.value.unit) ||
         (updateFields.value.critical_stock && form.value.critical_stock > 0)
})

// Unit label
function unitLabel(unit) {
  const labels = { adet: 'Adet', litre: 'Litre', kutu: 'Kutu', paket: 'Paket' }
  return labels[unit] || unit
}

// Handle new category creation
function handleNewCategory(name) {
  if (!categories.value.includes(name)) {
    categories.value.unshift(name)
  }
  form.value.category = name
}

// Handle new brand creation
function handleNewBrand(name) {
  if (brands && Array.isArray(brands.value) && !brands.value.includes(name)) {
    brands.value.unshift(name)
  }
  form.value.brand = name
}

// Save
async function save() {
  error.value = ''
  saving.value = true
  
  try {
    let successCount = 0
    let errorCount = 0
    
    for (const product of props.products) {
      const updates = { id: product.id }
      
      // Only include fields that are being updated
      if (updateFields.value.category && form.value.category) {
        updates.category = form.value.category
      }
      if (updateFields.value.brand && form.value.brand) {
        updates.brand = form.value.brand
      }
      if (updateFields.value.unit && form.value.unit) {
        updates.unit = form.value.unit
      }
      if (updateFields.value.critical_stock && form.value.critical_stock > 0) {
        updates.critical_stock = form.value.critical_stock
      }
      
      // Keep existing values for non-updated fields
      updates.name = product.name
      if (!updates.category) updates.category = product.category
      if (!updates.brand) updates.brand = product.brand
      if (!updates.unit) updates.unit = product.unit
      if (!updates.critical_stock) updates.critical_stock = product.critical_stock
      updates.oem_number = product.oem_number
      
      const result = await updateProduct(updates)
      
      if (result.error) {
        errorCount++
      } else {
        successCount++
      }
    }
    
    if (errorCount > 0) {
      error.value = `${errorCount} ürün güncellenemedi`
    }
    
    if (successCount > 0) {
      emit('success', `${successCount} ürün güncellendi`)
      emit('close')
    }
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadCategories()
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
  max-width: 550px;
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
  color: var(--accent-color);
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

.modal-body {
  padding: 1.5rem;
  overflow-y: auto;
  flex: 1;
}

.info-banner {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem;
  background: rgba(59, 130, 246, 0.1);
  border: 1px solid rgba(59, 130, 246, 0.3);
  border-radius: 0.5rem;
  color: #3b82f6;
  font-size: 0.875rem;
  margin-bottom: 1rem;
}

.selected-products {
  margin-bottom: 1.5rem;
  border: 1px solid var(--border-color);
  border-radius: 0.5rem;
  overflow: hidden;
}

.selected-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.75rem 1rem;
  background: var(--bg-secondary);
  cursor: pointer;
  font-weight: 500;
  color: var(--text-primary);
}

.selected-header svg {
  transition: transform 0.2s;
}

.selected-header svg.rotate {
  transform: rotate(180deg);
}

.product-list {
  max-height: 150px;
  overflow-y: auto;
}

.product-item {
  display: flex;
  justify-content: space-between;
  padding: 0.5rem 1rem;
  border-top: 1px solid var(--border-color);
  font-size: 0.875rem;
}

.product-name {
  color: var(--text-primary);
}

.product-meta {
  color: var(--text-muted);
  font-family: monospace;
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
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-primary);
  cursor: pointer;
}

.form-group label input[type="checkbox"] {
  width: 1rem;
  height: 1rem;
  cursor: pointer;
}

.preview-section {
  margin-top: 1.5rem;
  padding: 1rem;
  background: var(--bg-secondary);
  border-radius: 0.5rem;
}

.preview-section h4 {
  margin: 0 0 0.75rem 0;
  font-size: 0.875rem;
  color: var(--text-muted);
}

.changes-list {
  margin: 0;
  padding: 0;
  list-style: none;
}

.changes-list li {
  padding: 0.375rem 0;
  font-size: 0.875rem;
  color: var(--text-primary);
}

.changes-list li strong {
  color: var(--accent-color);
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
