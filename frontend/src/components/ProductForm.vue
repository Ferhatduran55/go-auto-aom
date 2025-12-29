<template>
  <div class="card">
    <!-- Edit Banner -->
    <div v-if="editingProductId" class="bg-accent text-white px-5 py-3 flex items-center justify-between font-semibold">
      ✏️ Ürün Düzenleniyor
      <button @click="cancelEdit" class="btn btn-sm bg-white/20">İptal</button>
    </div>
    
    <!-- Header -->
    <div class="px-6 py-5 border-b flex items-center gap-4 bg-gradient-to-r from-accent/10 to-purple/10" style="border-color: var(--border-color);">
      <div class="w-12 h-12 bg-gradient-to-r from-accent to-purple rounded-xl flex items-center justify-center text-xl text-white">
        ➕
      </div>
      <div>
        <h3 class="text-lg font-bold" style="color: var(--text-primary);">{{ editingProductId ? 'Ürünü Düzenle' : 'Ürün Ekle' }}</h3>
        <span class="text-sm" style="color: var(--text-muted);">Sipariş kalemi ekleyin</span>
      </div>
    </div>

    <!-- Form Body -->
    <div class="p-6 flex-1 overflow-y-auto">
      <form @submit.prevent="handleSubmit">
        <!-- Ürün Adı -->
        <div class="mb-4 relative">
          <label class="block mb-2 text-sm font-semibold" style="color: var(--text-muted);">Ürün Adı</label>
          <input 
            type="text" 
            v-model="form.name" 
            @input="onProductInput"
            @focus="showProductDropdown = true"
            class="form-input" 
            placeholder="Ürün adı yazın..."
            autocomplete="off"
            required
          >
          <div v-if="showProductDropdown && filteredProducts.length > 0" 
               class="absolute top-full left-0 right-0 mt-1 rounded-xl max-h-64 overflow-y-auto z-50 shadow-xl" style="background: var(--bg-card); border: 2px solid var(--border-color);">
            <div 
              v-for="p in filteredProducts" 
              :key="p.id"
              @click="selectProduct(p)"
              class="px-4 py-3 cursor-pointer border-b last:border-b-0 hover:bg-accent hover:text-white transition-colors" style="border-color: var(--border-color);"
            >
              <div class="font-semibold" style="color: var(--text-primary);">🔧 {{ p.name }}</div>
              <div class="text-sm" style="color: var(--text-muted);">OEM: {{ p.oem_number }} • {{ p.used_count || 0 }}x kullanıldı</div>
            </div>
          </div>
        </div>

        <!-- OEM Numarası -->
        <div class="mb-4 relative">
          <label class="block mb-2 text-sm font-semibold" style="color: var(--text-muted);">OEM Numarası</label>
          <input 
            type="text" 
            v-model="form.oem" 
            @input="onOEMInput"
            @focus="showOEMDropdown = true"
            class="form-input" 
            placeholder="OEM kodu"
            autocomplete="off"
          >
          <div v-if="showOEMDropdown && filteredOEMs.length > 0" 
               class="absolute top-full left-0 right-0 mt-1 rounded-xl max-h-64 overflow-y-auto z-50 shadow-xl" style="background: var(--bg-card); border: 2px solid var(--border-color);">
            <div 
              v-for="p in filteredOEMs" 
              :key="p.id"
              @click="selectOEM(p)"
              class="px-4 py-3 cursor-pointer border-b last:border-b-0 hover:bg-accent hover:text-white transition-colors" style="border-color: var(--border-color);"
            >
              <div class="font-semibold" style="color: var(--text-primary);">🔢 {{ p.oem_number }}</div>
              <div class="text-sm" style="color: var(--text-muted);">Ürün: {{ p.name }}</div>
            </div>
          </div>
        </div>

        <!-- Adet ve Birim Fiyat -->
        <div class="grid grid-cols-2 gap-4 mb-4">
          <div>
            <label class="block mb-2 text-sm font-semibold" style="color: var(--text-muted);">Adet</label>
            <input type="number" v-model.number="form.quantity" min="1" class="form-input" required>
          </div>
          <div>
            <label class="block mb-2 text-sm font-semibold" style="color: var(--text-muted);">Birim Fiyat (₺)</label>
            <input type="number" v-model.number="form.unitPrice" min="0" step="0.01" class="form-input" placeholder="0.00" required>
          </div>
        </div>

        <!-- Parça Durumu -->
        <div class="mb-4">
          <label class="block mb-2 text-sm font-semibold" style="color: var(--text-muted);">Parça Durumu</label>
          <AutocompleteSelect
            v-model="form.partStatus"
            :items="partStatusOptions"
            placeholder="Durum seçin..."
          />
        </div>

        <button type="submit" class="btn btn-primary btn-block">
          {{ editingProductId ? '✓ Güncelle' : '➕ Listeye Ekle' }}
        </button>
      </form>

      <div class="grid grid-cols-2 gap-3 mt-4">
        <button @click="clearForm" class="btn btn-secondary">🔄 Temizle</button>
        <button @click="$emit('clearAll')" class="btn btn-secondary">🗑️ Listeyi Sil</button>
      </div>
    </div>
  </div>

  <!-- OEM Conflict Modal -->
  <Teleport to="body">
    <div v-if="showConflictModal" class="fixed inset-0 bg-black/70 backdrop-blur-sm flex items-center justify-center z-[1000]">
      <div class="bg-card rounded-2xl w-[90%] max-w-[500px] border border-border overflow-hidden animate-modal">
        <div class="p-8 text-center bg-gradient-to-r from-warning/20 to-danger/20">
          <div class="text-5xl mb-4">⚠️</div>
          <div class="text-xl font-bold mb-2 text-slate-200">OEM Çakışması Tespit Edildi</div>
          <div class="text-slate-400">Bu OEM numarası farklı bir ürün adıyla kayıtlı</div>
        </div>
        <div class="p-6">
          <div class="bg-secondary p-4 rounded-xl mb-4">
            <div class="text-sm text-slate-400 mb-2">OEM Numarası</div>
            <div class="text-xl font-bold text-warning">{{ conflictData.oem }}</div>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div class="bg-secondary p-4 rounded-xl border-2 border-success">
              <div class="text-xs text-slate-400 mb-1">📦 Kayıtlı Ürün</div>
              <div class="font-semibold text-success">{{ conflictData.existingName }}</div>
            </div>
            <div class="bg-secondary p-4 rounded-xl border-2 border-accent">
              <div class="text-xs text-slate-400 mb-1">🆕 Yeni Girilen</div>
              <div class="font-semibold text-accent">{{ conflictData.newName }}</div>
            </div>
          </div>
          <div class="mt-4 p-3 bg-accent/10 rounded-lg text-sm text-slate-400">
            💡 Aynı OEM farklı ürünlerde kullanılabilir. Yeni isimle kaydetmek isterseniz "Yeni Oluştur" seçin.
          </div>
        </div>
        <div class="grid grid-cols-3 gap-3 p-5 bg-secondary">
          <button @click="showConflictModal = false" class="btn btn-secondary">İptal</button>
          <button @click="useExisting" class="btn btn-success">Mevcut Kullan</button>
          <button @click="createNew" class="btn btn-primary">Yeni Oluştur</button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted, watch } from 'vue'
