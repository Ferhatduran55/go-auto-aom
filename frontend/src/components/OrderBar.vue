<template>
  <div 
    class="border-b px-6 grid grid-cols-[1fr_1fr_1fr_1fr_auto] gap-4 items-end transition-all duration-300"
    :class="[
      isActive ? 'sticky top-[80px] z-30 backdrop-blur-md py-3' : 'py-4',
      isActive && isEditing ? 'bg-accent/10' : (isActive ? 'bg-success/10' : '')
    ]"
    :style="{
      background: isActive ? undefined : 'var(--bg-secondary)',
      borderColor: 'var(--border-color)'
    }"
  >
    <div class="flex flex-col gap-1.5">
      <label class="text-xs font-bold uppercase tracking-wide" style="color: var(--text-muted);">📋 {{ t('orders.orderTitle') }}</label>
      <input 
        type="text" 
        v-model="orderTitle" 
        class="form-input" 
        :placeholder="t('orders.orderTitlePlaceholder')"
      >
    </div>
    
    <div class="flex flex-col gap-1.5 relative">
      <label class="text-xs font-bold uppercase tracking-wide" style="color: var(--text-muted);">👤 {{ t('orders.customer') }}</label>
      <div class="flex gap-2">
        <input 
          type="text" 
          v-model="customerName" 
          @input="onCustomerInput"
          @focus="showCustomerDropdown = true"
          class="form-input flex-1" 
          :placeholder="t('orders.customerNamePlaceholder')"
          autocomplete="off"
        >
        <!-- Müşteri Filtre Toggle Butonu -->
        <button 
          v-if="currentCustomerId"
          @click="toggleCustomerFilter"
          :class="['px-3 rounded-xl border-2 transition-colors', customerFilterActive ? 'bg-accent border-accent text-white' : 'hover:bg-accent hover:text-white hover:border-accent']"
          :style="!customerFilterActive ? { borderColor: 'var(--border-color)', color: 'var(--text-muted)' } : {}"
          :title="customerFilterActive ? t('orders.clearCustomerFilter') : t('orders.filterByCustomer')"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z"/>
          </svg>
        </button>
      </div>
      <div v-if="showCustomerDropdown && (filteredCustomers.length > 0 || (customerName && customerName.length >= 2))" 
           class="absolute top-full left-0 right-0 mt-1 rounded-xl max-h-64 overflow-y-auto z-50 shadow-xl" style="background: var(--bg-card); border: 2px solid var(--border-color);">
        <div 
          v-for="c in filteredCustomers" 
          :key="c.id"
          @click="selectCustomer(c)"
          class="px-4 py-3 cursor-pointer border-b last:border-b-0 hover:bg-accent hover:text-white transition-colors" style="border-color: var(--border-color);"
        >
          <div class="flex items-center justify-between">
            <span class="font-semibold" style="color: var(--text-primary);">👤 {{ c.name }}</span>
            <span class="text-xs px-2 py-0.5 rounded bg-accent/20 text-accent font-mono">#{{ c.id?.substring(0, 6) }}</span>
          </div>
          <div class="text-sm" style="color: var(--text-muted);">
            <span v-if="c.phone">📞 {{ c.phone }} • </span>
            {{ c.order_count || 0 }} {{ t('common.order') }} • {{ currency }}{{ formatPrice(c.total_amount || 0) }}
          </div>
        </div>
        <!-- Yeni müşteri oluştur seçeneği -->
        <div 
          v-if="customerName && customerName.length >= 2"
          @click="createNewCustomer"
          class="px-4 py-3 cursor-pointer hover:bg-success hover:text-white transition-colors font-semibold" 
          :class="{ 'border-t-2': filteredCustomers.length > 0 }"
          style="border-color: var(--border-color); color: var(--text-primary);"
        >
          ➕ {{ t('orders.createCustomer', { name: customerName }) }}
        </div>
      </div>
    </div>
    
    <div class="flex flex-col gap-1.5 relative">
      <label class="text-xs font-bold uppercase tracking-wide" style="color: var(--text-muted);">📞 {{ t('orders.customerPhone') }}</label>
      <div class="flex items-stretch gap-2">
        <div class="flex items-stretch rounded-xl overflow-hidden border-2 flex-1" style="border-color: var(--border-color);">
          <input 
            type="text" 
            v-model="countryCode" 
            class="w-14 text-center text-base font-medium border-0 border-r-2 focus:ring-0 focus:outline-none px-2 py-3.5" 
            style="background: var(--bg-secondary); border-color: var(--border-color); color: var(--text-primary);"
            placeholder="+90"
          >
          <input 
            type="tel" 
            v-model="phoneInput" 
            @input="onPhoneInput"
            @focus="showPhoneDropdown = true"
            class="flex-1 min-w-0 px-4 py-3.5 text-base border-0 focus:ring-0 focus:outline-none" 
            style="background: var(--bg-secondary); color: var(--text-primary);"
            placeholder="5XX XXX XX XX"
            autocomplete="off"
          >
        </div>
        <button 
          v-if="phoneInput"
          @click="openWhatsApp"
          class="px-3 rounded-xl border-2 hover:bg-green-500 hover:text-white hover:border-green-500 transition-colors"
          style="border-color: var(--border-color); color: var(--text-muted);"
          :title="t('common.openInWhatsApp')"
        >
          <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
            <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413Z"/>
          </svg>
        </button>
        <button 
          v-if="phoneInput"
          @click="copyPhone"
          class="px-3 rounded-xl border-2 hover:bg-accent hover:text-white hover:border-accent transition-colors"
          style="border-color: var(--border-color); color: var(--text-muted);"
          :title="t('common.copyNumber')"
        >
          📋
        </button>
      </div>
      <!-- Telefon dropdown -->
      <div v-if="showPhoneDropdown && filteredByPhone.length > 0" 
           class="absolute top-full left-0 right-0 mt-1 rounded-xl max-h-64 overflow-y-auto z-50 shadow-xl" style="background: var(--bg-card); border: 2px solid var(--border-color);">
        <div 
          v-for="c in filteredByPhone" 
          :key="c.id"
          @click="selectCustomer(c)"
          class="px-4 py-3 cursor-pointer border-b last:border-b-0 hover:bg-accent hover:text-white transition-colors" style="border-color: var(--border-color);"
        >
          <div class="flex items-center justify-between">
            <span class="font-semibold" style="color: var(--text-primary);">👤 {{ c.name }}</span>
            <span class="text-xs px-2 py-0.5 rounded bg-accent/20 text-accent font-mono">#{{ c.id?.substring(0, 6) }}</span>
          </div>
          <div class="text-sm" style="color: var(--text-muted);">
            📞 {{ c.phone }} • {{ c.order_count || 0 }} {{ t('common.order') }}
          </div>
        </div>
      </div>
    </div>
    
    <div class="flex flex-col gap-1.5">
      <label class="text-xs font-bold uppercase tracking-wide" style="color: var(--text-muted);">📅 {{ t('orders.orderDate') }}</label>
      <input 
        type="text" 
        :value="currentDate" 
        class="form-input cursor-not-allowed opacity-70" 
        readonly
      >
    </div>
    
    <div>
      <span 
        v-if="isEditing" 
        class="bg-accent text-white px-5 py-3 rounded-lg font-bold text-sm inline-block"
      >
        ✏️ {{ t('orders.editing') }}
      </span>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useOrder } from '@/composables/useOrder'
