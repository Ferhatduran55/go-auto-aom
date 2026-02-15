<template>
  <div class="card">
    <!-- Edit Banner -->
    <div 
      v-if="isEditing" 
      class="sticky top-0 z-20 backdrop-blur-md bg-accent/85 text-white px-4 py-2 flex items-center justify-between text-sm font-medium border-b border-white/10"
    >
      <span class="flex items-center gap-2">
        <span class="w-2 h-2 bg-white rounded-full animate-pulse"></span>
        {{ t('orders.editingCurrent') }} <!--: #{{ currentOrderId.substring(0, 8) }}-->
      </span>
      <button @click="$emit('newOrder')" class="px-3 py-1 bg-white/20 hover:bg-white/30 rounded-lg text-xs font-semibold transition-colors">{{ t('orders.newStart') }}</button>
    </div>
    
    <!-- Header -->
    <div class="px-6 py-5 border-b flex items-center gap-4 bg-gradient-to-r from-accent/10 to-purple/10" style="border-color: var(--border-color);">
      <div class="w-12 h-12 bg-gradient-to-r from-accent to-purple rounded-xl flex items-center justify-center text-xl text-white">
        📋
      </div>
      <div class="flex-1">
        <h3 class="text-lg font-bold" style="color: var(--text-primary);">{{ t('orders.orderItems') }}</h3>
        <span class="text-sm" style="color: var(--text-muted);">
          {{ isEditing ? '#' + currentOrderId.substring(0, 8) : t('orders.new') }}
          <span v-if="products.length > 0"> • {{ products.length }} {{ t('common.product') }}</span>
        </span>
      </div>
      <!-- Search Box -->
      <div v-if="products.length > 0" class="flex items-center gap-2">
        <input 
          type="text" 
          v-model="searchTerm" 
          class="form-input !py-2 !text-sm w-48" 
          :placeholder="'🔍 ' + t('orders.searchProduct')"
        >
        <button 
          v-if="searchTerm" 
          @click="searchTerm = ''" 
          class="text-sm px-2 py-1 rounded hover:bg-accent/20"
          style="color: var(--text-muted);"
        >✕</button>
      </div>
    </div>

    <!-- Multi-select Actions Bar -->
    <div v-if="selectedIds.length > 0" class="px-4 py-3 bg-danger/10 border-b flex items-center justify-between" style="border-color: var(--border-color);">
      <span class="text-sm font-semibold" style="color: var(--text-primary);">
        {{ t('stock.productsSelected', { count: selectedIds.length }) }}
      </span>
      <div class="flex gap-2">
        <button @click="selectAll" class="btn btn-sm btn-secondary">
          {{ selectedIds.length === products.length ? '☐ ' + t('common.deselectAll') : '☑ ' + t('common.selectAll') }}
        </button>
        <button @click="deleteSelected" class="btn btn-sm btn-danger">🗑️ {{ t('orders.deleteSelectedProducts') }}</button>
      </div>
    </div>

    <!-- Filter Info -->
    <div v-if="hiddenCount > 0" class="px-4 py-2 text-sm text-center" style="background: var(--bg-secondary); color: var(--text-muted);">
      ⚠️ {{ t('orders.hiddenByFilter', { count: hiddenCount }) }}
    </div>

    <!-- Table Body -->
    <div class="flex-1 overflow-auto">
      <table v-if="filteredProducts.length > 0" class="w-full border-collapse">
        <thead class="sticky top-0 z-10" style="background: var(--bg-secondary);">
          <tr>
            <th class="px-2 py-4 w-8">
              <input 
                type="checkbox" 
                :checked="selectedIds.length === filteredProducts.length && filteredProducts.length > 0"
                @change="toggleAllFiltered"
                class="w-4 h-4 cursor-pointer accent-accent"
              >
            </th>
            <th class="px-4 py-4 text-left font-bold text-xs uppercase tracking-wide" style="color: var(--text-muted);">{{ t('products.name') }}</th>
            <th class="px-4 py-4 text-left font-bold text-xs uppercase tracking-wide" style="color: var(--text-muted);">{{ t('common.oem') }}</th>
            <th class="px-4 py-4 text-left font-bold text-xs uppercase tracking-wide" style="color: var(--text-muted);">{{ t('common.quantity') }}</th>
            <th class="px-4 py-4 text-left font-bold text-xs uppercase tracking-wide" style="color: var(--text-muted);">{{ t('common.unitPrice') }}</th>
            <th class="px-4 py-4 text-left font-bold text-xs uppercase tracking-wide" style="color: var(--text-muted);">{{ t('common.status') }}</th>
            <th class="px-4 py-4 text-left font-bold text-xs uppercase tracking-wide" style="color: var(--text-muted);">{{ t('common.total') }}</th>
            <th class="px-4 py-4"></th>
          </tr>
        </thead>
        <tbody>
          <tr 
            v-for="p in filteredProducts" 
            :key="p.id" 
            class="border-b hover:bg-accent/5 transition-colors" 
            :class="{ 'bg-accent/10': selectedIds.includes(p.id) }"
            style="border-color: var(--border-color);"
          >
            <td class="px-2 py-4">
              <input 
                type="checkbox" 
                :checked="selectedIds.includes(p.id)"
                @change="toggleSelect(p.id)"
                class="w-4 h-4 cursor-pointer accent-accent"
              >
            </td>
            <td class="px-4 py-4 font-semibold" style="color: var(--text-primary);">{{ p.product_name }}</td>
            <td class="px-4 py-4 text-sm" style="color: var(--text-muted);">{{ p.oem_number }}</td>
            <td class="px-4 py-4" style="color: var(--text-primary);">{{ p.quantity }}</td>
            <td class="px-4 py-4" style="color: var(--text-primary);">{{ currency }}{{ formatPrice(p.unit_price) }}</td>
            <td class="px-4 py-4">
              <span :class="['badge', p.part_status === 'original' ? 'badge-success' : (p.part_status === 'zero' ? 'badge-primary' : 'badge-warning')]">
                {{ p.part_status === 'original' ? t('orders.conditionOriginal') : (p.part_status === 'zero' ? t('orders.conditionNew') : t('orders.conditionUsed')) }}
              </span>
            </td>
            <td class="px-4 py-4 font-bold text-success">{{ currency }}{{ formatPrice(p.total_price) }}</td>
            <td class="px-4 py-4">
              <button @click="editProduct(p.id)" class="action-btn" :title="t('common.edit')">✏️</button>
              <button @click="$emit('deleteProduct', p.id)" class="action-btn" :title="t('common.delete')">🗑️</button>
            </td>
          </tr>
        </tbody>
      </table>

      <!-- Empty State -->
      <div v-else-if="products.length === 0" class="text-center py-16" style="color: var(--text-muted);">
        <div class="text-6xl mb-4 opacity-50">📦</div>
        <p class="text-lg">{{ t('orders.noItems') }}</p>
      </div>

      <!-- No Results State -->
      <div v-else class="text-center py-16" style="color: var(--text-muted);">
        <div class="text-6xl mb-4 opacity-50">🔍</div>
        <p class="text-lg">{{ t('orders.noSearchResults') }}</p>
        <button @click="searchTerm = ''" class="btn btn-sm btn-secondary mt-4">{{ t('common.clearSearch') }}</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useOrder } from '@/composables/useOrder'
