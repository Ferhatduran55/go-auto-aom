// i18n - Çoklu dil desteği modülü
// Bu modül genişletilebilir bir yapıda tasarlanmıştır

import { ref, computed } from 'vue'
import tr from './locales/tr.json'
import en from './locales/en.json'

// Mevcut diller
const locales = {
  tr,
  en
}

// Aktif dil (localStorage'dan al veya varsayılan tr)
const currentLocale = ref(localStorage.getItem('app_locale') || 'tr')

// Dil değiştirme fonksiyonu
export function setLocale(locale) {
  if (locales[locale]) {
    currentLocale.value = locale
    localStorage.setItem('app_locale', locale)
    document.documentElement.setAttribute('lang', locale)
  }
}

// Mevcut dili al
export function getLocale() {
  return currentLocale.value
}

// Çeviri fonksiyonu - nested keys destekler (örn: "chatbot.title")
export function t(key, params = {}) {
  const keys = key.split('.')
  let value = locales[currentLocale.value]
  
  for (const k of keys) {
    if (value && typeof value === 'object' && k in value) {
      value = value[k]
    } else {
      // Fallback: İngilizce'ye bak
      value = locales['en']
      for (const k2 of keys) {
        if (value && typeof value === 'object' && k2 in value) {
          value = value[k2]
        } else {
          return key // Bulunamazsa key'i döndür
        }
      }
      break
    }
  }
  
  // Parametre yerleştirme: {param} -> value
  if (typeof value === 'string' && Object.keys(params).length > 0) {
    for (const [paramKey, paramValue] of Object.entries(params)) {
      value = value.replace(new RegExp(`\\{${paramKey}\\}`, 'g'), paramValue)
    }
  }
  
  return value || key
}

// Reactive translation hook for Vue components
export function useI18n() {
  const locale = computed(() => currentLocale.value)
  
  const availableLocales = computed(() => Object.keys(locales))
  
  return {
    t,
    locale,
    setLocale,
    getLocale,
    availableLocales
  }
}

// Dil bilgisi
export const localeNames = {
  tr: 'Türkçe',
  en: 'English'
}
