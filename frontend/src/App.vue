<template>
  <div class="min-h-screen flex flex-col transition-colors duration-300" style="background: var(--bg-primary); color: var(--text-primary);">

    <!-- Loadout overlay -->
    <div v-if="initialLoading" class="fixed inset-0 z-[2000] flex items-center justify-center bg-black/90">
      <div class="bg-card p-6 rounded-lg text-center shadow-lg flex flex-col items-center gap-4">
        <div class="w-12 h-12 border-4 border-accent rounded-full border-t-transparent animate-spin"></div>
        <div class="text-2xl font-bold">{{ t('app.loading') }}</div>
        <div class="text-sm text-muted">{{ t('common.loading') }}</div>
      </div>
    </div>

    <!-- Main application content is rendered only after initial loading completes -->
    <div v-if="!initialLoading">
      <AppHeader 
        :activeTab="activeTab"
        @showHistory="modals.showHistory()" 
        @newOrder="handleNewOrder"
        @showCatalog="showCatalog = true"
        @showSettings="showSettings = true"
      >
      <template #critical-stock>
        <CriticalStockBadge 
          @stock-in="handleCriticalStockIn" 
          @show-all="handleShowAllCriticalStock"
        />
      </template>
    </AppHeader>
    
    <!-- Tab Navigation -->
    <nav class="tab-nav" style="background: var(--bg-secondary); border-bottom: 1px solid var(--border-color);">
      <div class="max-w-[1800px] mx-auto w-full px-6">
        <div class="tab-list">
          <button 
            @click="activeTab = 'orders'" 
            :class="['tab-btn', { active: activeTab === 'orders' }]"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M16 16h6"></path>
              <path d="M19 13v6"></path>
              <path d="M21 10V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l2-1.14"></path>
              <path d="M16.5 9.4 7.55 4.24"></path>
              <polyline points="3.29 7 12 12 20.71 7"></polyline>
              <line x1="12" y1="22" x2="12" y2="12"></line>
            </svg>
            {{ t('tabs.orders') }}
          </button>
          <button 
            @click="activeTab = 'stock'" 
            :class="['tab-btn', { active: activeTab === 'stock' }]"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="m7.5 4.27 9 5.15"></path>
              <path d="M21 8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16Z"></path>
              <path d="m3.3 7 8.7 5 8.7-5"></path>
              <path d="M12 22V12"></path>
            </svg>
            {{ t('tabs.stock') }}
          </button>
          <button 
            @click="activeTab = 'reports'" 
            :class="['tab-btn', { active: activeTab === 'reports' }]"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M3 3v18h18"></path>
              <path d="m19 9-5 5-4-4-3 3"></path>
            </svg>
            {{ t('tabs.reports') }}
          </button>
          <!-- Geliştirici Sekmesi (sadece developerMode aktifse) -->
          <button 
            v-if="developerMode"
            @click="activeTab = 'developer'" 
            :class="['tab-btn developer-tab', { active: activeTab === 'developer' }]"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="16 18 22 12 16 6"/>
              <polyline points="8 6 2 12 8 18"/>
            </svg>
            {{ t('settings.developer') }}
          </button>
          <button 
            @click="activeTab = 'notes'" 
            :class="['tab-btn', { active: activeTab === 'notes' }]"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z"></path>
              <polyline points="14 2 14 8 20 8"></polyline>
              <line x1="16" y1="13" x2="8" y2="13"></line>
              <line x1="16" y1="17" x2="8" y2="17"></line>
              <line x1="10" y1="9" x2="8" y2="9"></line>
            </svg>
            {{ t('tabs.notes') }}
          </button>
          <button 
            @click="activeTab = 'whatsapp'" 
            :class="['tab-btn', { active: activeTab === 'whatsapp' }]"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"></path>
            </svg>
            {{ t('tabs.whatsapp') }}
          </button>
        </div>
      </div>
    </nav>
    
    <!-- Siparişler View -->
    <template v-if="activeTab === 'orders'">
      <OrderBar />
      
      <main class="flex-1 grid grid-cols-4 gap-6 p-6 max-w-[1800px] mx-auto w-full">
        <!-- Left: Product Form -->
        <div class="col-span-1">
          <ProductForm @clearAll="handleClearAll" />
        </div>
        
        <!-- Center: Order Items -->
        <div class="col-span-2">
          <OrderItems 
            @newOrder="handleNewOrder"
            @deleteProduct="handleDeleteProduct"
          />
        </div>
        
        <!-- Right: Summary + Orders List -->
        <div class="col-span-1 flex flex-col gap-4">
          <!-- Order Summary -->
          <OrderSummary
            @saveOrder="handleSaveOrder"
            @saveAsNew="handleSaveAsNew"
            @exportTxt="exportTxt"
            @exportPng="exportPng"
            @exportTxtWhatsApp="exportTxtWhatsApp"
            @exportPngWhatsApp="exportPngWhatsApp"
          />
          
          <!-- Orders List -->
          <OrdersList 
            ref="ordersListRef"
            :refreshTrigger="refreshTrigger"
            :advancedSearchResults="advancedResults"
            @loadOrder="handleLoadOrder"
            @deleteOrder="handleDeleteOrder"
            @showAdvancedSearch="modals.showAdvancedSearch()"
            @clearAdvancedSearch="clearAdvancedResults"
          />
        </div>
      </main>
    </template>
    
    <!-- Stok Yönetimi View -->
    <template v-else-if="activeTab === 'stock'">
      <main class="flex-1 p-6 max-w-[1800px] mx-auto w-full">
        <StockList
          ref="stockListRef"
          :filter-critical="showOnlyCriticalStock"
          @stock-in="handleStockIn"
          @stock-out="handleStockOut"
          @edit-product="handleEditProduct"
          @delete-product="handleDeleteStockProduct"
          @movements="handleMovements"
          @new-product="handleNewProduct"
          @stokGiris="openBulkStockEntry"
          @stokCikis="openBulkStockExit"
          @yeniUrun="handleNewProduct"
          @bulk-stock-in="handleBulkStockIn"
          @bulk-stock-out="handleBulkStockOut"
          @bulk-edit="handleBulkEdit"
        />
      </main>
    </template>
    
    <!-- Raporlar View -->
    <template v-else-if="activeTab === 'reports'">
      <main class="flex-1 p-6 max-w-[1800px] mx-auto w-full">
        <ReportsView />
      </main>
    </template>
    
    <!-- Notlar View -->
    <template v-else-if="activeTab === 'notes'">
      <main class="flex-1 p-6 max-w-[1800px] mx-auto w-full">
        <NotesView />
      </main>
    </template>
    
    <!-- WhatsApp Siparişleri View -->
    <template v-else-if="activeTab === 'whatsapp'">
      <main class="flex-1 p-6 max-w-[1800px] mx-auto w-full">
        <WhatsAppOrdersView />
      </main>
    </template>
    
    <!-- Geliştirici View -->
    <template v-else-if="activeTab === 'developer' && developerMode">
      <main class="flex-1 p-6 max-w-[1800px] mx-auto w-full">
        <DeveloperConsole />
      </main>
    </template>
    
    <!-- Footer -->
    <footer class="border-t py-3 px-6 text-center text-sm" style="background: var(--bg-secondary); border-color: var(--border-color); color: var(--text-muted);">
      <div class="flex items-center justify-center gap-2 flex-wrap">
        <span class="opacity-70">🔓 {{ t('settings.openSource') }}</span>
        <span class="opacity-50">•</span>
        <span>{{ t('common.version') }}:</span>
        <span class="font-semibold text-success">{{ appVersion }}</span>
        <span class="opacity-50">•</span>
        <span>{{ t('settings.developer') }}:</span>
        <a 
          href="https://github.com/Ferhatduran55" 
          target="_blank" 
          rel="noopener noreferrer"
          class="font-semibold text-accent hover:underline cursor-pointer"
        >
          Ferhat Duran
        </a>
        <span class="opacity-50">•</span>
        <a 
          href="https://github.com/Ferhatduran55/go-auto-aom" 
          target="_blank" 
          rel="noopener noreferrer"
          class="text-accent hover:underline cursor-pointer flex items-center gap-1"
        >
          <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>
          GitHub
        </a>
      </div>
      <div class="text-xs opacity-50 mt-1">
        {{ t('settings.copyright', { year: new Date().getFullYear() }) }}
      </div>
    </footer>
    
    <Modals 
      ref="modals" 
      @loadOrder="handleLoadOrder"
      @advancedSearchResults="handleAdvancedSearchResults" 
    />
    
    <CatalogModal 
      :visible="showCatalog" 
      @close="showCatalog = false"
      @updated="handleCatalogUpdated"
    />
    
    <!-- Stok Modalleri -->
    <StockEntryModal
      v-if="showStockEntry"
      :product="selectedProduct"
      :products="stockProducts"
      :bulkProducts="bulkProducts"
      @close="closeStockModals"
      @success="handleStockOperationSuccess"
    />
    
    <StockExitModal
      v-if="showStockExit"
      :product="selectedProduct"
      :products="stockProducts"
      :bulkProducts="bulkProducts"
      @close="closeStockModals"
      @success="handleStockOperationSuccess"
    />
    
    <StockMovementsModal
      v-if="showStockMovements"
      :product="selectedProduct"
      @close="showStockMovements = false"
    />
    
    <ProductFormModal
      v-if="showProductForm"
      :product="editingProduct"
      @close="showProductForm = false"
      @success="handleProductFormSuccess"
    />
    
    <BulkEditModal
      v-if="showBulkEdit"
      :products="bulkEditProducts"
      @close="closeBulkEditModal"
      @success="handleBulkEditSuccess"
    />
    
    <Toast />
    
    <SettingsModal
      v-if="showSettings"
      @close="showSettings = false"
    />
  </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useOrder } from '@/composables/useOrder'
