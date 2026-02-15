<template>
  <div class="whatsapp-orders-container">
    <!-- Sol Panel: Sipariş Listesi -->
    <div class="orders-sidebar card">
      <!-- Header -->
      <div class="sidebar-header">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 bg-gradient-to-r from-green-500 to-green-600 rounded-xl flex items-center justify-center text-white text-lg">
            📱
          </div>
          <div>
            <h2 class="text-lg font-bold">{{ t('whatsapp.title') }}</h2>
            <span class="text-xs" style="color: var(--text-muted);">{{ t('whatsapp.orderCount', { count: orders.length }) }}</span>
          </div>
        </div>
        <button @click="handleNewOrder" class="btn btn-primary btn-sm">
          {{ t('whatsapp.new') }}
        </button>
      </div>

      <!-- Arama -->
      <div class="px-4 py-3 border-b" style="border-color: var(--border-color);">
        <div class="relative">
          <input 
            v-model="searchTerm"
            @input="handleSearch"
            type="text" 
            :placeholder="t('whatsapp.searchPlaceholder')"
            class="form-input w-full pl-10 text-sm"
          />
          <span class="absolute left-3 top-1/2 -translate-y-1/2 text-lg opacity-50">🔍</span>
        </div>
      </div>

      <!-- Durum Filtreleri -->
      <div class="px-4 py-2 border-b flex gap-2 flex-wrap" style="border-color: var(--border-color);">
        <button 
          v-for="(info, key) in statusFilters" 
          :key="key"
          @click="handleStatusFilter(key)"
          :class="['filter-chip', { active: statusFilter === key }]"
        >
          {{ info.label }}
        </button>
      </div>

      <!-- Sipariş Listesi -->
      <div class="orders-list">
        <div v-if="loading" class="text-center py-8">
          <div class="loading-spinner mx-auto"></div>
        </div>

        <div v-else-if="filteredOrders.length === 0" class="text-center py-12" style="color: var(--text-muted);">
          <div class="text-5xl mb-4 opacity-50">📭</div>
          <p>{{ t('whatsapp.noOrders') }}</p>
        </div>

        <div 
          v-else
          v-for="order in filteredOrders" 
          :key="order.id"
          @click="handleSelectOrder(order)"
          :class="['order-card', { active: currentOrder?.id === order.id }]"
        >
          <!-- Durum Badge -->
          <div class="flex items-center justify-between mb-2">
            <span :class="['status-badge', `status-${order.current_status || 'none'}`]">
              {{ getStatusInfo(order.current_status).emoji }} {{ getStatusInfo(order.current_status).label }}
            </span>
            <span class="text-xs" style="color: var(--text-muted);">{{ formatDate(order.date) }}</span>
          </div>

          <!-- Müşteri Bilgisi -->
          <div class="font-semibold mb-1">
            {{ order.customer_name || t('whatsapp.anonymousCustomer') }}
          </div>
          <div v-if="order.customer_phone" class="text-sm mb-2" style="color: var(--text-muted);">
            📱 {{ order.customer_phone }}
          </div>

          <!-- Özet -->
          <div class="flex justify-between items-center pt-2 border-t" style="border-color: var(--border-color);">
            <span class="text-xs" style="color: var(--text-muted);">
              📦 {{ t('whatsapp.itemCount', { count: order.items?.length || 0 }) }}
            </span>
            <span class="font-bold text-success">
              {{ formatCurrency(order.grand_total) }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Sağ Panel: Detay/Form -->
    <div class="orders-content card">
      <!-- Sipariş Seçilmemiş -->
      <div v-if="!currentOrder" class="flex flex-col items-center justify-center h-full" style="color: var(--text-muted);">
        <div class="text-6xl mb-4 opacity-30">📋</div>
        <p class="text-lg mb-2">{{ t('whatsapp.selectOrder') }}</p>
        <p class="text-sm">{{ t('whatsapp.or') }}</p>
        <button @click="handleNewOrder" class="btn btn-primary mt-4">
          {{ t('whatsapp.createOrder') }}
        </button>
      </div>

      <!-- Sipariş Detay/Form -->
      <template v-else>
        <!-- Header -->
        <div class="content-header">
          <div class="flex items-center gap-3">
            <div class="w-12 h-12 bg-gradient-to-r from-accent to-purple-500 rounded-xl flex items-center justify-center text-white text-xl">
              {{ currentOrder.id ? '📋' : '✨' }}
            </div>
            <div>
              <h3 class="text-lg font-bold">
                {{ currentOrder.id ? t('whatsapp.orderDetail') : t('whatsapp.newOrder') }}
              </h3>
              <span v-if="currentOrder.id" class="text-xs" style="color: var(--text-muted);">
                ID: {{ currentOrder.id.substring(0, 8) }}...
              </span>
            </div>
          </div>
          <div class="flex gap-2">
            <button v-if="currentOrder.id" @click="showStatusModal = true" class="btn btn-secondary btn-sm">
              📍 {{ t('whatsapp.addStatus') }}
            </button>
            <button @click="handleSave" :disabled="!hasUnsavedChanges && currentOrder.id" class="btn btn-primary btn-sm">
              💾 {{ t('whatsapp.save') }}
            </button>
          </div>
        </div>

        <!-- Form İçeriği -->
        <div class="content-body">
          <!-- Müşteri Bilgileri -->
          <div class="form-section">
            <h4 class="section-title">👤 {{ t('whatsapp.customerInfo') }}</h4>
            <div class="form-row">
              <div class="form-group">
                <label class="form-label">{{ t('whatsapp.dateRequired') }}</label>
                <input 
                  v-model="currentOrder.date" 
                  @change="markAsModified"
                  type="date" 
                  class="form-input"
                />
              </div>
              <div class="form-group">
                <label class="form-label">{{ t('whatsapp.customerName') }}</label>
                <input 
                  v-model="currentOrder.customer_name" 
                  @input="markAsModified"
                  type="text" 
                  :placeholder="t('whatsapp.customerNamePlaceholder')"
                  class="form-input"
                />
              </div>
              <div class="form-group">
                <label class="form-label">{{ t('whatsapp.phone') }}</label>
                <input 
                  v-model="currentOrder.customer_phone" 
                  @input="markAsModified"
                  type="text" 
                  :placeholder="t('whatsapp.phonePlaceholder')"
                  class="form-input"
                />
              </div>
            </div>
          </div>

          <!-- Ödeme Bilgileri -->
          <div class="form-section">
            <h4 class="section-title">💳 {{ t('whatsapp.paymentInfo') }}</h4>
            <div class="form-row">
              <div class="form-group">
                <label class="form-label">{{ t('whatsapp.paymentMethod') }}</label>
                <div class="payment-options">
                  <label 
                    v-for="(info, key) in PAYMENT_METHODS" 
                    :key="key"
                    :class="['payment-option', { active: currentOrder.payment_method === key }]"
                  >
                    <input 
                      type="radio" 
                      :value="key" 
                      v-model="currentOrder.payment_method"
                      @change="markAsModified"
                      class="hidden"
                    />
                    <span class="text-lg">{{ info.emoji }}</span>
                    <span>{{ info.label }}</span>
                  </label>
                </div>
              </div>
              <div v-if="currentOrder.payment_method === 'kredi_karti'" class="form-group">
                <label class="form-label">{{ t('whatsapp.paymentNote') }}</label>
                <input 
                  v-model="currentOrder.payment_note" 
                  @input="markAsModified"
                  type="text" 
                  :placeholder="t('whatsapp.paymentNotePlaceholder')"
                  class="form-input"
                />
              </div>
            </div>
          </div>

          <!-- Sipariş Kalemleri -->
          <div class="form-section">
            <div class="flex items-center justify-between mb-4">
              <h4 class="section-title mb-0">📦 {{ t('whatsapp.orderItems') }}</h4>
              <button @click="addItem" class="btn btn-secondary btn-sm">
                {{ t('whatsapp.addItem') }}
              </button>
            </div>

            <div v-if="currentOrder.items.length === 0" class="text-center py-8 border-2 border-dashed rounded-xl" style="border-color: var(--border-color); color: var(--text-muted);">
              <p>{{ t('whatsapp.noItems') }}</p>
              <button @click="addItem" class="btn btn-primary btn-sm mt-3">
                {{ t('whatsapp.addFirstItem') }}
              </button>
            </div>

            <table v-else class="items-table">
              <thead>
                <tr>
                  <th class="w-2/5">{{ t('whatsapp.productName') }}</th>
                  <th class="w-1/6">{{ t('whatsapp.type') }}</th>
                  <th class="w-1/8 text-center">{{ t('whatsapp.quantity') }}</th>
                  <th class="w-1/6 text-right">{{ t('whatsapp.price') }}</th>
                  <th class="w-1/6 text-right">{{ t('whatsapp.total') }}</th>
                  <th class="w-10"></th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(item, index) in currentOrder.items" :key="index">
                  <td>
                    <input 
                      :value="item.product_name"
                      @input="updateItem(index, 'product_name', $event.target.value)"
                      type="text" 
                      :placeholder="t('whatsapp.productNamePlaceholder')"
                      class="form-input text-sm"
                    />
                  </td>
                  <td>
                    <select 
                      :value="item.type"
                      @change="updateItem(index, 'type', $event.target.value)"
                      class="form-input text-sm"
                    >
                      <option value="cikma">🔵 {{ t('whatsapp.partTypes.cikma') }}</option>
                      <option value="sifir">🟢 {{ t('whatsapp.partTypes.sifir') }}</option>
                    </select>
                  </td>
                  <td>
                    <input 
                      :value="item.quantity"
                      @input="updateItem(index, 'quantity', parseInt($event.target.value) || 0)"
                      type="number" 
                      min="1"
                      class="form-input text-sm text-center"
                    />
                  </td>
                  <td>
                    <input 
                      :value="item.price"
                      @input="updateItem(index, 'price', parseFloat($event.target.value) || 0)"
                      type="number" 
                      min="0"
                      step="0.01"
                      placeholder="₺"
                      class="form-input text-sm text-right"
                    />
                  </td>
                  <td class="text-right font-semibold">
                    {{ formatCurrency(item.total) }}
                  </td>
                  <td class="text-center">
                    <button @click="removeItem(index)" class="text-danger hover:opacity-70 transition-opacity">
                      🗑️
                    </button>
                  </td>
                </tr>
              </tbody>
              <tfoot>
                <tr>
                  <td colspan="4" class="text-right font-bold">{{ t('whatsapp.grandTotal') }}</td>
                  <td class="text-right font-bold text-lg text-success">
                    {{ formatCurrency(currentOrder.grand_total) }}
                  </td>
                  <td></td>
                </tr>
              </tfoot>
            </table>
          </div>

          <!-- Durum Geçmişi (sadece kayıtlı siparişler için) -->
          <div v-if="currentOrder.id && currentOrder.status_history?.length > 0" class="form-section">
            <h4 class="section-title">📍 {{ t('whatsapp.statusHistory') }}</h4>
            <div class="status-timeline">
              <div 
                v-for="status in [...currentOrder.status_history].reverse()" 
                :key="status.id"
                class="timeline-item"
              >
                <div class="timeline-dot" :class="`status-${status.status}`"></div>
                <div class="timeline-content">
                  <div class="flex items-center gap-2">
                    <span :class="['status-badge', `status-${status.status}`]">
                      {{ getStatusInfo(status.status).emoji }} {{ getStatusInfo(status.status).label }}
                    </span>
                    <span class="text-xs" style="color: var(--text-muted);">
                      {{ formatDateTime(status.timestamp) }}
                    </span>
                  </div>
                  <p v-if="status.note" class="text-sm mt-1" style="color: var(--text-muted);">
                    {{ status.note }}
                  </p>
                </div>
              </div>
            </div>
          </div>

          <!-- Aksiyon Butonları -->
          <div v-if="currentOrder.id" class="form-section">
            <div class="flex gap-3">
              <button @click="handleSave" :disabled="!hasUnsavedChanges" class="btn btn-primary">
                💾 {{ t('whatsapp.save') }}
              </button>
              <button @click="handleDelete" class="btn btn-danger">
                🗑️ {{ t('whatsapp.deleteOrder') }}
              </button>
              <button @click="clearSelection" class="btn btn-secondary">
                ✕ {{ t('whatsapp.close') }}
              </button>
            </div>
          </div>
        </div>
      </template>
    </div>

    <!-- Durum Ekleme Modalı -->
    <WhatsAppOrderStatusModal 
      v-if="showStatusModal"
      :order-id="currentOrder?.id"
      @close="showStatusModal = false"
      @saved="handleStatusSaved"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useWhatsAppOrders, ORDER_STATUSES, PAYMENT_METHODS } from '@/composables/useWhatsAppOrders'
