<template>
  <div class="min-h-screen flex flex-col transition-colors duration-300" style="background: var(--bg-primary); color: var(--text-primary);">
    <AppHeader 
      @showHistory="modals.showHistory()" 
      @newOrder="handleNewOrder"
      @showCatalog="showCatalog = true"
    />
    
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
    
    <!-- Footer -->
    <footer class="border-t py-3 px-6 text-center text-sm" style="background: var(--bg-secondary); border-color: var(--border-color); color: var(--text-muted);">
      <div class="flex items-center justify-center gap-2 flex-wrap">
        <span class="opacity-70">🔓 Open Source</span>
        <span class="opacity-50">•</span>
        <span>Sürüm:</span>
        <span class="font-semibold text-success">25.12-stable</span>
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
        © 2025 Durasoft • MIT License<br>
        <b>Son Özellikler:</b>
        <ul class="list-disc list-inside text-left mx-auto max-w-xl mt-1">
          <li>Gelişmiş arama + tarih filtresi birlikte çalışır</li>
          <li>Sipariş kalemlerinde içerik arama, gizlenen ürün sayısı</li>
          <li>Çoklu seçim ve silme</li>
          <li>WhatsApp'a resim kopyalama</li>
          <li>Tema uyumlu toast ve dropdown</li>
          <li>Karanlık/aydınlık tema desteği</li>
        </ul>
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
    
    <Toast />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useOrder } from '@/composables/useOrder'
import { useToast } from '@/composables/useToast'
import AppHeader from '@/components/AppHeader.vue'
import OrderBar from '@/components/OrderBar.vue'
import ProductForm from '@/components/ProductForm.vue'
import OrderItems from '@/components/OrderItems.vue'
import OrdersList from '@/components/OrdersList.vue'
import OrderSummary from '@/components/OrderSummary.vue'
import Modals from '@/components/Modals.vue'
import CatalogModal from '@/components/CatalogModal.vue'
import Toast from '@/components/Toast.vue'

const showCatalog = ref(false)

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

// Init
onMounted(() => {
  loadData()
})
</script>