import { useI18n } from '@/i18n'
import { useSettings } from '@/composables/useSettings'

const { t, locale } = useI18n()
const { settings } = useSettings()
const emit = defineEmits(['newOrder', 'deleteProduct', 'deleteProducts'])

const { products, currentOrderId, isEditing, editingProductId, deleteProduct } = useOrder()

const searchTerm = ref('')
const selectedIds = ref([])

// Filtrelenmiş ürünler
const filteredProducts = computed(() => {
  if (!searchTerm.value.trim()) return products.value
  
  const search = searchTerm.value.toLowerCase().trim()
  return products.value.filter(p => {
    const statusText = p.part_status === 'original' ? 'orijinal' : (p.part_status === 'zero' ? 'sıfır' : 'çıkma')
    return (
      p.product_name.toLowerCase().includes(search) ||
      p.oem_number.toLowerCase().includes(search) ||
      p.unit_price.toString().includes(search) ||
      statusText.includes(search)
    )
  })
})

// Gizlenen ürün sayısı
const hiddenCount = computed(() => {
  return products.value.length - filteredProducts.value.length
})

function formatPrice(n) {
  const loc = locale.value || navigator.language || 'en-US'
  return (n || 0).toLocaleString(loc, { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

function editProduct(id) {
  editingProductId.value = id
}

function toggleSelect(id) {
  const index = selectedIds.value.indexOf(id)
  if (index > -1) {
    selectedIds.value.splice(index, 1)
  } else {
    selectedIds.value.push(id)
  }
}

function toggleAllFiltered() {
  if (selectedIds.value.length === filteredProducts.value.length) {
    selectedIds.value = []
  } else {
    selectedIds.value = filteredProducts.value.map(p => p.id)
  }
}

function selectAll() {
  if (selectedIds.value.length === products.value.length) {
    selectedIds.value = []
  } else {
    selectedIds.value = products.value.map(p => p.id)
  }
}

function deleteSelected() {
  if (selectedIds.value.length === 0) return
  
  // Seçili ürünleri sil
  selectedIds.value.forEach(id => {
    deleteProduct(id)
  })
  selectedIds.value = []
}
</script>

<style scoped>
.action-btn {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  color: var(--text-muted);
  @apply cursor-pointer p-2 rounded-lg text-base transition-all ml-1 hover:bg-accent hover:text-white hover:border-accent;
}
</style>