import WhatsAppOrderStatusModal from './WhatsAppOrderStatusModal.vue'
import { useI18n } from '@/i18n'
import { useToast } from '@/composables/useToast'

const { t } = useI18n()
const { showToast } = useToast()

const {
  orders,
  currentOrder,
  loading,
  searchTerm,
  statusFilter,
  filteredOrders,
  hasUnsavedChanges,
  loadOrders,
  saveOrder,
  deleteOrder,
  newOrder,
  selectOrder,
  clearSelection,
  addItem,
  removeItem,
  updateItem,
  markAsModified,
  formatDate,
  formatDateTime,
  formatCurrency,
  getStatusInfo,
  filterByStatus
} = useWhatsAppOrders()

const showStatusModal = ref(false)

const statusFilters = {
  all: { label: t('whatsapp.filters.all') },
  tedarik_surecinde: { label: t('whatsapp.filters.tedarik') },
  hazirlaniyor: { label: t('whatsapp.filters.hazirlaniyor') },
  tamamlandi: { label: t('whatsapp.filters.tamamlandi') },
  iade: { label: t('whatsapp.filters.iade') }
}

onMounted(async () => {
  await loadOrders()
})

async function handleSearch() {
  await loadOrders(searchTerm.value)
}

async function handleStatusFilter(status) {
  await filterByStatus(status)
}

