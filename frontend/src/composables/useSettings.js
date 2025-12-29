import { ref, watch } from 'vue'

// Shared settings state
const settings = ref({
  // Display settings
  itemsPerPage: 25,
  theme: 'dark',
  // Show loadout overlay on every app start (optional)
  showLoadoutAlways: false,
  
  // Stock settings
  autoDeductStock: false,
  defaultUnit: 'adet',
  
  // Developer settings
  developerMode: false,
  
  // Table columns visibility
  stockListColumns: {
    name: true,
    oem_number: true,
    brand: true,
    category: true,
    unit: true,
    stock_quantity: true,
    critical_stock: true,
    status: true
  }
})

const isLoaded = ref(false)

// Load settings from localStorage
function loadSettings() {
  if (isLoaded.value) return
  
  try {
    const saved = localStorage.getItem('appSettings')
    if (saved) {
      const parsed = JSON.parse(saved)
      settings.value = { ...settings.value, ...parsed }
    }
    
    // Also load theme from old key for backward compatibility
    const savedTheme = localStorage.getItem('theme')
    if (savedTheme) {
      settings.value.theme = savedTheme
    }
    
    // Apply theme
    document.documentElement.classList.toggle('dark', settings.value.theme === 'dark')
    
    isLoaded.value = true
  } catch (e) {
    console.error('Settings load error:', e)
  }
}

// Save settings to localStorage
function saveSettings() {
  try {
    localStorage.setItem('appSettings', JSON.stringify(settings.value))
    // Also save theme to old key for backward compatibility
    localStorage.setItem('theme', settings.value.theme)
  } catch (e) {
    console.error('Settings save error:', e)
  }
}

// Toggle theme
function toggleTheme() {
  settings.value.theme = settings.value.theme === 'dark' ? 'light' : 'dark'
  document.documentElement.classList.toggle('dark', settings.value.theme === 'dark')
  saveSettings()
}

// Set theme directly
function setTheme(theme) {
  settings.value.theme = theme
  document.documentElement.classList.toggle('dark', theme === 'dark')
  saveSettings()
}

// Update a setting
function updateSetting(key, value) {
  if (key.includes('.')) {
    // Handle nested keys like 'stockListColumns.name'
    const keys = key.split('.')
    let obj = settings.value
    for (let i = 0; i < keys.length - 1; i++) {
      obj = obj[keys[i]]
    }
    obj[keys[keys.length - 1]] = value
  } else {
    settings.value[key] = value
  }
  saveSettings()
}

// Watch for changes and auto-save
watch(settings, () => {
  if (isLoaded.value) {
    saveSettings()
  }
}, { deep: true })

// Changelog data for About section
// Status types (can be array for multiple badges):
// - release: Son Sürüm (purple)
// - stable: Kararlı (green)
// - unstable: Geliştirme (red)
// - latest: En Son (blue)
// - compatible: Uyumlu (teal)
// - incompatible: Uyumsuz (orange)
// - first: İlk Sürüm (gray)
// - pre-release: Ön Sürüm (amber)
const changelog = [
  {
    version: '25.12.2',
    status: ['latest', 'stable'],
    date: '2025-12-29',
    changes: [
      'Yeni kayıt sırasında input artık boşalmıyor ve v-model hemen güncelleniyor',
      'Yeni oluşturulan kategori/marka uygulama listelerine ekleniyor',
      'Eski tabloları yeni formata migrate edecek tek seferlik göç fonksiyonu eklendi',
      'Loadout screen eklendi'
    ]
  },
  {
    version: '25.12.1',
    status: ['pre-release', 'unstable'],
    date: '2025-12-29',
    changes: [
      'Stok ve Raporlama sistemi',
      'Tüm select elementleri için AutocompleteSelect bileşeni',
      'Toplu ürün düzenleme',
      'Stok giriş/çıkış modallarında autocomplete ürün seçimi',
      'Kategori, marka ve birim için akıllı dropdown',
      'Stok listesinde sütun sıralama ve gizleme',
      'CSV export özelliği',
      'Server-side pagination ile performans optimizasyonu',
      'Ayarlar modalı ve tema yönetimi',
      'Uygulama adı AutoManagement olarak güncellendi'
    ]
  },
  {
    version: '25.12',
    status: ['first', 'incompatible'],
    date: '2025-12-24',
    changes: [
      'Gelişmiş autocomplete bileşeni',
      'Sipariş kalemlerinde içerik arama, gizlenen ürün sayısı',
      'Sipariş Yönetiminde Çoklu seçim ve silme',
      'WhatsApp\'a resim kopyalama',
    ]
  },
  {
    version: '25.11',
    status: ['pre-release', 'incompatible'],
    date: '2025-11-15',
    changes: [
      'Gelişmiş arama + tarih filtresi birlikte çalışır',
      'Tema uyumlu toast ve dropdown',
      'Sipariş Yönetim sistemi',
      'Karanlık/aydınlık tema desteği'
    ]
  }
]

export function useSettings() {
  // Auto-load on first use
  loadSettings()
  
  return {
    settings,
    changelog,
    loadSettings,
    saveSettings,
    toggleTheme,
    setTheme,
    updateSetting
  }
}
