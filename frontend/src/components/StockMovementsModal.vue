<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-container large">
      <div class="modal-header">
        <h2>
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 20h9"></path>
            <path d="M16.5 3.5a2.12 2.12 0 0 1 3 3L7 19l-4 1 1-4Z"></path>
          </svg>
          Stok Hareketleri
          <span v-if="selectedProduct" class="product-badge">{{ selectedProduct.name }}</span>
        </h2>
        <button @click="$emit('close')" class="close-btn">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M18 6 6 18"></path>
            <path d="m6 6 12 12"></path>
          </svg>
        </button>
      </div>
      
      <!-- Filters -->
      <div class="filter-bar">
        <div class="date-filters">
          <div class="form-group">
            <label>Başlangıç</label>
            <input type="date" v-model="filter.start" class="form-input" />
          </div>
          <div class="form-group">
            <label>Bitiş</label>
            <input type="date" v-model="filter.end" class="form-input" />
          </div>
        </div>
        
        <div class="quick-filters">
          <button 
            :class="{ active: quickFilter === 'today' }" 
            @click="applyQuickFilter('today')"
            class="quick-btn"
          >
            Bugün
          </button>
          <button 
            :class="{ active: quickFilter === 'week' }" 
            @click="applyQuickFilter('week')"
            class="quick-btn"
          >
            Bu Hafta
          </button>
          <button 
            :class="{ active: quickFilter === 'month' }" 
            @click="applyQuickFilter('month')"
            class="quick-btn"
          >
            Bu Ay
          </button>
          <button 
            :class="{ active: quickFilter === 'all' }" 
            @click="applyQuickFilter('all')"
            class="quick-btn"
          >
            Tümü
          </button>
        </div>
        
        <button @click="loadMovements" class="btn btn-primary btn-sm">
          <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 12a9 9 0 1 1-6.219-8.56"></path>
          </svg>
          Filtrele
        </button>
      </div>
      
      <div class="modal-body">
        <!-- Statistics -->
        <div class="stats-bar">
          <div class="stat-item entry">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 5v14"></path>
              <path d="M5 12h14"></path>
            </svg>
            <div class="stat-content">
              <span class="stat-value">{{ totalIn }}</span>
              <span class="stat-label">Toplam Giriş</span>
            </div>
          </div>
          <div class="stat-item exit">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M5 12h14"></path>
            </svg>
            <div class="stat-content">
              <span class="stat-value">{{ totalOut }}</span>
              <span class="stat-label">Toplam Çıkış</span>
            </div>
          </div>
          <div class="stat-item movement">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 20h9"></path>
              <path d="M16.5 3.5a2.12 2.12 0 0 1 3 3L7 19l-4 1 1-4Z"></path>
            </svg>
            <div class="stat-content">
              <span class="stat-value">{{ movements.length }}</span>
              <span class="stat-label">Hareket Sayısı</span>
            </div>
          </div>
        </div>
        
        <!-- Movements Table -->
        <div class="table-container">
          <table class="movements-table">
            <thead>
              <tr>
                <th>Tarih</th>
                <th v-if="!selectedProduct">Ürün</th>
                <th class="text-center">Tip</th>
                <th class="text-right">Miktar</th>
                <th>Açıklama</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="loading">
                <td :colspan="selectedProduct ? 4 : 5" class="loading-row">
                  <div class="loading-spinner"></div>
                  Yükleniyor...
                </td>
              </tr>
              <tr v-else-if="movements.length === 0">
                <td :colspan="selectedProduct ? 4 : 5" class="empty-row">
                  <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
                    <path d="M12 20h9"></path>
                    <path d="M16.5 3.5a2.12 2.12 0 0 1 3 3L7 19l-4 1 1-4Z"></path>
                  </svg>
                  <p>Bu tarih aralığında hareket bulunamadı</p>
                </td>
              </tr>
              <tr 
                v-else
                v-for="movement in paginatedMovements" 
                :key="movement.id"
                :class="movement.movement_type"
              >
                <td class="date-cell">
                  <div class="date-full">{{ formatDate(movement.date) }}</div>
                  <div class="date-time">{{ formatTime(movement.date) }}</div>
                </td>
                <td v-if="!selectedProduct" class="product-cell">{{ movement.product_name }}</td>
                <td class="text-center">
                  <span :class="['type-badge', movement.movement_type]">
                    <svg v-if="movement.movement_type === 'in'" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M12 5v14"></path>
                      <path d="M5 12h14"></path>
                    </svg>
                    <svg v-else xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M5 12h14"></path>
                    </svg>
                    {{ movement.movement_type === 'in' ? 'Giriş' : 'Çıkış' }}
                  </span>
                </td>
                <td class="text-right amount-cell">
                  <span :class="movement.movement_type === 'in' ? 'text-success' : 'text-danger'">
                    {{ movement.movement_type === 'in' ? '+' : '-' }}{{ movement.amount }}
                  </span>
                </td>
                <td class="note-cell">{{ movement.note || '-' }}</td> 
              </tr>
            </tbody>
          </table>
        </div>
        
        <!-- Pagination -->
        <div v-if="movements.length > 0" class="pagination-bar">
          <div class="pagination-info">
            <span v-if="!showAll">
              {{ (currentPage - 1) * pageSize + 1 }}-{{ Math.min(currentPage * pageSize, movements.length) }} / {{ movements.length }} hareket
            </span>
            <span v-else>
              Tüm {{ movements.length }} hareket gösteriliyor
            </span>
          </div>
          
          <div class="pagination-controls">
            <button 
              @click="toggleShowAll" 
              class="btn btn-sm"
              :class="showAll ? 'btn-primary' : 'btn-secondary'"
            >
              {{ showAll ? 'Sayfalı Göster' : 'Tümünü Göster' }}
            </button>
            
            <template v-if="!showAll && totalPages > 1">
              <button 
                @click="goToPage(1)" 
                :disabled="currentPage === 1"
                class="btn btn-sm btn-secondary"
              >
                &laquo;
              </button>
              <button 
                @click="goToPage(currentPage - 1)" 
                :disabled="currentPage === 1"
                class="btn btn-sm btn-secondary"
              >
                &lsaquo;
              </button>
              
              <span class="page-indicator">{{ currentPage }} / {{ totalPages }}</span>
              
              <button 
                @click="goToPage(currentPage + 1)" 
                :disabled="currentPage === totalPages"
                class="btn btn-sm btn-secondary"
              >
                &rsaquo;
              </button>
              <button 
                @click="goToPage(totalPages)" 
                :disabled="currentPage === totalPages"
                class="btn btn-sm btn-secondary"
              >
                &raquo;
              </button>
            </template>
          </div>
        </div>
      </div>
      
      <div class="modal-footer">
        <button @click="$emit('close')" class="btn btn-secondary">Kapat</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useStock } from '../composables/useStock'