function handleNewOrder() {
  newOrder()
}

function handleSelectOrder(order) {
  selectOrder(order)
}

async function handleSave() {
  const result = await saveOrder()
  if (result.success) {
    showToast(t('whatsapp.orderSaved'), 'success')
  } else if (result.error) {
    showToast(t('whatsapp.saveError', { error: result.error }), 'error')
  }
}

async function handleDelete() {
  if (!currentOrder.value?.id) return
  
  if (confirm(t('whatsapp.deleteConfirm'))) {
    const result = await deleteOrder(currentOrder.value.id)
    if (result.success) {
      showToast(t('whatsapp.orderDeleted'), 'success')
    }
  }
}

async function handleStatusSaved() {
  showStatusModal.value = false
  // Sipariş listesini ve detayı yenile
  await loadOrders()
}
</script>

<style scoped>
.whatsapp-orders-container {
  display: grid;
  grid-template-columns: 380px 1fr;
  gap: 1.5rem;
  height: calc(100vh - 140px);
}

.orders-sidebar {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.sidebar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 1.25rem;
  border-bottom: 1px solid var(--border-color);
  background: linear-gradient(to right, rgba(34, 197, 94, 0.1), rgba(22, 163, 74, 0.05));
}

.orders-list {
  flex: 1;
  overflow-y: auto;
  padding: 0.75rem;
}

