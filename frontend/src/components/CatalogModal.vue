<template>
  <Teleport to="body">
    <div v-if="visible" class="fixed inset-0 bg-black/70 backdrop-blur-sm flex items-center justify-center z-[1000] p-4">
      <div class="rounded-2xl w-full max-w-4xl max-h-[90vh] flex flex-col overflow-hidden border" style="background: var(--bg-card); border-color: var(--border-color);">
        <!-- Header -->
        <div class="p-6 bg-gradient-to-r from-accent/10 to-purple/10 border-b" style="border-color: var(--border-color);">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-4">
              <div class="w-12 h-12 bg-gradient-to-r from-accent to-purple rounded-xl flex items-center justify-center text-xl text-white">
                📚
              </div>
              <div>
                <h2 class="text-xl font-bold" style="color: var(--text-primary);">Ürün & Müşteri Kataloğu</h2>
                <p class="text-sm" style="color: var(--text-muted);">İndekslenmiş verileri görüntüle ve düzenle</p>
              </div>
            </div>
            <button @click="close" class="btn btn-secondary btn-sm">✕ Kapat</button>
          </div>
          
          <!-- Tabs -->
          <div class="flex gap-2 mt-4">
            <button 
              @click="activeTab = 'products'" 
              :class="['tab-btn', { active: activeTab === 'products' }]"
            >
              🔧 Ürünler ({{ products.length }})
            </button>
            <button 
              @click="activeTab = 'customers'" 
              :class="['tab-btn', { active: activeTab === 'customers' }]"
            >
              👤 Müşteriler ({{ customers.length }})
            </button>
          </div>
        </div>
        
        <!-- Search -->
        <div class="p-4 border-b" style="border-color: var(--border-color);">
          <input 
            type="text" 
            v-model="searchQuery" 
            class="form-input" 
            :placeholder="activeTab === 'products' ? '🔍 Ürün adı veya OEM ara...' : '🔍 Müşteri adı ara...'"
          >
        </div>
        
        <!-- Content -->
        <div class="flex-1 overflow-y-auto p-4">
          <!-- Products Tab -->
          <template v-if="activeTab === 'products'">
            <div v-if="filteredProducts.length === 0" class="text-center py-12" style="color: var(--text-muted);">
              <div class="text-5xl mb-4 opacity-50">📦</div>
              <p>Ürün bulunamadı</p>
            </div>
            
            <div v-else class="grid gap-3">
              <div 
                v-for="product in filteredProducts" 
                :key="product.id"
                class="p-4 rounded-xl border flex items-center gap-4 transition-all hover:border-accent"
                style="background: var(--bg-secondary); border-color: var(--border-color);"
              >
                <div class="flex-1">
                  <div class="font-bold mb-1" style="color: var(--text-primary);">{{ product.name }}</div>
                  <div class="text-sm" style="color: var(--text-muted);">
                    OEM: {{ product.oem_number || '-' }} • {{ product.brand || '-' }}
                  </div>
                </div>
                <div class="stock-info">
                  <div class="stock-current" :class="{ 'stock-critical': product.stock_quantity <= (product.critical_stock || 3) }">
                    {{ product.stock_quantity || 0 }} {{ product.unit || 'adet' }}
                  </div>
                  <div class="stock-critical-level">
                    Kritik: {{ product.critical_stock || 3 }}
                  </div>
                </div>
                <div class="flex gap-2">
                  <button @click="editProduct(product)" class="btn btn-sm btn-secondary">✏️</button>
                  <button @click="deleteProduct(product)" class="btn btn-sm btn-danger">🗑️</button>
                </div>
              </div>
            </div>
          </template>
          
          <!-- Customers Tab -->
          <template v-if="activeTab === 'customers'">
            <div v-if="filteredCustomers.length === 0" class="text-center py-12" style="color: var(--text-muted);">
              <div class="text-5xl mb-4 opacity-50">👥</div>
              <p>Müşteri bulunamadı</p>
            </div>
            
            <div v-else class="grid gap-3">
              <div 
                v-for="customer in filteredCustomers" 
                :key="customer.id"
                class="p-4 rounded-xl border flex items-center gap-4 transition-all hover:border-accent"
                style="background: var(--bg-secondary); border-color: var(--border-color);"
              >
                <div class="flex-1">
                  <div class="flex items-center gap-2 mb-1">
                    <span class="font-bold" style="color: var(--text-primary);">{{ customer.name }}</span>
                    <span class="text-xs px-2 py-0.5 rounded bg-accent/20 text-accent font-mono">#{{ customer.id?.substring(0, 6) }}</span>
                  </div>
                  <div class="text-sm" style="color: var(--text-muted);">
                    <span v-if="customer.phone" class="mr-2">📞 {{ customer.phone }}</span>
                    <span>{{ customer.order_count || 0 }} sipariş • Toplam: ₺{{ formatPrice(customer.total_amount || 0) }}</span>
                  </div>
                </div>
                <div class="flex gap-2">
                  <button @click="editCustomer(customer)" class="btn btn-sm btn-secondary">✏️</button>
                  <button @click="deleteCustomer(customer)" class="btn btn-sm btn-danger">🗑️</button>
                </div>
              </div>
            </div>
          </template>
        </div>
      </div>
    </div>
    
    <!-- Edit Modal -->
    <div v-if="editModal.visible" class="fixed inset-0 bg-black/80 flex items-center justify-center z-[1001] p-4">
      <div class="rounded-2xl w-full max-w-md overflow-hidden border" style="background: var(--bg-card); border-color: var(--border-color);">
        <div class="p-6 bg-gradient-to-r from-accent/10 to-purple/10 text-center">
          <div class="text-4xl mb-2">{{ editModal.type === 'product' ? '🔧' : '👤' }}</div>
          <h3 class="text-lg font-bold" style="color: var(--text-primary);">
            {{ editModal.type === 'product' ? 'Ürün Düzenle' : 'Müşteri Düzenle' }}
          </h3>
        </div>
        <div class="p-6">
          <div class="mb-4">
            <label class="block mb-2 text-sm font-semibold" style="color: var(--text-muted);">
              {{ editModal.type === 'product' ? 'Ürün Adı' : 'Müşteri Adı' }}
            </label>
            <input type="text" v-model="editModal.name" class="form-input">
          </div>
          <div v-if="editModal.type === 'product'" class="mb-4">
            <label class="block mb-2 text-sm font-semibold" style="color: var(--text-muted);">OEM Numarası</label>
            <input type="text" v-model="editModal.oem" class="form-input">
          </div>
          <div v-if="editModal.type === 'customer'" class="mb-4">
            <label class="block mb-2 text-sm font-semibold" style="color: var(--text-muted);">📞 Telefon Numarası</label>
            <div class="flex items-stretch rounded-xl overflow-hidden border-2" style="border-color: var(--border-color);">
              <input 
                type="text" 
                v-model="editModal.countryCode" 
                class="w-14 text-center text-base font-medium border-0 border-r-2 focus:ring-0 focus:outline-none px-2 py-3.5" 
                style="background: var(--bg-secondary); border-color: var(--border-color); color: var(--text-primary);"
                placeholder="+90"
              >
              <input 
                type="tel" 
                v-model="editModal.phone" 
                class="flex-1 min-w-0 px-4 py-3.5 text-base border-0 focus:ring-0 focus:outline-none" 
                style="background: var(--bg-secondary); color: var(--text-primary);"
                placeholder="5XX XXX XX XX"
              >
            </div>
          </div>
        </div>
        <div class="grid grid-cols-2 gap-3 p-4" style="background: var(--bg-secondary);">
          <button @click="editModal.visible = false" class="btn btn-secondary">İptal</button>
          <button @click="saveEdit" class="btn btn-primary">💾 Kaydet</button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { api } from '@/api'
