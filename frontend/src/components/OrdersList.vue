<template>
  <div class="card">
    <!-- Header -->
    <div class="px-6 py-5 border-b flex items-center gap-4 bg-gradient-to-r from-accent/10 to-purple/10" style="border-color: var(--border-color);">
      <div class="w-12 h-12 bg-gradient-to-r from-accent to-purple rounded-xl flex items-center justify-center text-xl text-white">
        📚
      </div>
      <div class="flex-1">
        <h3 class="text-lg font-bold" style="color: var(--text-primary);">
          {{ isAdvancedSearch ? '🔎 Arama Sonuçları' : (customerFilterActive ? '👤 ' + customerName : 'Kayıtlı Siparişler') }}
        </h3>
        <span class="text-sm" style="color: var(--text-muted);">
          {{ orders.length }} sipariş
          <span v-if="customerFilterActive && !isAdvancedSearch" class="ml-2 px-2 py-0.5 bg-accent/20 text-accent rounded-full text-xs font-semibold">
            Müşteri Filtreli
          </span>
          <span v-if="isAdvancedSearch" class="ml-2 px-2 py-0.5 bg-warning/20 text-warning rounded-full text-xs font-semibold">
            Gelişmiş Arama
          </span>
        </span>
      </div>
      <button 
        v-if="isAdvancedSearch" 
        @click="clearAdvancedSearch" 
        class="btn btn-sm btn-secondary"
        title="Aramayı temizle"
      >
        ✕ Temizle
      </button>
    </div>

    <!-- Search & Filters -->
    <div class="p-4 border-b" style="border-color: var(--border-color);">
      <div class="flex gap-2 mb-3">
        <input 
          type="text" 
          v-model="searchTerm" 
          @input="onSearch"
          class="form-input flex-1" 
          placeholder="🔍 Hızlı ara..."
        >
        <button @click="$emit('showAdvancedSearch')" class="btn btn-secondary btn-sm" title="Gelişmiş Arama">🔎</button>
      </div>
      
      <div class="flex gap-2 flex-wrap">
        <button 
          v-for="f in filters" 
          :key="f.value"
          @click="setFilter(f.value)"
          :class="['filter-chip', { active: currentFilter === f.value }]"
        >
          {{ f.label }}
        </button>
      </div>
      
      <div v-if="currentFilter === 'range'" class="mt-3 space-y-2">
        <div class="flex flex-col gap-2">
          <div class="flex items-center gap-2">
            <span class="text-xs whitespace-nowrap" style="color: var(--text-muted);">Başlangıç:</span>
            <input type="date" v-model="startDate" class="form-input flex-1 !py-2 !text-sm">
          </div>
          <div class="flex items-center gap-2">
            <span class="text-xs whitespace-nowrap" style="color: var(--text-muted);">Bitiş:</span>
            <input type="date" v-model="endDate" class="form-input flex-1 !py-2 !text-sm">
          </div>
        </div>
        <button @click="loadOrders" class="btn btn-primary btn-sm w-full">📅 Uygula</button>
      </div>
    </div>

    <!-- Orders List -->
    <div class="flex-1 overflow-y-auto p-4">
      <div v-if="orders.length === 0" class="text-center py-12" style="color: var(--text-muted);">
        <div class="text-6xl mb-4 opacity-50">📭</div>
        <p>Sipariş bulunamadı</p>
      </div>
      
      <div 
        v-for="order in orders" 
        :key="order.id"
        @click="$emit('loadOrder', order.id)"
        :class="['order-card', { active: order.id === currentOrderId }]"
      >
        <div class="flex items-center gap-2 font-bold mb-2" style="color: var(--text-primary);">
          <span class="text-lg">📋</span>
          {{ order.title || 'İsimsiz Sipariş' }}
        </div>
        <div class="flex items-center gap-2 text-sm mb-3 pb-3 border-b" style="color: var(--text-muted); border-color: var(--border-color);">
          👤 {{ order.customer_name || 'Müşteri belirtilmemiş' }}
        </div>
        <div class="flex justify-between items-center">
          <div class="flex gap-4 text-xs" style="color: var(--text-muted);">
            <span>📦 {{ order.items?.length || 0 }}</span>
            <span>📅 {{ formatDate(order.created_at) }}</span>
          </div>
          <div class="font-extrabold text-success text-lg">₺{{ formatPrice(order.grand_total || 0) }}</div>
        </div>
        <div class="grid grid-cols-2 gap-2 mt-3 pt-3 border-t" style="border-color: var(--border-color);">
          <button @click.stop="$emit('loadOrder', order.id)" class="btn btn-sm btn-secondary">📝 Düzenle</button>
          <button @click.stop="$emit('deleteOrder', order.id)" class="btn btn-sm btn-danger">🗑️ Sil</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import { api } from '@/api'