.order-card {
  background: var(--bg-secondary);
  border: 2px solid var(--border-color);
  border-radius: 0.75rem;
  padding: 1rem;
  margin-bottom: 0.75rem;
  cursor: pointer;
  transition: all 0.2s ease;
}

.order-card:hover {
  border-color: var(--accent-color);
  transform: translateX(4px);
}

.order-card.active {
  border-color: var(--accent-color);
  background: rgba(14, 165, 233, 0.1);
}

.filter-chip {
  padding: 0.35rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 600;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  cursor: pointer;
  transition: all 0.2s ease;
}

.filter-chip:hover {
  border-color: var(--accent-color);
}

.filter-chip.active {
  background: var(--accent-color);
  border-color: var(--accent-color);
  color: white;
}

.orders-content {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.content-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 1.5rem;
  border-bottom: 1px solid var(--border-color);
  background: linear-gradient(to right, rgba(14, 165, 233, 0.1), rgba(168, 85, 247, 0.05));
}

.content-body {
  flex: 1;
  overflow-y: auto;
  padding: 1.5rem;
}

.form-section {
  margin-bottom: 2rem;
}

.section-title {
  font-size: 1rem;
  font-weight: 700;
  margin-bottom: 1rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.form-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-label {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--text-muted);
}

.payment-options {
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.payment-option {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1rem;
  border: 2px solid var(--border-color);
  border-radius: 0.75rem;
  cursor: pointer;
  transition: all 0.2s ease;
  background: var(--bg-secondary);
}

.payment-option:hover {
  border-color: var(--accent-color);
}

.payment-option.active {
  border-color: var(--accent-color);
  background: rgba(14, 165, 233, 0.1);
}

.items-table {
  width: 100%;
  border-collapse: collapse;
}

.items-table th {
  padding: 0.75rem;
  text-align: left;
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: uppercase;
  color: var(--text-muted);
  border-bottom: 2px solid var(--border-color);
}

.items-table td {
  padding: 0.5rem;
  border-bottom: 1px solid var(--border-color);
}

.items-table tfoot td {
  padding: 1rem 0.5rem;
  border-bottom: none;
  border-top: 2px solid var(--border-color);
}

/* Status Badge Styles */
.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 600;
}