import { useToast } from '@/composables/useToast'
import { useStock } from '@/composables/useStock'
import { useSettings } from '@/composables/useSettings'
import AppHeader from '@/components/AppHeader.vue'
import OrderBar from '@/components/OrderBar.vue'
import ProductForm from '@/components/ProductForm.vue'
import OrderItems from '@/components/OrderItems.vue'
import OrdersList from '@/components/OrdersList.vue'
import OrderSummary from '@/components/OrderSummary.vue'
import Modals from '@/components/Modals.vue'
import CatalogModal from '@/components/CatalogModal.vue'
import Toast from '@/components/Toast.vue'
import SettingsModal from '@/components/SettingsModal.vue'
// Stok components
import StockList from '@/components/StockList.vue'
import StockEntryModal from '@/components/StockEntryModal.vue'
import StockExitModal from '@/components/StockExitModal.vue'
import StockMovementsModal from '@/components/StockMovementsModal.vue'
import ProductFormModal from '@/components/ProductFormModal.vue'
import DeveloperConsole from '@/components/DeveloperConsole.vue'
import BulkEditModal from '@/components/BulkEditModal.vue'
import CriticalStockBadge from '@/components/CriticalStockBadge.vue'
import ReportsView from '@/components/ReportsView.vue'
import NotesView from '@/components/NotesView.vue'
import WhatsAppOrdersView from '@/components/WhatsAppOrdersView.vue'
import { useI18n } from '@/i18n'