import { useOrder } from '@/composables/useOrder'

const props = defineProps(['refreshTrigger', 'advancedSearchResults'])
const emit = defineEmits(['loadOrder', 'deleteOrder', 'showAdvancedSearch', 'clearAdvancedSearch'])

const { currentOrderId, currentCustomerId, customerFilterActive, customerName, advancedSearchFilter } = useOrder()

const allOrders = ref([])
const searchTerm = ref('')
const currentFilter = ref('today')
const startDate = ref('')
const endDate = ref('')
const isAdvancedSearch = ref(false)
let searchTimeout = null

const filters = [
  { label: 'Bugün', value: 'today' },
  { label: 'Tümü', value: 'all' },
  { label: 'Tarih Aralığı', value: 'range' },
]

// Müşteri filtresine göre siparişleri filtrele
const orders = computed(() => {
  let filtered = allOrders.value
  
  // Müşteri filtresi
  if (customerFilterActive.value && currentCustomerId.value) {
    filtered = filtered.filter(o => o.customer_id === currentCustomerId.value)
  }
  
  return filtered
})

function formatPrice(n) {
  return (n || 0).toLocaleString('tr-TR', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

function formatDate(d) {
  return d ? new Date(d).toLocaleDateString('tr-TR', { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' }) : '-'
}

async function loadOrders() {
  // Eğer gelişmiş arama aktifse, tarih filtresi ile birlikte gelişmiş arama yap
  if (isAdvancedSearch.value && advancedSearchFilter.value) {
    const filter = {
      ...advancedSearchFilter.value,
      date_filter: currentFilter.value,
      start_date: startDate.value,
      end_date: endDate.value
    }
    allOrders.value = await api.searchOrdersAdvanced(filter)
  } else {
    const filter = {
      type: currentFilter.value,
      start_date: startDate.value,
      end_date: endDate.value
    }
    allOrders.value = await api.loadOrders(filter)
  }
}

function setFilter(filter) {
  currentFilter.value = filter
  if (filter !== 'range') {
    loadOrders()
  }
}

function onSearch() {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(async () => {
    if (searchTerm.value.trim().length > 0) {
      allOrders.value = await api.searchOrders(searchTerm.value.trim())
    } else {
      loadOrders()
    }
  }, 300)
}

onMounted(() => {
  const today = new Date().toISOString().split('T')[0]
  endDate.value = today
  const lastMonth = new Date()
  lastMonth.setMonth(lastMonth.getMonth() - 1)
  startDate.value = lastMonth.toISOString().split('T')[0]
  
  loadOrders()
})

watch(() => props.refreshTrigger, () => {
  if (!isAdvancedSearch.value) {
    loadOrders()
  }
})

// Müşteri filtresi değiştiğinde siparişleri yeniden yükle
watch(customerFilterActive, () => {
  if (!isAdvancedSearch.value) {
    loadOrders()
  }
})

// Gelişmiş arama sonuçları geldiğinde
watch(() => props.advancedSearchResults, (results) => {
  if (results && results.length >= 0) {
    allOrders.value = results
    isAdvancedSearch.value = true
  }
})

function clearAdvancedSearch() {
  isAdvancedSearch.value = false
  advancedSearchFilter.value = null
  emit('clearAdvancedSearch')
  loadOrders()
}

defineExpose({ loadOrders })
</script>

<style scoped>
.filter-chip {
  background: var(--bg-secondary);
  border: 2px solid var(--border-color);
  color: var(--text-muted);
  @apply px-4 py-2 rounded-full text-sm font-semibold cursor-pointer transition-all hover:bg-accent hover:border-accent hover:text-white;
}
.filter-chip.active {
  @apply bg-accent border-accent text-white;
}
.order-card {
  background: var(--bg-secondary);
  border: 2px solid var(--border-color);
  @apply rounded-xl p-4 mb-3 cursor-pointer transition-all hover:border-accent hover:translate-x-1;
}
.order-card.active {
  @apply border-accent bg-accent/10;
}
</style>