.status-badge.status-tedarik_surecinde,
.status-tedarik_surecinde { background: rgba(14, 165, 233, 0.15); color: #0ea5e9; }

.status-badge.status-hazirlaniyor,
.status-hazirlaniyor { background: rgba(245, 158, 11, 0.15); color: #f59e0b; }

.status-badge.status-hazirlandi,
.status-hazirlandi { background: rgba(34, 197, 94, 0.15); color: #22c55e; }

.status-badge.status-kargoya_verildi,
.status-kargoya_verildi { background: rgba(168, 85, 247, 0.15); color: #a855f7; }

.status-badge.status-tamamlandi,
.status-tamamlandi { background: rgba(34, 197, 94, 0.2); color: #16a34a; }

.status-badge.status-iade,
.status-iade { background: rgba(239, 68, 68, 0.15); color: #ef4444; }

.status-badge.status-hukumsuz,
.status-hukumsuz { background: rgba(100, 116, 139, 0.15); color: #64748b; }

.status-badge.status-none,
.status-none { background: rgba(100, 116, 139, 0.1); color: #94a3b8; }

/* Status Timeline */
.status-timeline {
  position: relative;
  padding-left: 1.5rem;
}

.timeline-item {
  position: relative;
  padding-bottom: 1.5rem;
  padding-left: 1.5rem;
  border-left: 2px solid var(--border-color);
}

.timeline-item:last-child {
  border-left-color: transparent;
  padding-bottom: 0;
}

.timeline-dot {
  position: absolute;
  left: -0.5rem;
  top: 0;
  width: 1rem;
  height: 1rem;
  border-radius: 50%;
  border: 2px solid var(--bg-card);
}

.timeline-dot.status-tedarik_surecinde { background: #0ea5e9; }
.timeline-dot.status-hazirlaniyor { background: #f59e0b; }
.timeline-dot.status-hazirlandi { background: #22c55e; }
.timeline-dot.status-kargoya_verildi { background: #a855f7; }
.timeline-dot.status-tamamlandi { background: #16a34a; }
.timeline-dot.status-iade { background: #ef4444; }
.timeline-dot.status-hukumsuz { background: #64748b; }

.timeline-content {
  padding-top: 0;
}

.loading-spinner {
  width: 2rem;
  height: 2rem;
  border: 3px solid var(--border-color);
  border-top-color: var(--accent-color);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