import { useToast } from '@/composables/useToast'
import { useI18n } from '@/i18n'
import { useSettings } from '@/composables/useSettings'

const { t, locale } = useI18n()
const { settings } = useSettings()
const currency = computed(() => settings?.value?.currency || '₺')

const { orderTitle, customerName, customerPhone, currentCustomerId, allCustomers, isEditing, products, customerFilterActive } = useOrder()
const { showToast } = useToast()

// Aktif durum: ürün varsa veya düzenleme modundaysa
const isActive = computed(() => {
  return products.value.length > 0 || isEditing.value
})

// Müşteri filtresini toggle et
function toggleCustomerFilter() {
  customerFilterActive.value = !customerFilterActive.value
  if (customerFilterActive.value) {
    showToast(`📋 ${customerName.value} ${t('orders.customerFiltered')}`, 'success')
  } else {
    showToast(t('orders.filterRemoved'), 'info')
  }
}

const showCustomerDropdown = ref(false)
const showPhoneDropdown = ref(false)
const phoneInput = ref('')
const countryCode = ref('+90')

// customerPhone değiştiğinde phoneInput'u güncelle
watch(customerPhone, (newVal) => {
  if (newVal) {
    // Ülke kodunu ve numarayı ayır (örn: +90 5321234567 veya +905321234567)
    const match = newVal.match(/^(\+\d{1,4})\s*(.*)$/)
    if (match) {
      countryCode.value = match[1]
      // Numaradan tekrar ülke kodu varsa kaldır
      let phone = match[2]
      if (phone.startsWith('0')) {
        phone = phone.substring(1) // Başındaki 0'ı kaldır
      }
      phoneInput.value = phone
    } else if (newVal.startsWith('0')) {
      // 05XX formatı - ülke kodu yok
      phoneInput.value = newVal.substring(1)
    } else {
      phoneInput.value = newVal
    }
  } else {
    phoneInput.value = ''
  }
}, { immediate: true })

