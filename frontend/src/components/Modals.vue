<template>
  <Teleport to="body">
    <!-- Confirm Modal -->
    <div v-if="confirmVisible" class="modal-overlay">
      <div class="modal max-w-md">
        <div class="modal-header">
          <div class="text-5xl">{{ confirmIcon }}</div>
          <div class="text-xl font-bold mt-4" style="color: var(--text-primary);">{{ confirmTitle }}</div>
          <div class="mt-2" style="color: var(--text-muted);">{{ confirmMessage }}</div>
        </div>
        <div class="modal-actions grid-cols-2">
          <button @click="closeConfirm" class="btn btn-secondary">İptal</button>
          <button @click="executeConfirm" class="btn btn-danger">Onayla</button>
        </div>
      </div>
    </div>

    <!-- Customer History Modal -->
    <div v-if="historyVisible" class="modal-overlay">
      <div class="modal max-w-xl">
        <div class="modal-header">
          <div class="text-5xl">📊</div>
          <div class="text-xl font-bold mt-4" style="color: var(--text-primary);">{{ historyCustomer?.name }}</div>
          <div class="mt-2" style="color: var(--text-muted);">
            {{ historyOrders.length }} sipariş • Toplam: ₺{{ formatPrice(historyCustomer?.total_amount || 0) }}
          </div>
        </div>
        <div class="modal-body max-h-[400px] overflow-y-auto">
          <div v-if="historyOrders.length === 0" class="text-center py-8" style="color: var(--text-muted);">
            <p>Sipariş bulunamadı</p>
          </div>
          <div 
            v-for="order in historyOrders" 
            :key="order.id"
            @click="loadFromHistory(order.id)"
            class="order-history-card"
          >
            <div class="font-bold" style="color: var(--text-primary);">📋 {{ order.title || 'İsimsiz' }}</div>
            <div class="flex justify-between items-center mt-2">
              <div class="flex gap-4 text-xs" style="color: var(--text-muted);">
                <span>📦 {{ order.items?.length || 0 }}</span>
                <span>📅 {{ formatDate(order.created_at) }}</span>
              </div>
              <div class="font-bold text-success">₺{{ formatPrice(order.grand_total || 0) }}</div>
            </div>
          </div>
        </div>
        <div class="modal-actions grid-cols-1">
          <button @click="closeHistory" class="btn btn-secondary btn-block">Kapat</button>
        </div>
      </div>
    </div>

    <!-- Advanced Search Modal -->
    <div v-if="advancedSearchVisible" class="modal-overlay">
      <div class="modal max-w-xl">
        <div class="modal-header">
          <div class="text-5xl">🔎</div>
          <div class="text-xl font-bold mt-4" style="color: var(--text-primary);">Gelişmiş Sipariş Arama</div>
          <div class="mt-2" style="color: var(--text-muted);">Birden fazla kriter ile filtreleme yapın</div>
        </div>
        <div class="modal-body">
          <!-- Product Name -->
          <div class="mb-4 relative">
            <label class="block mb-2 text-sm font-semibold" style="color: var(--text-muted);">📦 Ürün Adı İçeren</label>
            <input 
              type="text" 
              v-model="advSearch.productName" 
              @input="onAdvProductInput"
              class="form-input" 
              placeholder="Ürün adı..."
              autocomplete="off"
            >
            <div v-if="advProductDropdown.length > 0" class="dropdown-menu">
              <div v-for="p in advProductDropdown" :key="p.id" @click="selectAdvProduct(p)" class="dropdown-item">
                <div class="font-semibold">📦 {{ p.name }}</div>
                <div class="text-sm" style="color: var(--text-muted);">OEM: {{ p.oem_number }}</div>
              </div>
            </div>
          </div>
          
          <!-- OEM Number -->
          <div class="mb-4 relative">
            <label class="block mb-2 text-sm font-semibold" style="color: var(--text-muted);">🔢 OEM Numarası İçeren</label>
            <input 
              type="text" 
              v-model="advSearch.oemNumber" 
              @input="onAdvOEMInput"
              class="form-input" 
              placeholder="OEM kodu..."
              autocomplete="off"
            >
            <div v-if="advOEMDropdown.length > 0" class="dropdown-menu">
              <div v-for="p in advOEMDropdown" :key="p.id" @click="selectAdvOEM(p)" class="dropdown-item">
                <div class="font-semibold">🔢 {{ p.oem_number }}</div>
                <div class="text-sm" style="color: var(--text-muted);">{{ p.name }}</div>
              </div>
            </div>
          </div>
          
          <!-- Customer Name -->
          <div class="mb-4 relative">
            <label class="block mb-2 text-sm font-semibold" style="color: var(--text-muted);">👤 Müşteri Adı İçeren</label>
            <input 
              type="text" 
              v-model="advSearch.customerName" 
              @input="onAdvCustomerInput"
              class="form-input" 
              placeholder="Müşteri adı..."
              autocomplete="off"
            >
            <div v-if="advCustomerDropdown.length > 0" class="dropdown-menu">
              <div v-for="c in advCustomerDropdown" :key="c.id" @click="selectAdvCustomer(c)" class="dropdown-item">
                <div class="font-semibold">👤 {{ c.name }}</div>
                <div class="text-sm" style="color: var(--text-muted);">{{ c.order_count || 0 }} sipariş</div>
              </div>
            </div>
          </div>
          
          <!-- Quantity Filters -->
          <div class="grid grid-cols-2 gap-4 mb-4">
            <div>
              <label class="block mb-2 text-sm font-semibold" style="color: var(--text-muted);">📊 Min. Adet</label>
              <input type="number" v-model.number="advSearch.minQty" min="0" class="form-input" placeholder="0">
            </div>
            <div>
              <label class="block mb-2 text-sm font-semibold" style="color: var(--text-muted);">📊 Max. Adet</label>
              <input type="number" v-model.number="advSearch.maxQty" min="0" class="form-input" placeholder="Sınırsız">
            </div>
          </div>
          
          <!-- Total Filters -->
          <div class="grid grid-cols-2 gap-4 mb-4">
            <div>
              <label class="block mb-2 text-sm font-semibold" style="color: var(--text-muted);">💰 Min. Toplam (₺)</label>
              <input type="number" v-model.number="advSearch.minTotal" min="0" step="0.01" class="form-input" placeholder="0">
            </div>
            <div>
              <label class="block mb-2 text-sm font-semibold" style="color: var(--text-muted);">💰 Max. Toplam (₺)</label>
              <input type="number" v-model.number="advSearch.maxTotal" min="0" step="0.01" class="form-input" placeholder="Sınırsız">
            </div>
          </div>
          
          <!-- Unit Price Filters -->
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block mb-2 text-sm font-semibold" style="color: var(--text-muted);">💵 Min. Birim Fiyat (₺)</label>
              <input type="number" v-model.number="advSearch.minUnitPrice" min="0" step="0.01" class="form-input" placeholder="0">
            </div>
            <div>
              <label class="block mb-2 text-sm font-semibold" style="color: var(--text-muted);">💵 Max. Birim Fiyat (₺)</label>
              <input type="number" v-model.number="advSearch.maxUnitPrice" min="0" step="0.01" class="form-input" placeholder="Sınırsız">
            </div>
          </div>
        </div>
        <div class="modal-actions grid-cols-2">
          <button @click="closeAdvancedSearch" class="btn btn-secondary">İptal</button>
          <button @click="executeAdvancedSearch" class="btn btn-primary">🔍 Ara</button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { api } from '@/api'