import { useOrder } from '@/composables/useOrder'
import { useToast } from '@/composables/useToast'
import AutocompleteSelect from './AutocompleteSelect.vue'

const emit = defineEmits(['clearAll'])

const { allProducts, editingProductId, products, addProduct, updateProduct, checkOEMConflict } = useOrder()
const { showToast } = useToast()

const form = reactive({
  name: '',
  oem: '',
  quantity: 1,
  unitPrice: '',
  partStatus: 'original'
})

// Part status options for autocomplete
const partStatusOptions = [
  { label: '🟢 Orijinal', value: 'original' },
  { label: '🟡 Çıkma', value: 'used' },
  { label: '🔵 Sıfır', value: 'zero' }
]

const showProductDropdown = ref(false)
const showOEMDropdown = ref(false)
const showConflictModal = ref(false)
const conflictData = reactive({
  oem: '',
  existingName: '',
  newName: ''
})

const filteredProducts = computed(() => {
  if (!form.name || form.name.length < 2) return []
  const search = form.name.toLowerCase()
  return allProducts.value
    .filter(p => p.name.toLowerCase().includes(search) || p.oem_number?.toLowerCase().includes(search))
    .slice(0, 6)
})

const filteredOEMs = computed(() => {
  if (!form.oem || form.oem.length < 2) return []
  const search = form.oem.toLowerCase()
  return allProducts.value
    .filter(p => p.oem_number?.toLowerCase().includes(search))
    .slice(0, 6)
})