import { useToast } from '@/composables/useToast'

const props = defineProps(['visible'])
const emit = defineEmits(['close', 'updated'])

const { showToast } = useToast()

const activeTab = ref('products')
const searchQuery = ref('')
const products = ref([])
const customers = ref([])

const editModal = reactive({
  visible: false,
  type: '',
  id: '',
  name: '',
  oem: '',
  phone: '',
  countryCode: '+90'
})

const filteredProducts = computed(() => {
  if (!searchQuery.value) return products.value
  const q = searchQuery.value.toLowerCase()
  return products.value.filter(p => 
    p.name.toLowerCase().includes(q) || 
    (p.oem_number && p.oem_number.toLowerCase().includes(q))
  )
})

const filteredCustomers = computed(() => {
  if (!searchQuery.value) return customers.value
  const q = searchQuery.value.toLowerCase()
  return customers.value.filter(c => 
    c.name.toLowerCase().includes(q) ||
    (c.phone && c.phone.includes(q))
  )
})

function formatPrice(n) {
  return (n || 0).toLocaleString('tr-TR', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

async function loadData() {
  products.value = await api.listProducts()
  customers.value = await api.listCustomers()
}

function close() {
  emit('close')
}

function editProduct(product) {
  editModal.type = 'product'
  editModal.id = product.id
  editModal.name = product.name
  editModal.oem = product.oem_number || ''
  editModal.visible = true
}

function editCustomer(customer) {
  editModal.type = 'customer'
  editModal.id = customer.id
  editModal.name = customer.name
  // Ülke kodunu ve numarayı ayır
  if (customer.phone) {
    const match = customer.phone.match(/^(\+\d{1,4})\s*(.*)$/)
    if (match) {
      editModal.countryCode = match[1]
      editModal.phone = match[2]
    } else {
      editModal.countryCode = '+90'
      editModal.phone = customer.phone
    }
  } else {
    editModal.countryCode = '+90'
    editModal.phone = ''
  }
  editModal.visible = true
}

async function saveEdit() {
  if (!editModal.name.trim()) {
    showToast('İsim boş olamaz', 'error')
    return
  }
  
  try {
    if (editModal.type === 'product') {
      await api.updateProduct({
        id: editModal.id,
        name: editModal.name.trim(),
        oem_number: editModal.oem.trim()
      })
      showToast('Ürün güncellendi')
    } else {
      // Ülke kodu ve numara birleştir
      const fullPhone = editModal.phone.trim() 
        ? (editModal.countryCode || '+90') + ' ' + editModal.phone.trim()
        : ''
      await api.updateCustomer({
        id: editModal.id,
        name: editModal.name.trim(),
        phone: fullPhone
      })
      showToast('Müşteri güncellendi')
    }
    editModal.visible = false
    await loadData()
    emit('updated')
  } catch (e) {
    showToast('Güncelleme hatası', 'error')
  }
}

async function deleteProduct(product) {
  if (!confirm(`"${product.name}" ürününü silmek istediğinize emin misiniz?`)) return
  
  try {
    await api.deleteProduct(product.id)
    showToast('Ürün silindi')
    await loadData()
    emit('updated')
  } catch (e) {
    showToast('Silme hatası', 'error')
  }
}

async function deleteCustomer(customer) {
  if (!confirm(`"${customer.name}" müşterisini silmek istediğinize emin misiniz?`)) return
  
  try {
    await api.deleteCustomer(customer.id)
    showToast('Müşteri silindi')
    await loadData()
    emit('updated')
  } catch (e) {
    showToast('Silme hatası', 'error')
  }
}

watch(() => props.visible, (val) => {
  if (val) {
    loadData()
    searchQuery.value = ''
  }
})

onMounted(() => {
  if (props.visible) {
    loadData()
  }
})
</script>

<style scoped>
.tab-btn {
  @apply px-4 py-2 rounded-lg text-sm font-semibold transition-all;
  background: var(--bg-secondary);
  color: var(--text-muted);
  border: 2px solid transparent;
}
.tab-btn:hover {
  border-color: var(--border-color);
}
.tab-btn.active {
  @apply bg-accent text-white border-accent;
}

.stock-info {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  min-width: 80px;
}

.stock-current {
  font-weight: 600;
  font-size: 1rem;
  color: var(--success-color);
}

.stock-current.stock-critical {
  color: var(--danger-color);
}

.stock-critical-level {
  font-size: 0.75rem;
  color: var(--text-muted);
}
</style>