import { useOrder } from '@/composables/useOrder'
import { useToast } from '@/composables/useToast'

const emit = defineEmits(['advancedSearchResults', 'loadOrder'])

const { allProducts, allCustomers, customerName, currentCustomerId, loadOrder, advancedSearchFilter } = useOrder()
const { showToast } = useToast()

// Confirm Modal
const confirmVisible = ref(false)
const confirmTitle = ref('')
const confirmMessage = ref('')
const confirmIcon = ref('⚠️')
let confirmCallback = null

function showConfirm(title, message, callback, icon = '⚠️') {
  confirmTitle.value = title
  confirmMessage.value = message
  confirmIcon.value = icon
  confirmCallback = callback
  confirmVisible.value = true
}

function closeConfirm() {
  confirmVisible.value = false
  confirmCallback = null
}

function executeConfirm() {
  if (confirmCallback) confirmCallback()
  closeConfirm()
}

// Customer History Modal
const historyVisible = ref(false)
const historyCustomer = ref(null)
const historyOrders = ref([])

async function showHistory() {
  const name = customerName.value?.trim()
  if (!name && !currentCustomerId.value) {
    showToast('Önce müşteri seçin', 'error')
    return
  }

  const customer = allCustomers.value.find(c => 
    c.id === currentCustomerId.value || 
    c.name.toLowerCase() === name?.toLowerCase()
  )
  
  if (!customer) {
    showToast('Müşteri bulunamadı', 'error')
    return
  }

  historyCustomer.value = customer
  historyOrders.value = await api.getCustomerOrders(customer.id)
  historyVisible.value = true
}

function closeHistory() {
  historyVisible.value = false
}

function loadFromHistory(id) {
  closeHistory()
  emit('loadOrder', id)
}

// Advanced Search Modal
const advancedSearchVisible = ref(false)
const advSearch = reactive({
  productName: '',
  oemNumber: '',
  customerName: '',
  minQty: null,
  maxQty: null,
  minTotal: null,
  maxTotal: null,
  minUnitPrice: null,
  maxUnitPrice: null
})