const currentDate = computed(() => {
  const loc = locale.value || navigator.language || 'en-US'
  return new Date().toLocaleDateString(loc, { 
    day: 'numeric', 
    month: 'long', 
    year: 'numeric' 
  })
})

const filteredCustomers = computed(() => {
  if (!customerName.value || customerName.value.length < 1) return []
  const search = customerName.value.toLowerCase().trim()
  return allCustomers.value
    .filter(c => c.name.toLowerCase().includes(search))
    .slice(0, 6)
})

// Telefon ile arama
const filteredByPhone = computed(() => {
  if (!phoneInput.value || phoneInput.value.length < 3) return []
  const search = phoneInput.value.replace(/\s/g, '')
  return allCustomers.value
    .filter(c => c.phone && c.phone.replace(/\D/g, '').includes(search))
    .slice(0, 6)
})

// Tam eşleşme var mı? (aynı isimde müşteri zaten seçili mi)
const exactMatch = computed(() => {
  if (!customerName.value) return false
  const search = customerName.value.toLowerCase().trim()
  return allCustomers.value.some(c => c.name.toLowerCase() === search)
})

function onCustomerInput() {
  showCustomerDropdown.value = true
  showPhoneDropdown.value = false
  // Yeni isim girildiğinde customer ID'yi temizle (yeni müşteri oluşturulacak)
  currentCustomerId.value = null
  // Filtreyi kapat
  customerFilterActive.value = false
}

function onPhoneInput() {
  showPhoneDropdown.value = true
  showCustomerDropdown.value = false
  
  // Telefon numarasını formatla
  let digits = phoneInput.value.replace(/\D/g, '')
  
  // customerPhone'u güncelle (ülke kodu ile)
  if (digits) {
    customerPhone.value = countryCode.value + ' ' + digits
  } else {
    customerPhone.value = ''
  }
  
  // Tam eşleşme varsa otomatik seç
  const exactPhoneMatch = allCustomers.value.find(c => {
    if (!c.phone) return false
    const cDigits = c.phone.replace(/\D/g, '')
    return cDigits === '90' + digits || cDigits === digits
  })
  
  if (exactPhoneMatch && digits.length >= 10) {
    selectCustomer(exactPhoneMatch)
  }
}

function selectCustomer(customer) {
  customerName.value = customer.name
  if (customer.phone) {
    customerPhone.value = customer.phone
    // Ülke kodunu ve numarayı ayır
    const match = customer.phone.match(/^(\+\d{1,4})\s*(.*)$/)
    if (match) {
      countryCode.value = match[1]
      phoneInput.value = match[2]
    } else {
      phoneInput.value = customer.phone
    }
  } else {
    customerPhone.value = ''
    phoneInput.value = ''
  }
  currentCustomerId.value = customer.id
  showCustomerDropdown.value = false
  showPhoneDropdown.value = false
}

function createNewCustomer() {
  // ID'yi null bırakarak yeni müşteri oluşturulmasını sağla
  currentCustomerId.value = null
  showCustomerDropdown.value = false
  showPhoneDropdown.value = false
}

function copyPhone() {
  const fullPhone = countryCode.value + ' ' + phoneInput.value
  navigator.clipboard.writeText(fullPhone).then(() => {
    showToast(t('app.toasts.phoneCopied', { phone: fullPhone }), 'success')
  }).catch(err => {
    showToast(t('app.toasts.copyError'), 'error')
  })
}

function openWhatsApp() {
  if (!phoneInput.value) return
  
  // Ülke kodunu temizle ve formatla
  let phone = phoneInput.value.replace(/\D/g, '')
  let code = countryCode.value.replace(/\D/g, '')
  
  // Tam numarayı oluştur
  const fullPhone = code + phone
  
  // WhatsApp linkini aç
  window.open(`https://wa.me/${fullPhone}`, '_blank')
}

function formatPrice(n) {
  const loc = locale.value || navigator.language || 'en-US'
  return (n || 0).toLocaleString(loc, { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

function handleClickOutside(e) {
  if (!e.target.closest('.relative')) {
    showCustomerDropdown.value = false
    showPhoneDropdown.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>