const props = defineProps({
  selectedProduct: {
    type: Object,
    default: null
  }
})

defineEmits(['close'])

const { stockMovements, loadStockMovements, loading } = useStock()

const movements = ref([])
const quickFilter = ref('month')
const filter = ref({
  start: '',
  end: ''
})

// Pagination
const currentPage = ref(1)
const pageSize = ref(25)
const showAll = ref(false)

// Paginated movements
const paginatedMovements = computed(() => {
  if (showAll.value) return movements.value
  const start = (currentPage.value - 1) * pageSize.value
  return movements.value.slice(start, start + pageSize.value)
})

const totalPages = computed(() => Math.ceil(movements.value.length / pageSize.value))

function goToPage(page) {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
  }
}

function toggleShowAll() {
  showAll.value = !showAll.value
  if (!showAll.value) {
    currentPage.value = 1
  }
}

// Statistics
const totalIn = computed(() => {
  return movements.value
    .filter(m => m.movement_type === 'in')
    .reduce((acc, m) => acc + m.amount, 0)
    .toFixed(1)
})

const totalOut = computed(() => {
  return movements.value
    .filter(m => m.movement_type === 'out')
    .reduce((acc, m) => acc + m.amount, 0)
    .toFixed(1)
})

// Apply quick filter
const applyQuickFilter = (type) => {
  quickFilter.value = type
  const today = new Date()
  
  switch (type) {
    case 'today':
      filter.value.start = formatDateInput(today)
      filter.value.end = formatDateInput(today)
      break
    case 'week':
      const weekStart = new Date(today)
      weekStart.setDate(today.getDate() - today.getDay() + 1) // Monday
      filter.value.start = formatDateInput(weekStart)
      filter.value.end = formatDateInput(today)
      break
    case 'month':
      const monthStart = new Date(today.getFullYear(), today.getMonth(), 1)
      filter.value.start = formatDateInput(monthStart)
      filter.value.end = formatDateInput(today)
      break
    case 'all':
      filter.value.start = ''
      filter.value.end = ''
      break
  }
  
  loadMovements()
}

// Load movements
const loadMovements = async () => {
  await loadStockMovements(
    props.selectedProduct?.id || '',
    filter.value.start,
    filter.value.end
  )
  movements.value = stockMovements.value
} 

// Date formatting
const formatDate = (dateStr) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('tr-TR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric'
  })
}

