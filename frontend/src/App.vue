<template>
  <div class="min-h-screen flex flex-col transition-colors duration-300" style="background: var(--bg-primary); color: var(--text-primary);">

    <!-- Loadout overlay -->
    <div v-if="initialLoading" class="fixed inset-0 z-[2000] flex items-center justify-center bg-black/90">
      <div class="bg-card p-6 rounded-lg text-center shadow-lg flex flex-col items-center gap-4">
        <div class="w-12 h-12 border-4 border-accent rounded-full border-t-transparent animate-spin"></div>
        <div class="text-2xl font-bold">Verileriniz kurtarılıyor</div>
        <div class="text-sm text-muted">Lütfen bekleyin...</div>
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
      <template #kritik-stok>
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
            Siparişler
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
            Stok Yönetimi
          </button>
          <button 
            @click="activeTab = 'reports'" 
            :class="['tab-btn', { active: activeTab === 'reports' }]"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M3 3v18h18"></path>
              <path d="m19 9-5 5-4-4-3 3"></path>
            </svg>
            Raporlar
          </button>
        </div>
      </div>
    </nav>
    
    <!-- Siparişler View -->
    <template v-if="activeTab === 'orders'">
      <OrderBar />
      
      <main class="flex-1 grid grid-cols-[1fr_2fr_1fr] gap-6 p-6 max-w-[1800px] mx-auto w-full">
        <!-- Left: Product Form -->
        <ProductForm @clearAll="handleClearAll" />
        
        <!-- Center: Order Items -->
        <OrderItems 
          @newOrder="handleNewOrder"
          @deleteProduct="handleDeleteProduct"
        />
        
        <!-- Right: Summary + Orders List -->
        <div class="flex flex-col gap-4">
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
    
    <!-- Footer -->
    <footer class="border-t py-3 px-6 text-center text-sm" style="background: var(--bg-secondary); border-color: var(--border-color); color: var(--text-muted);">
      <div class="flex items-center justify-center gap-2 flex-wrap">
        <span class="opacity-70">🔓 Open Source</span>
        <span class="opacity-50">•</span>
        <span>Sürüm:</span>
        <span class="font-semibold text-success">25.12.3</span>
        <span class="opacity-50">•</span>
        <span>Geliştirici:</span>
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
        © 2025 Durasoft • MIT License
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
import { ref, onMounted } from 'vue'
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
import BulkEditModal from '@/components/BulkEditModal.vue'
import CriticalStockBadge from '@/components/CriticalStockBadge.vue'
import ReportsView from '@/components/ReportsView.vue'

const showCatalog = ref(false)
const showSettings = ref(false)
const activeTab = ref('orders') // orders | stock | reports

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
  isEditing
} = useOrder()

const { settings } = useSettings()
const initialLoading = ref(true)
const showProductForm = ref(false)
const showBulkEdit = ref(false)
const showOnlyCriticalStock = ref(false)
const selectedProduct = ref(null)
const editingProduct = ref(null)
const bulkEditProducts = ref([])
const stockListRef = ref(null)

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
    showToast('Liste boş', 'error')
    return
  }
  
  const result = await saveOrder()
  if (result.error) {
    showToast(result.error, 'error')
  } else {
    showToast('Sipariş kaydedildi')
    refreshTrigger.value++
  }
}

async function handleSaveAsNew() {
  currentOrderId.value = null
  await handleSaveOrder()
}

async function handleLoadOrder(id) {
  const result = await loadOrder(id)
  if (result.error) {
    showToast('Sipariş yüklenemedi', 'error')
  } else {
    showToast('Sipariş yüklendi')
  }
}

function handleDeleteProduct(id) {
  modals.value.showConfirm(
    'Ürün Silme',
    'Bu ürünü listeden silmek istediğinize emin misiniz?',
    () => {
      deleteProduct(id)
      showToast('Ürün silindi')
    },
    '🗑️'
  )
}