function onProductInput() {
  showProductDropdown.value = true
}

function onOEMInput() {
  showOEMDropdown.value = true
}

function selectProduct(p) {
  form.name = p.name
  form.oem = p.oem_number || ''
  showProductDropdown.value = false
}

function selectOEM(p) {
  form.name = p.name
  form.oem = p.oem_number
  showOEMDropdown.value = false
}

function handleSubmit() {
  if (!form.name || form.quantity < 1 || !form.unitPrice) {
    showToast('Geçerli değerler girin', 'error')
    return
  }

  if (editingProductId.value) {
    updateProduct(editingProductId.value, {
      name: form.name,
      oem: form.oem,
      quantity: form.quantity,
      unitPrice: parseFloat(form.unitPrice),
      partStatus: form.partStatus
    })
    showToast('Ürün güncellendi')
    clearForm()
    return
  }

  // OEM çakışma kontrolü
  const conflict = checkOEMConflict(form.name, form.oem)
  if (conflict) {
    conflictData.oem = form.oem
    conflictData.existingName = conflict.name
    conflictData.newName = form.name
    showConflictModal.value = true
    return
  }

  addProduct({
    name: form.name,
    oem: form.oem,
    quantity: form.quantity,
    unitPrice: parseFloat(form.unitPrice),
    partStatus: form.partStatus
  })
  showToast('Ürün eklendi')
  clearForm()
}

function useExisting() {
  form.name = conflictData.existingName
  showConflictModal.value = false
  showToast('Mevcut ürün adı kullanıldı')
}

function createNew() {
  addProduct({
    name: form.name,
    oem: form.oem,
    quantity: form.quantity,
    unitPrice: parseFloat(form.unitPrice),
    partStatus: form.partStatus
  })
  showConflictModal.value = false
  showToast('Farklı isimle yeni ürün oluşturuldu')
  clearForm()
}

function clearForm() {
  form.name = ''
  form.oem = ''
  form.quantity = 1
  form.unitPrice = ''
  form.partStatus = 'original'
  editingProductId.value = null
}

function cancelEdit() {
  clearForm()
}

// Düzenleme moduna girildiğinde formu doldur
watch(editingProductId, (id) => {
  if (id) {
    const p = products.value.find(x => x.id === id)
    if (p) {
      form.name = p.product_name
      form.oem = p.oem_number === '-' ? '' : p.oem_number
      form.quantity = p.quantity
      form.unitPrice = p.unit_price
      form.partStatus = p.part_status
    }
  }
})

function handleClickOutside(e) {
  if (!e.target.closest('.relative')) {
    showProductDropdown.value = false
    showOEMDropdown.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.animate-modal {
  animation: modalIn 0.25s ease;
}
@keyframes modalIn {
  from { transform: scale(0.9) translateY(20px); opacity: 0; }
  to { transform: scale(1) translateY(0); opacity: 1; }
}
</style>