const advProductDropdown = ref([])
const advOEMDropdown = ref([])
const advCustomerDropdown = ref([])

function showAdvancedSearch() {
  Object.keys(advSearch).forEach(k => advSearch[k] = k.includes('min') || k.includes('max') ? null : '')
  advProductDropdown.value = []
  advOEMDropdown.value = []
  advCustomerDropdown.value = []
  advancedSearchVisible.value = true
}

function closeAdvancedSearch() {
  advancedSearchVisible.value = false
}

function onAdvProductInput() {
  if (advSearch.productName.length < 2) {
    advProductDropdown.value = []
    return
  }
  const search = advSearch.productName.toLowerCase()
  advProductDropdown.value = allProducts.value.filter(p => p.name.toLowerCase().includes(search)).slice(0, 5)
}

function onAdvOEMInput() {
  if (advSearch.oemNumber.length < 2) {
    advOEMDropdown.value = []
    return
  }
  const search = advSearch.oemNumber.toLowerCase()
  advOEMDropdown.value = allProducts.value.filter(p => p.oem_number?.toLowerCase().includes(search)).slice(0, 5)
}

function onAdvCustomerInput() {
  if (advSearch.customerName.length < 1) {
    advCustomerDropdown.value = []
    return
  }
  const search = advSearch.customerName.toLowerCase()
  advCustomerDropdown.value = allCustomers.value.filter(c => c.name.toLowerCase().includes(search)).slice(0, 5)
}

function selectAdvProduct(p) {
  advSearch.productName = p.name
  advProductDropdown.value = []
}

function selectAdvOEM(p) {
  advSearch.oemNumber = p.oem_number
  advOEMDropdown.value = []
}

function selectAdvCustomer(c) {
  advSearch.customerName = c.name
  advCustomerDropdown.value = []
}

async function executeAdvancedSearch() {
  const filter = {
    product_name: advSearch.productName,
    oem_number: advSearch.oemNumber,
    customer_name: advSearch.customerName,
    min_quantity: advSearch.minQty || 0,
    max_quantity: advSearch.maxQty || 0,
    min_total: advSearch.minTotal || 0,
    max_total: advSearch.maxTotal || 0,
    min_unit_price: advSearch.minUnitPrice || 0,
    max_unit_price: advSearch.maxUnitPrice || 0
  }

  const hasFilter = filter.product_name || filter.oem_number || filter.customer_name ||
                    filter.min_quantity > 0 || filter.max_quantity > 0 ||
                    filter.min_total > 0 || filter.max_total > 0 ||
                    filter.min_unit_price > 0 || filter.max_unit_price > 0

  if (!hasFilter) {
    showToast('En az bir filtre kriteri girin', 'error')
    return
  }

  // Filtreleri global state'e kaydet
  advancedSearchFilter.value = { ...filter }

  const orders = await api.searchOrdersAdvanced(filter)
  emit('advancedSearchResults', orders)
  closeAdvancedSearch()
  showToast(`${orders.length} sonuç bulundu`, orders.length > 0 ? 'success' : 'error')
}

// Utilities
function formatPrice(n) {
  return (n || 0).toLocaleString('tr-TR', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

function formatDate(d) {
  return d ? new Date(d).toLocaleDateString('tr-TR', { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' }) : '-'
}

defineExpose({
  showConfirm,
  showHistory,
  showAdvancedSearch
})
</script>

<style scoped>
.modal-overlay {
  @apply fixed inset-0 bg-black/70 backdrop-blur-sm flex items-center justify-center z-[1000] p-4 overflow-y-auto;
}
.modal {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  @apply rounded-2xl w-full max-h-[90vh] overflow-y-auto my-auto;
  animation: modalIn 0.25s ease;
}
.modal-header {
  @apply p-6 text-center bg-gradient-to-r from-accent/10 to-purple/10 sticky top-0 z-10;
}
.modal-body {
  @apply p-6 overflow-y-auto;
}
.modal-actions {
  background: var(--bg-secondary);
  @apply grid gap-3 p-5 sticky bottom-0;
}
.dropdown-menu {
  background: var(--bg-card);
  border: 2px solid var(--border-color);
  @apply absolute top-full left-0 right-0 mt-1 rounded-xl max-h-64 overflow-y-auto z-50 shadow-xl;
}
.dropdown-item {
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
  color: var(--text-primary);
  @apply px-4 py-3 cursor-pointer last:border-b-0 hover:bg-accent hover:text-white transition-colors;
}
.order-history-card {
  background: var(--bg-secondary);
  @apply rounded-xl p-4 mb-3 cursor-pointer border-2 border-transparent transition-all hover:border-accent;
}
@keyframes modalIn {
  from { transform: scale(0.9) translateY(20px); opacity: 0; }
  to { transform: scale(1) translateY(0); opacity: 1; }
}
</style>