const formatTime = (dateStr) => {
  const date = new Date(dateStr)
  return date.toLocaleTimeString('tr-TR', {
    hour: '2-digit',
    minute: '2-digit'
  })
}

const formatDateInput = (date) => {
  return date.toISOString().split('T')[0]
}

// Reload when product changes
watch(() => props.selectedProduct, () => {
  applyQuickFilter('month')
})

onMounted(() => {
  applyQuickFilter('month')
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

.modal-container.large {
  max-width: 900px;
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
  color: #3b82f6;
}

.product-badge {
  background: var(--accent-color);
  color: white;
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.875rem;
  font-weight: 500;
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

.filter-bar {
  display: flex;
  align-items: flex-end;
  gap: 1rem;
  padding: 1rem 1.5rem;
  border-bottom: 1px solid var(--border-color);
  flex-wrap: wrap;
}

.date-filters {
  display: flex;
  gap: 0.75rem;
}

.date-filters .form-group {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.date-filters label {
  font-size: 0.75rem;
  color: var(--text-muted);
}

.date-filters input {
  padding: 0.5rem;
  font-size: 0.875rem;
}

.quick-filters {
  display: flex;
  gap: 0.25rem;
}

.quick-btn {
  padding: 0.5rem 0.75rem;
  border: 1px solid var(--border-color);
  border-radius: 0.375rem;
  background: transparent;
  color: var(--text-muted);
  font-size: 0.75rem;
  cursor: pointer;
  transition: all 0.15s;
}

.quick-btn:hover {
  background: var(--bg-secondary);
  color: var(--text-primary);
}

.quick-btn.active {
  background: var(--accent-color);
  border-color: var(--accent-color);
  color: white;
}

.modal-body {
  padding: 1.5rem;
  overflow-y: auto;
  flex: 1;
}

.stats-bar {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem;
  border-radius: 0.5rem;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
}

.stat-item.entry svg {
  color: #22c55e;
}

.stat-item.exit svg {
  color: #ef4444;
}

.stat-item.movement svg {
  color: #3b82f6;
}

.stat-content {
  display: flex;
  flex-direction: column;
}

.stat-value {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-primary);
}

.stat-label {
  font-size: 0.75rem;
  color: var(--text-muted);
}

.table-container {
  background: var(--bg-card);
  border-radius: 0.5rem;
  border: 1px solid var(--border-color);
  overflow: auto;
  max-height: 400px;
}

.movements-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.875rem;
}

.movements-table th {
  background: var(--bg-secondary);
  padding: 0.75rem 1rem;
  text-align: left;
  font-weight: 600;
  color: var(--text-primary);
  border-bottom: 2px solid var(--border-color);
  position: sticky;
  top: 0;
  z-index: 10;
}

.movements-table td {
  padding: 0.75rem 1rem;
  border-bottom: 1px solid var(--border-color);
  color: var(--text-primary);
}

.movements-table tbody tr:hover {
  background: var(--bg-secondary);
}

.movements-table tbody tr.entry {
  background: rgba(34, 197, 94, 0.03);
}

.movements-table tbody tr.exit {
  background: rgba(239, 68, 68, 0.03);
}

.text-center {
  text-align: center;
}

.text-right {
  text-align: right;
}

.date-cell {
  white-space: nowrap;
}

.date-full {
  font-weight: 500;
}

.date-time {
  font-size: 0.75rem;
  color: var(--text-muted);
}

.type-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.25rem 0.5rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 500;
}

.type-badge.entry {
  background: rgba(34, 197, 94, 0.1);
  color: #22c55e;
}

.type-badge.exit {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
}

.amount-cell {
  font-weight: 600;
  font-size: 1rem;
}

.text-success {
  color: #22c55e;
}

.text-danger {
  color: #ef4444;
}

.note-cell {
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.loading-row, .empty-row {
  text-align: center;
  padding: 2rem !important;
  color: var(--text-muted);
}

.empty-row {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
}

.loading-spinner {
  width: 2rem;
  height: 2rem;
  border: 3px solid var(--border-color);
  border-top-color: var(--accent-color);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 0.5rem;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  padding: 1.25rem 1.5rem;
  border-top: 1px solid var(--border-color);
}

.pagination-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem;
  background: var(--bg-secondary);
  border-radius: 0.5rem;
  margin-top: 1rem;
}

.pagination-info {
  font-size: 0.875rem;
  color: var(--text-muted);
}

.pagination-controls {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.page-indicator {
  padding: 0 0.75rem;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-primary);
}

@media (max-width: 768px) {
  .filter-bar {
    flex-direction: column;
    align-items: stretch;
  }
  
  .stats-bar {
    grid-template-columns: 1fr;
  }
}
</style>
