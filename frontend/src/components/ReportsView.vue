<template>
  <div class="reports-view">
    <!-- Report Selection -->
    <div class="report-header">
      <div class="period-tabs">
        <button 
          :class="{ active: activePeriod === 'daily' }" 
          @click="activePeriod = 'daily'"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect width="18" height="18" x="3" y="4" rx="2" ry="2"></rect>
            <line x1="16" x2="16" y1="2" y2="6"></line>
            <line x1="8" x2="8" y1="2" y2="6"></line>
            <line x1="3" x2="21" y1="10" y2="10"></line>
          </svg>
          Günlük Rapor
        </button>
        <button 
          :class="{ active: activePeriod === 'monthly' }" 
          @click="activePeriod = 'monthly'"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect width="18" height="18" x="3" y="4" rx="2" ry="2"></rect>
            <line x1="16" x2="16" y1="2" y2="6"></line>
            <line x1="8" x2="8" y1="2" y2="6"></line>
            <line x1="3" x2="21" y1="10" y2="10"></line>
            <path d="M8 14h.01"></path>
            <path d="M12 14h.01"></path>
            <path d="M16 14h.01"></path>
            <path d="M8 18h.01"></path>
            <path d="M12 18h.01"></path>
            <path d="M16 18h.01"></path>
          </svg>
          Aylık Rapor
        </button>
      </div>
      
      <div class="date-selection">
        <input 
          v-if="activePeriod === 'daily'"
          type="date" 
          v-model="selectedDate" 
          class="form-input"
        />
        <input 
          v-else
          type="month" 
          v-model="selectedMonth" 
          class="form-input"
        />
        <button @click="generateReportHandler" class="btn btn-primary" :disabled="loading">
          <svg v-if="!loading" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 12a9 9 0 1 1-6.219-8.56"></path>
          </svg>
          <span v-if="loading" class="loading-spinner-sm"></span>
          Rapor Oluştur
        </button>
      </div>
    </div>
    
    <!-- Report Content -->
    <div v-if="report" class="report-content">
      <!-- Summary Cards -->
      <div class="summary-cards">
        <div class="summary-card entry">
          <div class="card-icon">
            <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 5v14"></path>
              <path d="M5 12h14"></path>
            </svg>
          </div>
          <div class="card-content">
            <span class="card-value">{{ report.total_in?.toFixed(1) || 0 }}</span>
            <span class="card-label">Toplam Giriş</span>
            <span class="card-count">{{ report.in_movement_count || 0 }} işlem</span>
          </div>
        </div>
        
        <div class="summary-card exit">
          <div class="card-icon">
            <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M5 12h14"></path>
            </svg>
          </div>
          <div class="card-content">
            <span class="card-value">{{ report.total_out?.toFixed(1) || 0 }}</span>
            <span class="card-label">Toplam Çıkış</span>
            <span class="card-count">{{ report.out_movement_count || 0 }} işlem</span>
          </div>
        </div>
        
        <div class="summary-card net">
          <div class="card-icon">
            <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"></path>
              <circle cx="9" cy="7" r="4"></circle>
              <line x1="19" x2="19" y1="8" y2="14"></line>
              <line x1="22" x2="16" y1="11" y2="11"></line>
            </svg>
          </div>
          <div class="card-content">
            <span class="card-value" :class="netChange >= 0 ? 'positive' : 'negative'">
              {{ netChange >= 0 ? '+' : '' }}{{ netChange.toFixed(1) }}
            </span>
            <span class="card-label">Net Değişim</span>
            <span class="card-count">Giriş - Çıkış</span>
          </div>
        </div>
      </div>
      
      <!-- Most Used Items -->
      <div class="section">
        <h3>
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 20V10"></path>
            <path d="M18 20V4"></path>
            <path d="M6 20v-4"></path>
          </svg>
          En Çok Kullanılan Ürünler
        </h3>
        
        <div v-if="report.most_used_items?.length > 0" class="top-products">
          <div 
            v-for="(item, index) in report.most_used_items" 
            :key="item.product_id"
            class="product-row"
          >
            <div class="rank" :class="getRankClass(index)">{{ index + 1 }}</div>
            <div class="product-info">
              <span class="product-name">{{ item.product_name }}</span>
              <span class="product-stats">{{ item.movement_count }} çıkış işlemi</span>
            </div>
            <div class="product-amount">
              <span class="amount-value">{{ item.total_out?.toFixed(1) }}</span>
              <span class="amount-label">toplam</span>
            </div>
            <div class="product-bar">
              <div 
                class="bar-fill" 
                :style="{ width: getBarWidth(item.total_out) + '%' }"
              ></div>
            </div>
          </div>
        </div>
        <div v-else class="empty-state">
          <p>Bu dönemde çıkış hareketi bulunmuyor</p>
        </div>
      </div>
      
      <!-- Movement Details -->
      <div class="section">
        <h3>
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 20h9"></path>
            <path d="M16.5 3.5a2.12 2.12 0 0 1 3 3L7 19l-4 1 1-4Z"></path>
          </svg>
          Hareket Detayları ({{ report.movements?.length || 0 }} hareket)
        </h3>
        
        <div v-if="report.movements?.length > 0" class="movements-table-container">
          <table class="movements-table">
            <thead>
              <tr>
                <th>Tarih/Saat</th>
                <th>Ürün</th>
                <th class="text-center">Tip</th>
                <th class="text-right">Miktar</th>
                <th>Açıklama</th>
              </tr>
            </thead>
            <tbody>
              <tr 
                v-for="movement in report.movements.slice(0, 20)" 
                :key="movement.id"
                :class="movement.movement_type"
              >
                <td class="date-cell">
                  {{ formatDateTime(movement.date) }}
                </td>
                <td>{{ movement.product_name }}</td>
                <td class="text-center">
                  <span :class="['type-badge', movement.movement_type]">
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
          <div v-if="report.movements.length > 20" class="more-info">
            ... ve {{ report.movements.length - 20 }} hareket daha
          </div>
        </div>
        <div v-else class="empty-state">
          <p>Bu dönemde hareket bulunmuyor</p>
        </div>
      </div>
    </div>
    
    <!-- No Report -->
    <div v-else class="no-report">
      <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
        <path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z"></path>
        <polyline points="14 2 14 8 20 8"></polyline>
        <line x1="16" x2="8" y1="13" y2="13"></line>
        <line x1="16" x2="8" y1="17" y2="17"></line>
        <line x1="10" x2="8" y1="9" y2="9"></line>
      </svg>
      <h3>Stok Raporu Oluşturun</h3>
      <p>Tarih seçerek günlük veya aylık stok hareketlerini görüntüleyin</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useStock } from '../composables/useStock'