const { t, locale } = useI18n()

const appVersion = ref('')

onMounted(async () => {
  try {
    appVersion.value = await api.getAppVersion()
  } catch (e) {
    // If API fails, try to use injected version directly or fallback
    appVersion.value = typeof __APP_VERSION__ !== 'undefined' ? __APP_VERSION__ : '0.0.0'
  }
})
const showCatalog = ref(false)
const showSettings = ref(false)
const activeTab = ref('orders') // orders | stock | reports | notes

const { 
  products, 
  orderTitle, 
  customerName,
  customerPhone,
  grandTotal,
  loadData, 
  loadOrder, 
  saveOrder, 
  resetOrder, 
  deleteProduct, 
  clearProducts,
  currentOrderId,
  isEditing,
  hasUnsavedChanges
} = useOrder()

// Stock state and helpers
const { products: stockProducts, loadProducts: loadStockProducts, deleteProduct: deleteStockProduct, loadStockMovements } = useStock()

const { settings } = useSettings()
const currency = computed(() => (settings && settings.value && settings.value.currency) || '₺')
const initialLoading = ref(true)
const developerMode = ref(false)
const showProductForm = ref(false)
const showBulkEdit = ref(false)
const showOnlyCriticalStock = ref(false)
const showStockEntry = ref(false)
const showStockExit = ref(false)
const showStockMovements = ref(false)
const selectedProduct = ref(null)
const editingProduct = ref(null)
const bulkEditProducts = ref([])
const stockListRef = ref(null)

