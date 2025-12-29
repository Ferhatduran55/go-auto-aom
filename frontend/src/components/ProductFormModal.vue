<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-container">
      <div class="modal-header">
        <h2>
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="m7.5 4.27 9 5.15"></path>
            <path d="M21 8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16Z"></path>
            <path d="m3.3 7 8.7 5 8.7-5"></path>
            <path d="M12 22V12"></path>
          </svg>
          {{ editMode ? 'Ürün Düzenle' : 'Yeni Ürün' }}
        </h2>
        <button @click="$emit('close')" class="close-btn">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M18 6 6 18"></path>
            <path d="m6 6 12 12"></path>
          </svg>
        </button>
      </div>
      
      <div class="modal-body">
        <div class="form-section">
          <!-- Ürün Adı -->
          <div class="form-group">
            <label>Ürün Adı *</label>
            <input 
              type="text" 
              v-model="form.name" 
              class="form-input"
              placeholder="Örn: Motor Yağı 5W-30"
              ref="firstInput"
            />
          </div>
          
          <!-- OEM No ve Marka -->
          <div class="form-row">
            <div class="form-group">
              <label>OEM Numarası</label>
              <input 
                type="text" 
                v-model="form.oem_number" 
                class="form-input"
                placeholder="Örn: 1520865F0A"
              />
            </div>
            <div class="form-group">
              <label>Marka</label>
              <AutocompleteSelect
                v-model="form.brand"
                :items="brandOptions"
                placeholder="Marka seçin veya yazın..."
                :allow-create="true"
                @create="handleNewBrand"
              />
            </div>
          </div>
          
          <!-- Category -->
          <div class="form-group">
            <label>Kategori</label>
            <AutocompleteSelect
              v-model="form.category"
              :items="categoryOptions"
              placeholder="Kategori seçin veya yazın..."
              :allow-create="true"
              @create="handleNewCategory"
            />
          </div>
          
          <!-- Unit and Critical Stock -->
          <div class="form-row">
            <div class="form-group">
              <label>Birim *</label>
              <AutocompleteSelect
                v-model="form.unit"
                :items="unitOptions"
                placeholder="Birim seçin..."
              />
            </div>
            <div class="form-group">
              <label>Kritik Stok Seviyesi</label>
              <input 
                type="number" 
                v-model.number="form.critical_stock" 
                class="form-input"
                min="1"
                placeholder="Varsayılan: 3"
              />
            </div>
          </div>
          
          <!-- Initial Stock (only for new product) -->
          <div v-if="!editMode" class="form-group">
            <label>Başlangıç Stoku</label>
            <div class="stock-input-wrapper">
              <input 
                type="number" 
                v-model.number="form.stock_quantity" 
                class="form-input"
                :step="form.unit === 'litre' ? '0.1' : '1'"
                min="0"
                placeholder="0"
              />
              <span class="unit-suffix">{{ form.unit }}</span>
            </div>
            <span class="input-hint">Stok girişi daha sonra da yapılabilir</span>
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
          class="btn btn-primary"
          :disabled="saving || !formValid"
        >
          <span v-if="saving" class="loading-spinner-sm"></span>
          {{ editMode ? 'Güncelle' : 'Ürün Ekle' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useStock } from '../composables/useStock'
import AutocompleteSelect from './AutocompleteSelect.vue'

const props = defineProps({
  product: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['close', 'success'])

const { 
  categories,
  brands,
  units, 
  createProduct,
  updateProduct,
  loadCategories,
  loadBrands
} = useStock()

const firstInput = ref(null)
const saving = ref(false)
const error = ref('')

// Edit mode check
const editMode = computed(() => props.product !== null)

// Form
const form = ref({
  name: '',
  oem_number: '',
  brand: '',
  category: '',
  unit: 'adet',
  critical_stock: 3,
  stock_quantity: 0
})

// Category and unit options for autocomplete
const categoryOptions = computed(() => 
  categories.value.map(c => ({ label: c, value: c }))
)

// Brand options for autocomplete
const brandOptions = computed(() => 
  brands.value.map(b => ({ label: b, value: b }))
)

const unitOptions = [
  { label: 'Adet', value: 'adet' },
  { label: 'Litre', value: 'litre' },
  { label: 'Kutu', value: 'kutu' },
  { label: 'Paket', value: 'paket' }
]

// Handle new category creation
function handleNewCategory(name) {
  if (!categories.value.includes(name)) {
    categories.value.unshift(name)
  }
  form.value.category = name
}

// Handle new brand creation
function handleNewBrand(name) {
  if (!brands.value.includes(name)) {
    brands.value.unshift(name)
  }
  form.value.brand = name
}

// Unit label
const unitLabel = (unit) => {
  const labels = {
    adet: 'Adet',
    litre: 'Litre',
    kutu: 'Kutu',
    paket: 'Paket'
  }
  return labels[unit] || unit
}

// Form validity
const formValid = computed(() => {
  return form.value.name.trim() !== ''
})

// Save
const save = async () => {
  error.value = ''
  saving.value = true
  
  try {
    let result
    
    if (editMode.value) {
      // Update
      result = await updateProduct({
        id: props.product.id,
        name: form.value.name,
        oem_number: form.value.oem_number,
        brand: form.value.brand,
        category: form.value.category,
        unit: form.value.unit,
        critical_stock: form.value.critical_stock
      })
    } else {
      // Create new product
      result = await createProduct({
        name: form.value.name,
        oem_number: form.value.oem_number,
        brand: form.value.brand,
        category: form.value.category,
        unit: form.value.unit,
        stock_quantity: form.value.stock_quantity,
        critical_stock: form.value.critical_stock
      })
    }
    
    if (result.error) {
      error.value = result.error
    } else {
      emit('success', editMode.value ? 'Ürün güncellendi' : 'Ürün eklendi')
      emit('close')
    }
  } finally {
    saving.value = false
  }
}

// Fill form when product changes
watch(() => props.product, (newProduct) => {
  if (newProduct) {
    form.value = {
      name: newProduct.name || '',
      oem_number: newProduct.oem_number || '',
      brand: newProduct.brand || '',
      category: newProduct.category || '',
      unit: newProduct.unit || 'adet',
      critical_stock: newProduct.critical_stock || 3,
      stock_quantity: newProduct.stock_quantity || 0
    }
  }
}, { immediate: true })

onMounted(async () => {
  await Promise.all([loadCategories(), loadBrands()])
  
  // Focus first input
  setTimeout(() => {
    firstInput.value?.focus()
  }, 100)
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
  max-width: 500px;
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
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.category-wrapper {
  display: flex;
  flex-direction: column;
}

.mt-2 {
  margin-top: 0.5rem;
}

.stock-input-wrapper {
  position: relative;
}

.stock-input-wrapper input {
  padding-right: 4rem;
}

.unit-suffix {
  position: absolute;
  right: 0.75rem;
  top: 50%;
  transform: translateY(-50%);
  color: var(--text-muted);
  font-size: 0.875rem;
}

.input-hint {
  font-size: 0.75rem;
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