const { stockReport, generateReport, loading } = useStock()

const activePeriod = ref('daily')
const selectedDate = ref(new Date().toISOString().split('T')[0])
const selectedMonth = ref(new Date().toISOString().slice(0, 7))
const report = ref(null)

// Net change calculation
const netChange = computed(() => {
  if (!report.value) return 0
  return (report.value.total_in || 0) - (report.value.total_out || 0)
})

// Maximum out amount (for bar chart)
const maxOut = computed(() => {
  if (!report.value?.most_used_items?.length) return 1
  return Math.max(...report.value.most_used_items.map(item => item.total_out))
})

// Calculate bar width
const getBarWidth = (value) => {
  return (value / maxOut.value) * 100
}

// Rank class
const getRankClass = (index) => {
  if (index === 0) return 'gold'
  if (index === 1) return 'silver'
  if (index === 2) return 'bronze'
  return ''
}

// Format date time
const formatDateTime = (dateStr) => {
  const date = new Date(dateStr)
  return date.toLocaleString('tr-TR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// Generate report
const generateReportHandler = async () => {
  const date = activePeriod.value === 'daily' ? selectedDate.value : selectedMonth.value
  const result = await generateReport(activePeriod.value, date)
  
  if (result.success) {
    report.value = result.report
  }
}

onMounted(() => {
  // Show today's report by default
  generateReportHandler()
})
</script>

<style scoped>
.reports-view {
  display: flex;
  flex-direction: column;
  height: 100%;
  gap: 1.5rem;
}

.report-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 1.5rem;
  background: var(--bg-card);
  border-radius: 0.75rem;
  border: 1px solid var(--border-color);
  flex-wrap: wrap;
  gap: 1rem;
}

.period-tabs {
  display: flex;
  gap: 0.5rem;
}

.period-tabs button {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.25rem;
  border: 1px solid var(--border-color);
  border-radius: 0.5rem;
  background: transparent;
  color: var(--text-muted);
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s;
}

.period-tabs button:hover {
  background: var(--bg-secondary);
  color: var(--text-primary);
}

.period-tabs button.active {
  background: var(--accent-color);
  border-color: var(--accent-color);
  color: white;
}

.date-selection {
  display: flex;
  gap: 0.75rem;
  align-items: center;
}

.date-selection input {
  padding: 0.75rem;
}

.report-content {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  flex: 1;
  overflow-y: auto;
}

/* Özet Kartları */
.summary-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1rem;
}