// Toast helper
const { showToast } = useToast()

const modals = ref(null)
const ordersListRef = ref(null)
const refreshTrigger = ref(0)

// Banner için fiyat formatlama
function formatPriceBanner(n) {
  return (n || 0).toLocaleString('tr-TR', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

// Handlers
async function handleSaveOrder() {
  if (products.value.length === 0) {
    showToast(t('app.toasts.listEmpty'), 'error')
    return
  }
  
  const result = await saveOrder()
  if (result.error) {
    showToast(result.error, 'error')
  } else {
    showToast(t('app.toasts.orderSaved'))
    refreshTrigger.value++
  }
}

async function handleSaveAsNew() {
  currentOrderId.value = null
  await handleSaveOrder()
}

async function handleLoadOrder(id) {
  // Aynı siparişe tıklandıysa hiçbir şey yapma
  if (currentOrderId.value === id) {
    return
  }
  
  // Kaydedilmemiş değişiklik varsa uyar
  if (hasUnsavedChanges.value) {
    modals.value.showConfirm(
      t('app.confirms.unsavedChangesTitle'),
      t('app.confirms.unsavedChangesLoad'),
      async () => {
        const result = await loadOrder(id)
        if (result.error) {
          showToast(t('app.toasts.orderLoadError'), 'error')
        } else {
          showToast(t('app.toasts.orderLoaded'))
        }
      },
      '⚠️'
    )
    return
  }
  
  const result = await loadOrder(id)
  if (result.error) {
    showToast(t('orders.loadError'), 'error')
  } else {
    showToast(t('orders.loaded'))
  }
}

function handleDeleteProduct(id) {
  modals.value.showConfirm(
    t('products.deleteProduct'),
    t('products.deleteProductConfirm'),
    () => {
      deleteProduct(id)
      showToast(t('products.productRemoved'))
    },
    '🗑️'
  )
}

function handleDeleteOrder(id) {
  modals.value.showConfirm(
    t('orders.deleteOrder'),
    t('orders.deleteConfirm'),
    async () => {
      const { api } = await import('@/api')
      const result = await api.deleteOrder(id)
      if (!result.error) {
        showToast(t('orders.deleted'))
        if (currentOrderId.value === id) {
          resetOrder()
        }
        refreshTrigger.value++
      }
    },
    '🗑️'
  )
}

function handleNewOrder() {
  // Kaydedilmemiş değişiklik varsa uyar
  if (hasUnsavedChanges.value) {
    modals.value.showConfirm(
      t('app.confirms.unsavedChangesTitle'),
      t('app.confirms.unsavedChangesNew'),
      () => {
        activeTab.value = 'orders'
        resetOrder()
        showToast(t('app.toasts.newOrderStarted'))
      },
      '⚠️'
    )
    return
  }
  
  // Switch to orders tab first if not already there
  activeTab.value = 'orders'
  
  // Reset order form
  resetOrder()
  showToast(t('app.toasts.newOrderStarted'))
}

function handleClearAll() {
  if (products.value.length === 0) return
  modals.value.showConfirm(
    t('app.confirms.clearListTitle'),
    t('app.confirms.clearListMessage'),
    () => {
      clearProducts()
      showToast(t('app.toasts.listCleared'))
    },
    '🗑️'
  )
}

function handleCatalogUpdated() {
  loadData()
  refreshTrigger.value++
}

function handleAdvancedSearchResults(orders) {
  advancedResults.value = orders
}

function clearAdvancedResults() {
  advancedResults.value = null
}

const advancedResults = ref(null)

// Export functions
function exportTxt() {
  if (products.value.length === 0) {
    showToast(t('app.toasts.listEmpty'), 'error')
    return
  }

  // WhatsApp uyumlu format
  let content = `📋 *SİPARİŞ: ${orderTitle.value || 'İsimsiz'}*\n`
  content += `👤 *Müşteri:* ${customerName.value || '-'}\n`
  const loc = locale.value || navigator.language || 'en-US'
  content += `📅 *Tarih:* ${new Date().toLocaleDateString(loc)}\n`
  content += `\n${'─'.repeat(35)}\n\n`

  products.value.forEach((p, i) => {
    const durum = p.part_status === 'original' ? '🟢 Orijinal' : (p.part_status === 'zero' ? '🔵 Sıfır' : '🟡 Çıkma')
    content += `*${i + 1}. ${p.product_name}*\n`
    content += `   🔢 OEM: ${p.oem_number}\n`
    content += `   📦 Adet: ${p.quantity}\n`
    content += `   ${durum}\n`
    content += `   💰 *${currency.value}${formatPrice(p.total_price)}*\n\n`
  })

  content += `${'─'.repeat(35)}\n`
  content += `\n💵 *GENEL TOPLAM: ${currency.value}${formatPrice(grandTotal.value)}*\n`

  const blob = new Blob([content], { type: 'text/plain;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `siparis_${orderTitle.value || 'isimsiz'}_${Date.now()}.txt`
  a.click()
  URL.revokeObjectURL(url)
  showToast(t('export.txtDownloaded'))
}

function exportTxtWhatsApp() {
  if (products.value.length === 0) {
    showToast(t('app.toasts.listEmpty'), 'error')
    return
  }

  // WhatsApp uyumlu format - Türkçe karakterli
  let content = `*SİPARİŞ: ${orderTitle.value || 'İsimsiz'}*\n`
  content += `*Müşteri:* ${customerName.value || '-'}\n`
  content += `*Tarih:* ${new Date().toLocaleDateString(locale.value || navigator.language || 'en-US')}\n`
  content += `\n${'─'.repeat(30)}\n\n`

  products.value.forEach((p, i) => {
    const durum = p.part_status === 'original' ? '[Orijinal]' : (p.part_status === 'zero' ? '[Sıfır]' : '[Çıkma]')
    content += `*${i + 1}. ${p.product_name}*\n`
    content += `   OEM: ${p.oem_number}\n`
    content += `   Adet: ${p.quantity}\n`
    content += `   ${durum}\n`
    content += `   *${currency.value}${formatPrice(p.total_price)}*\n\n`
  })

  content += `${'─'.repeat(30)}\n`
  content += `\n*GENEL TOPLAM: ${currency.value}${formatPrice(grandTotal.value)}*\n`

  // Telefon numarasını formatla
  if (!customerPhone.value) {
    showToast(t('app.toasts.customerPhoneNotFound'), 'error')
    return
  }
  
  const phone = customerPhone.value.replace(/\D/g, '')
  const encodedText = encodeURIComponent(content)
  
  // WhatsApp Web linkini aç - mesajı taslak olarak hazırlar
  window.open(`https://wa.me/${phone}?text=${encodedText}`, '_blank')
  showToast(t('export.whatsappOpening'))
}

function exportPng() {
  if (products.value.length === 0) {
    showToast(t('app.toasts.listEmpty'), 'error')
    return
  }

  const canvas = document.createElement('canvas')
  const ctx = canvas.getContext('2d')
  
  // Canvas boyutları
  const padding = 40
  const lineHeight = 50
  const headerHeight = 120
  const footerHeight = 80
  const width = 650
  const height = headerHeight + (products.value.length * lineHeight) + footerHeight + padding * 2
  
  canvas.width = width
  canvas.height = height
  
  // Beyaz arka plan
  ctx.fillStyle = '#ffffff'
  ctx.fillRect(0, 0, width, height)
  
  // Header arka plan
  const gradient = ctx.createLinearGradient(0, 0, width, 0)
  gradient.addColorStop(0, '#0ea5e9')
  gradient.addColorStop(1, '#8b5cf6')
  ctx.fillStyle = gradient
  ctx.fillRect(0, 0, width, headerHeight)
  
  // Header yazıları
  ctx.fillStyle = '#ffffff'
  ctx.font = 'bold 24px Segoe UI, sans-serif'
  ctx.textAlign = 'center'
  ctx.fillText('FİYAT TEKLİFİ', width / 2, 45)
  
  ctx.font = '16px Segoe UI, sans-serif'
  ctx.fillText(orderTitle.value || 'İsimsiz Sipariş', width / 2, 75)
  ctx.fillText(new Date().toLocaleDateString('tr-TR'), width / 2, 100)
  
  // Ürün listesi
  let y = headerHeight + padding
  ctx.textAlign = 'left'
  
  products.value.forEach((p, i) => {
    // Satır arka plan (alternatif)
    if (i % 2 === 0) {
      ctx.fillStyle = '#f8fafc'
      ctx.fillRect(0, y - 30, width, lineHeight)
    }
    
    // Ürün adı
    ctx.fillStyle = '#0f172a'
    ctx.font = '16px Segoe UI, sans-serif'
    const productName = `${i + 1}. ${p.product_name}`
    ctx.fillText(productName, padding, y)
    
    // Adet bilgisi (ürün adının yanında, farklı renk ve kalınlık)
    const nameWidth = ctx.measureText(productName).width
    ctx.fillStyle = '#8b5cf6' // Mor renk
    ctx.font = 'bold 14px Segoe UI, sans-serif'
    ctx.fillText(`  ${p.quantity} Adet`, padding + nameWidth, y)
    
    // Fiyat (sağda)
    ctx.textAlign = 'right'
    ctx.font = 'bold 18px Segoe UI, sans-serif'
    ctx.fillStyle = '#10b981'
    ctx.fillText(`${currency.value}${formatPrice(p.total_price)}`, width - padding, y)
    ctx.textAlign = 'left'
    
    y += lineHeight
  })
  
  // Toplam bölümü
  y += 10
  ctx.fillStyle = '#0f172a'
  ctx.fillRect(0, y - 10, width, 2)
  
  y += 30
  ctx.font = 'bold 22px Segoe UI, sans-serif'
  ctx.fillStyle = '#0f172a'
  ctx.fillText('TOPLAM:', padding, y)
  
  ctx.textAlign = 'right'
  ctx.fillStyle = '#0ea5e9'
  ctx.fillText(`${currency.value}${formatPrice(grandTotal.value)}`, width - padding, y)
  
  // PNG olarak indir
  canvas.toBlob((blob) => {
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `fiyat_teklifi_${orderTitle.value || 'isimsiz'}_${Date.now()}.png`
    a.click()
    URL.revokeObjectURL(url)
    showToast(t('export.pngDownloaded'))
  }, 'image/png')
}

function exportPngWhatsApp() {
  if (products.value.length === 0) {
    showToast(t('app.toasts.listEmpty'), 'error')
    return
  }

  if (!customerPhone.value) {
    showToast(t('app.toasts.customerPhoneNotFound'), 'error')
    return
  }

  const canvas = document.createElement('canvas')
  const ctx = canvas.getContext('2d')
  
  // Canvas boyutları
  const padding = 40
  const lineHeight = 50
  const headerHeight = 120
  const footerHeight = 80
  const width = 650
  const height = headerHeight + (products.value.length * lineHeight) + footerHeight + padding * 2
  
  canvas.width = width
  canvas.height = height
  
  // Beyaz arka plan
  ctx.fillStyle = '#ffffff'
  ctx.fillRect(0, 0, width, height)
  
  // Header arka plan
  const gradient = ctx.createLinearGradient(0, 0, width, 0)
  gradient.addColorStop(0, '#0ea5e9')
  gradient.addColorStop(1, '#8b5cf6')
  ctx.fillStyle = gradient
  ctx.fillRect(0, 0, width, headerHeight)
  
  // Header yazıları
  ctx.fillStyle = '#ffffff'
  ctx.font = 'bold 24px Segoe UI, sans-serif'
  ctx.textAlign = 'center'
  ctx.fillText('FİYAT TEKLİFİ', width / 2, 45)
  
  ctx.font = '16px Segoe UI, sans-serif'
  ctx.fillText(orderTitle.value || 'İsimsiz Sipariş', width / 2, 75)
  ctx.fillText(new Date().toLocaleDateString('tr-TR'), width / 2, 100)
  
  // Ürün listesi
  let y = headerHeight + padding
  ctx.textAlign = 'left'
  
  products.value.forEach((p, i) => {
    // Satır arka plan (alternatif)
    if (i % 2 === 0) {
      ctx.fillStyle = '#f8fafc'
      ctx.fillRect(0, y - 30, width, lineHeight)
    }
    
    // Ürün adı
    ctx.fillStyle = '#0f172a'
    ctx.font = '16px Segoe UI, sans-serif'
    const productName = `${i + 1}. ${p.product_name}`
    ctx.fillText(productName, padding, y)
    
    // Adet bilgisi
    const nameWidth = ctx.measureText(productName).width
    ctx.fillStyle = '#8b5cf6'
    ctx.font = 'bold 14px Segoe UI, sans-serif'
    ctx.fillText(`  ${p.quantity} Adet`, padding + nameWidth, y)
    
    // Fiyat
    ctx.textAlign = 'right'
    ctx.font = 'bold 18px Segoe UI, sans-serif'
    ctx.fillStyle = '#10b981'
    ctx.fillText(`${currency.value}${formatPrice(p.total_price)}`, width - padding, y)
    ctx.textAlign = 'left'
    
    y += lineHeight
  })
  
  // Toplam bölümü
  y += 10
  ctx.fillStyle = '#0f172a'
  ctx.fillRect(0, y - 10, width, 2)
  
  y += 30
  ctx.font = 'bold 22px Segoe UI, sans-serif'
  ctx.fillStyle = '#0f172a'
  ctx.fillText('TOPLAM:', padding, y)
  
  ctx.textAlign = 'right'
  ctx.fillStyle = '#0ea5e9'
  ctx.fillText(`${currency.value}${formatPrice(grandTotal.value)}`, width - padding, y)
  
  // PNG'yi clipboard'a kopyala ve WhatsApp'ı aç
  canvas.toBlob(async (blob) => {
    const fileName = `fiyat_teklifi_${orderTitle.value || 'isimsiz'}_${Date.now()}.png`
    
    try {
      // Resmi clipboard'a kopyala
      await navigator.clipboard.write([
        new ClipboardItem({
          'image/png': blob
        })
      ])
      
      // WhatsApp'ı aç
      const phone = customerPhone.value.replace(/\D/g, '')
      window.open(`https://wa.me/${phone}`, '_blank')
      
      showToast(t('export.imageCopied'), 'success')
    } catch (err) {
      // Clipboard API desteklenmiyorsa dosyayı indir
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = fileName
      a.click()
      URL.revokeObjectURL(url)
      
      const phone = customerPhone.value.replace(/\D/g, '')
      setTimeout(() => {
        window.open(`https://wa.me/${phone}`, '_blank')
        showToast(t('export.pngDownloadedManual'), 'success')
      }, 500)
    }
  }, 'image/png')
}

function formatPrice(n) {
  const loc = locale.value || navigator.language || 'en-US'
  return (n || 0).toLocaleString(loc, { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

// Stock Handlers
function handleStockIn(product) {
  selectedProduct.value = product
  showStockEntry.value = true
}

function handleCriticalStockIn(product) {
  selectedProduct.value = product
  showStockEntry.value = true
}

function handleShowAllCriticalStock() {
  // Switch to stock tab and filter by critical
  activeTab.value = 'stock'
  showOnlyCriticalStock.value = true
  
  // Reset after a short delay so it can be triggered again
  setTimeout(() => {
    showOnlyCriticalStock.value = false
  }, 500)
}

function handleStockOut(product) {
  selectedProduct.value = product
  showStockExit.value = true
}

function handleMovements(product) {
  selectedProduct.value = product
  showStockMovements.value = true
}

function handleNewProduct() {
  editingProduct.value = null
  showProductForm.value = true
}

function handleEditProduct(product) {
  editingProduct.value = product
  showProductForm.value = true
}

function handleDeleteStockProduct(product) {
  modals.value.showConfirm(
    t('app.confirms.deleteStockProduct'),
    t('app.confirms.deleteStockProductMessage', { name: product.name }),
    async () => {
      const result = await deleteStockProduct(product.id)
      if (result.error) {
        showToast(result.error, 'error')
      } else {
        showToast(t('app.toasts.productDeleted'))
      }
    },
    '🗑️'
  )
}

function handleStockOperationSuccess(message) {
  showToast(message || 'Operation successful')
  loadStockProducts()
}

// Bulk stock operations
const bulkProducts = ref([])

function openBulkStockEntry() {
  selectedProduct.value = null
  showStockEntry.value = true
}

function openBulkStockExit() {
  selectedProduct.value = null
  showStockExit.value = true
}

function handleBulkStockIn(products) {
  bulkProducts.value = products
  selectedProduct.value = null
  showStockEntry.value = true
}

function handleBulkStockOut(products) {
  bulkProducts.value = products
  selectedProduct.value = null
  showStockExit.value = true
}

function handleBulkEdit(products) {
  bulkEditProducts.value = products
  showBulkEdit.value = true
}

function closeBulkEditModal() {
  showBulkEdit.value = false
  bulkEditProducts.value = []
}

function handleBulkEditSuccess(message) {
  showToast(message || 'Ürünler güncellendi')
  loadStockProducts()
}

function closeStockModals() {
  showStockEntry.value = false
  showStockExit.value = false
  bulkProducts.value = []
  selectedProduct.value = null
}
function handleProductFormSuccess(message) {
  showToast(message || 'Operation successful')
  loadStockProducts()
}

// Check for updates on startup (if enabled)
async function checkForUpdatesOnStartup() {
  try {
    // Wait for settings to load
    const settings = JSON.parse(localStorage.getItem('appSettings') || '{}')
    if (settings.autoUpdateCheck === false) {
      console.log('[App] Auto update check disabled')
      return
    }
    
    // Check if API is available
    if (typeof checkForUpdates === 'undefined') {
      console.log('[App] Update API not available')
      return
    }
    
    console.log('[App] Checking for updates...')
    const result = await checkForUpdates()
    const updateInfo = JSON.parse(result)
    
    if (updateInfo.available) {
      console.log('[App] Update available:', updateInfo.latest_version)
      // Show toast notification
      const { showToast } = await import('./composables/useToast.js')
      showToast(
        t('app.updateAvailable', { version: updateInfo.latest_version }), 
        'info',
        8000
      )
    } else if (updateInfo.error) {
      console.warn('[App] Update check error:', updateInfo.error)
    } else {
      console.log('[App] Running latest version:', updateInfo.current_version)
    }
  } catch (e) {
    console.warn('[App] Update check failed:', e)
  }
}

// Init with loadout overlay
onMounted(async () => {
  initialLoading.value = true
  try {
    // Load developer mode setting
    if (typeof getDeveloperMode !== 'undefined') {
      try {
        const result = await getDeveloperMode()
        const parsed = JSON.parse(result)
        developerMode.value = parsed.enabled || false
      } catch (e) {
        console.warn('[App] Failed to load developer mode:', e)
      }
    }
    
    await Promise.all([
      loadData(),
      loadStockProducts()
    ])
    
    // Check for updates after data loads (non-blocking)
    checkForUpdatesOnStartup()
  } catch (e) {
    // Log but continue to hide overlay so UI is usable
    console.error('Initial load error:', e)
  } finally {
    // ensure overlay hides even if one of the loads fails; keep it visible briefly for UX
    setTimeout(() => {
      initialLoading.value = false
    }, 300)
  }
})
</script>

<style scoped>
.tab-nav {
  position: sticky;
  top: 0;
  z-index: 40;
}

.tab-list {
  display: flex;
  gap: 0.25rem;
  padding: 0.5rem 0;
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.625rem 1.25rem;
  border: none;
  background: transparent;
  color: var(--text-muted);
  font-size: 0.9rem;
  font-weight: 500;
  border-radius: 0.375rem;
  cursor: pointer;
  transition: all 0.15s;
}

.tab-btn:hover {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.tab-btn.active {
  background: var(--accent-color);
  color: white;
}

.tab-btn.active:hover {
  background: var(--accent-hover);
}

.tab-btn svg {
  opacity: 0.8;
}

.tab-btn.active svg {
  opacity: 1;
}

/* Developer tab special styling */
.tab-btn.developer-tab {
  color: #e94560;
}

.tab-btn.developer-tab:hover {
  background: rgba(233, 69, 96, 0.1);
}

.tab-btn.developer-tab.active {
  background: linear-gradient(135deg, #e94560, #c23a51);
}
</style>