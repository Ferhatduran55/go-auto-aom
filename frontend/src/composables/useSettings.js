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
  hideCriticalStockWarning: false, // Kritik stok uyarısını gizle
  currency: '₺',

  // Developer settings
  developerMode: false,
  logBufferSize: 1000, // Geliştirici konsolunda tutulacak log sayısı

  // Update settings
  autoUpdateCheck: true, // Uygulama açılışında otomatik güncelleme kontrolü

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

// Changelog is provided by i18n locales as `settings.changelogEntries` (so it becomes locale-aware)
import { computed } from 'vue'
import { useI18n } from '../i18n'

// changelog will be a computed value returning an array from the current locale
function getLocalizedChangelog() {
  const { t } = useI18n()
  return computed(() => {
    const val = t('settings.changelogEntries')
    return Array.isArray(val) ? val : []
  })
}

const changelog = getLocalizedChangelog()

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