.summary-card {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.5rem;
  background: var(--bg-card);
  border-radius: 0.75rem;
  border: 1px solid var(--border-color);
}

.card-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 4rem;
  height: 4rem;
  border-radius: 0.75rem;
}

.summary-card.entry .card-icon {
  background: rgba(34, 197, 94, 0.1);
  color: #22c55e;
}

.summary-card.exit .card-icon {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
}

.summary-card.net .card-icon {
  background: rgba(59, 130, 246, 0.1);
  color: #3b82f6;
}

.card-content {
  display: flex;
  flex-direction: column;
}

.card-value {
  font-size: 2rem;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1;
}

.card-value.positive {
  color: #22c55e;
}

.card-value.negative {
  color: #ef4444;
}

.card-label {
  font-size: 0.875rem;
  color: var(--text-muted);
  margin-top: 0.25rem;
}

.card-count {
  font-size: 0.75rem;
  color: var(--text-muted);
}

/* Section */
.section {
  background: var(--bg-card);
  border-radius: 0.75rem;
  border: 1px solid var(--border-color);
  padding: 1.5rem;
}

.section h3 {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin: 0 0 1rem;
  font-size: 1rem;
  color: var(--text-primary);
}

.section h3 svg {
  color: var(--accent-color);
}

/* Top Products */
.top-products {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.product-row {
  display: grid;
  grid-template-columns: 40px 1fr 80px 100px;
  align-items: center;
  gap: 1rem;
  padding: 0.75rem;
  background: var(--bg-secondary);
  border-radius: 0.5rem;
}

.rank {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  border-radius: 50%;
  background: var(--border-color);
  font-weight: 600;
  font-size: 0.875rem;
  color: var(--text-muted);
}

.rank.gold {
  background: #fef3c7;
  color: #b45309;
}

.rank.silver {
  background: #f3f4f6;
  color: #4b5563;
}

.rank.bronze {
  background: #fef3c7;
  color: #92400e;
}

.product-info {
  display: flex;
  flex-direction: column;
}

.product-name {
  font-weight: 500;
  color: var(--text-primary);
}

.product-stats {
  font-size: 0.75rem;
  color: var(--text-muted);
}

.product-amount {
  text-align: right;
}

.amount-value {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-primary);
}

.amount-label {
  display: block;
  font-size: 0.7rem;
  color: var(--text-muted);
}

.product-bar {
  height: 0.5rem;
  background: var(--border-color);
  border-radius: 0.25rem;
  overflow: hidden;
}

.bar-fill {
  height: 100%;
  background: linear-gradient(90deg, var(--accent-color), #8b5cf6);
  border-radius: 0.25rem;
  transition: width 0.3s ease;
}

/* Movements Table */
.movements-table-container {
  overflow-x: auto;
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
}

.movements-table td {
  padding: 0.75rem 1rem;
  border-bottom: 1px solid var(--border-color);
  color: var(--text-primary);
}

.movements-table tbody tr:hover {
  background: var(--bg-secondary);
}

.movements-table tbody tr.in {
  background: rgba(34, 197, 94, 0.03);
}

.movements-table tbody tr.out {
  background: rgba(239, 68, 68, 0.03);
}

.text-center {
  text-align: center;
}

.text-right {
  text-align: right;
}

.type-badge {
  display: inline-block;
  padding: 0.25rem 0.5rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 500;
}

.type-badge.in {
  background: rgba(34, 197, 94, 0.1);
  color: #22c55e;
}

.type-badge.out {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
}

.amount-cell {
  font-weight: 600;
}

.text-success {
  color: #22c55e;
}

.text-danger {
  color: #ef4444;
}

.note-cell {
  max-width: 150px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.more-info {
  text-align: center;
  padding: 0.75rem;
  color: var(--text-muted);
  font-size: 0.875rem;
}

.empty-state {
  text-align: center;
  padding: 2rem;
  color: var(--text-muted);
}

/* No Report State */
.no-report {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  color: var(--text-muted);
  padding: 3rem;
}

.no-report svg {
  margin-bottom: 1rem;
  opacity: 0.5;
}

.no-report h3 {
  margin: 0 0 0.5rem;
  color: var(--text-primary);
}

.no-report p {
  margin: 0;
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

@media (max-width: 768px) {
  .report-header {
    flex-direction: column;
    align-items: stretch;
  }
  
  .period-tabs {
    justify-content: center;
  }
  
  .date-selection {
    justify-content: center;
  }
  
  .product-row {
    grid-template-columns: 32px 1fr 60px;
  }
  
  .product-bar {
    display: none;
  }
}
</style>