function handleDeleteOrder(id) {
  modals.value.showConfirm(
    'Sipariş Silme',
    'Bu siparişi kalıcı olarak silmek istediğinize emin misiniz?',
    async () => {
      const { api } = await import('@/api')
      const result = await api.deleteOrder(id)
      if (!result.error) {
        showToast('Sipariş silindi')
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
  // Switch to orders tab first if not already there
  activeTab.value = 'orders'
  
  // Reset order form
  resetOrder()
  showToast('Yeni sipariş başlatıldı')
}

function handleClearAll() {
  if (products.value.length === 0) return
  modals.value.showConfirm(
    'Listeyi Temizle',
    'Tüm ürünleri listeden silmek istediğinize emin misiniz?',
    () => {
      clearProducts()
      showToast('Liste temizlendi')
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
    showToast('Liste boş', 'error')
    return
  }

  // WhatsApp uyumlu format
  let content = `📋 *SİPARİŞ: ${orderTitle.value || 'İsimsiz'}*\n`
  content += `👤 *Müşteri:* ${customerName.value || '-'}\n`
  content += `📅 *Tarih:* ${new Date().toLocaleDateString('tr-TR')}\n`
  content += `\n${'─'.repeat(35)}\n\n`

  products.value.forEach((p, i) => {
    const durum = p.part_status === 'original' ? '🟢 Orijinal' : (p.part_status === 'zero' ? '🔵 Sıfır' : '🟡 Çıkma')
    content += `*${i + 1}. ${p.product_name}*\n`
    content += `   🔢 OEM: ${p.oem_number}\n`
    content += `   📦 Adet: ${p.quantity}\n`
    content += `   ${durum}\n`
    content += `   💰 *₺${formatPrice(p.total_price)}*\n\n`
  })

  content += `${'─'.repeat(35)}\n`
  content += `\n💵 *GENEL TOPLAM: ₺${formatPrice(grandTotal.value)}*\n`

  const blob = new Blob([content], { type: 'text/plain;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `siparis_${orderTitle.value || 'isimsiz'}_${Date.now()}.txt`
  a.click()
  URL.revokeObjectURL(url)
  showToast('TXT dosyası indirildi')
}

function exportTxtWhatsApp() {
  if (products.value.length === 0) {
    showToast('Liste boş', 'error')
    return
  }

  // WhatsApp uyumlu format - Türkçe karakterli
  let content = `*SİPARİŞ: ${orderTitle.value || 'İsimsiz'}*\n`
  content += `*Müşteri:* ${customerName.value || '-'}\n`
  content += `*Tarih:* ${new Date().toLocaleDateString('tr-TR')}\n`
  content += `\n${'─'.repeat(30)}\n\n`

  products.value.forEach((p, i) => {
    const durum = p.part_status === 'original' ? '[Orijinal]' : (p.part_status === 'zero' ? '[Sıfır]' : '[Çıkma]')
    content += `*${i + 1}. ${p.product_name}*\n`
    content += `   OEM: ${p.oem_number}\n`
    content += `   Adet: ${p.quantity}\n`
    content += `   ${durum}\n`
    content += `   *${formatPrice(p.total_price)} TL*\n\n`
  })

  content += `${'─'.repeat(30)}\n`
  content += `\n*GENEL TOPLAM: ${formatPrice(grandTotal.value)} TL*\n`

  // Telefon numarasını formatla
  if (!customerPhone.value) {
    showToast('Müşteri telefonu bulunamadı', 'error')
    return
  }
  
  const phone = customerPhone.value.replace(/\D/g, '')
  const encodedText = encodeURIComponent(content)
  
  // WhatsApp Web linkini aç - mesajı taslak olarak hazırlar
  window.open(`https://wa.me/${phone}?text=${encodedText}`, '_blank')
  showToast('WhatsApp açılıyor...')
}

function exportPng() {
  if (products.value.length === 0) {
    showToast('Liste boş', 'error')
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
    ctx.fillText(`₺${formatPrice(p.total_price)}`, width - padding, y)
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
  ctx.fillText(`₺${formatPrice(grandTotal.value)}`, width - padding, y)
  
  // PNG olarak indir
  canvas.toBlob((blob) => {
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `fiyat_teklifi_${orderTitle.value || 'isimsiz'}_${Date.now()}.png`
    a.click()
    URL.revokeObjectURL(url)
    showToast('PNG dosyası indirildi')
  }, 'image/png')
}

function exportPngWhatsApp() {
  if (products.value.length === 0) {
    showToast('Liste boş', 'error')
    return
  }

  if (!customerPhone.value) {
    showToast('Müşteri telefonu bulunamadı', 'error')
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
    ctx.fillText(`₺${formatPrice(p.total_price)}`, width - padding, y)
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
  ctx.fillText(`₺${formatPrice(grandTotal.value)}`, width - padding, y)
  
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
      
      showToast('Resim panoya kopyalandı! WhatsApp\'ta Ctrl+V ile yapıştırın.', 'success')
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
        showToast('PNG indirildi. WhatsApp\'ta manuel paylaşabilirsiniz.', 'success')
      }, 500)
    }
  }, 'image/png')
}

function formatPrice(n) {
  return (n || 0).toLocaleString('tr-TR', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
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
    'Ürün Silme',
    `"${product.name}" ürününü silmek istediğinize emin misiniz?`,
    async () => {
      const result = await deleteStockProduct(product.id)
      if (result.error) {
        showToast(result.error, 'error')
      } else {
        showToast('Ürün silindi')
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

// Init with loadout overlay
onMounted(async () => {
  initialLoading.value = true
  try {
    await Promise.all([
      loadData(),
      loadStockProducts()
    ])
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
</style>