<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-container settings-modal">
      <div class="modal-header">
        <h2>
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"/>
            <circle cx="12" cy="12" r="3"/>
          </svg>
          Ayarlar
        </h2>
        <button @click="$emit('close')" class="close-btn">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M18 6 6 18"></path>
            <path d="m6 6 12 12"></path>
          </svg>
        </button>
      </div>
      
      <!-- Tab Navigation -->
      <div class="settings-tabs">
        <button 
          v-for="tab in tabs" 
          :key="tab.id"
          :class="['tab-btn', { active: activeTab === tab.id }]"
          @click="activeTab = tab.id"
        >
          <component :is="tab.icon" />
          {{ tab.label }}
        </button>
      </div>
      
      <div class="modal-body">
        <!-- Genel Tab -->
        <div v-if="activeTab === 'general'" class="settings-section">
          <div class="setting-group">
            <label class="setting-label">Sayfa Başına Öğe</label>
            <p class="setting-desc">Listelerde gösterilecek öğe sayısı</p>
            <AutocompleteSelect
              :model-value="settings.itemsPerPage"
              @update:model-value="v => { settings.itemsPerPage = Number(v); saveSettings() }"
              :items="itemsPerPageOptions"
              placeholder="Seçin..."
            />
          </div>
          
          <div class="setting-group">
            <label class="setting-label">Varsayılan Birim</label>
            <p class="setting-desc">Yeni ürünler için varsayılan birim</p>
            <AutocompleteSelect
              v-model="settings.defaultUnit"
              :items="unitOptions"
              placeholder="Birim seçin..."
              @update:model-value="saveSettings"
            />
          </div>
          
          <div class="setting-group">
            <label class="setting-toggle">
              <input 
                type="checkbox" 
                v-model="settings.autoDeductStock"
                @change="saveSettings"
              />
              <span class="toggle-slider"></span>
              <span class="toggle-label">
                <strong>Otomatik Stok Düşürme</strong>
                <small>Sipariş kaydedildiğinde stoktan otomatik düş</small>
              </span>
            </label>
          </div>
        </div>
        
        <!-- Görünüm Tab -->
        <div v-if="activeTab === 'appearance'" class="settings-section">
          <div class="setting-group">
            <label class="setting-label">Tema</label>
            <p class="setting-desc">Uygulama renk teması</p>
            <div class="theme-options">
              <button 
                :class="['theme-option', { active: settings.theme === 'light' }]"
                @click="setTheme('light')"
              >
                <span class="theme-icon">☀️</span>
                <span>Aydınlık</span>
              </button>
              <button 
                :class="['theme-option', { active: settings.theme === 'dark' }]"
                @click="setTheme('dark')"
              >
                <span class="theme-icon">🌙</span>
                <span>Karanlık</span>
              </button>
            </div>
          </div>
        </div>
        
        <!-- Hakkında Tab -->
        <div v-if="activeTab === 'about'" class="settings-section">
          <div class="about-header">
            <div class="app-icon">
              <img src="/favicon.ico" alt="AutoManagement" width="48" height="48" />
            </div>
            <div class="app-info">
              <h3>AutoManagement</h3>
              <p class="app-subtitle">Oto Yönetim Sistemi</p>
              <p class="app-version">
                <span class="version-badge">v{{ currentVersion }}</span>
                <template v-if="currentStatus && currentStatus.length">
                  <span v-for="s in currentStatus" :key="s" :class="['version-label', s]">{{ getStatusLabel(s) }}</span>
                </template>
                <span v-if="currentDate" class="changelog-date">{{ currentDate }}</span>
              </p>
            </div>
          </div>
          
          <div class="about-links">
            <a href="https://github.com/Ferhatduran55" target="_blank" class="about-link">
              <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>
              Ferhat Duran
            </a>
            <a href="https://github.com/Ferhatduran55/go-auto-aom" target="_blank" class="about-link">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M15 22v-4a4.8 4.8 0 0 0-1-3.5c3 0 6-2 6-5.5.08-1.25-.27-2.48-1-3.5.28-1.15.28-2.35 0-3.5 0 0-1 0-3 1.5-2.64-.5-5.36-.5-8 0C6 2 5 2 5 2c-.3 1.15-.3 2.35 0 3.5A5.403 5.403 0 0 0 4 9c0 3.5 3 5.5 6 5.5-.39.49-.68 1.05-.85 1.65-.17.6-.22 1.23-.15 1.85v4"/>
                <path d="M9 18c-4.51 2-5-2-7-2"/>
              </svg>
              GitHub
            </a>
          </div>
          
          <div class="changelog">
            <h4>Değişiklik Günlüğü</h4>
            <div v-for="release in changelog" :key="release.version" class="changelog-item">
              <div class="changelog-header">
                <span class="changelog-version">
                  <template v-if="release.version !== 'İlk Sürüm'">v</template>{{ release.version }}
                </span>
                <template v-if="Array.isArray(release.status)">
                  <span v-for="s in release.status" :key="s" :class="['status-badge', s]">{{ getStatusLabel(s) }}</span>
                </template>
                <span v-else :class="['status-badge', release.status]">{{ getStatusLabel(release.status) }}</span>
                <span class="changelog-date">{{ release.date }}</span>
              </div>
              <ul class="changelog-changes">
                <li v-for="(change, i) in release.changes" :key="i">{{ change }}</li>
              </ul>
            </div>
          </div>
          
          <div class="about-footer">
            <p>© 2025 Durasoft • MIT License</p>
            <p>🔓 Open Source</p>
          </div>
        </div>
        
        <!-- Geliştirici Tab -->
        <div v-if="activeTab === 'developer'" class="settings-section">
          <div class="developer-warning">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"></path>
              <path d="M12 9v4"></path>
              <path d="M12 17h.01"></path>
            </svg>
            <div>
              <strong>Geliştirici Seçenekleri</strong>
              <p>Bu ayarlar geliştiriciler içindir. Yanlış kullanım uygulamanın performansını etkileyebilir.</p>
            </div>
          </div>
          
          <div class="setting-group">
            <label class="setting-toggle">
              <input 
                type="checkbox" 
                v-model="settings.developerMode"
                @change="handleDeveloperModeChange"
              />
              <span class="toggle-slider"></span>
              <span class="toggle-label">
                <strong>Geliştirici Modu</strong>
                <small>F12 ile DevTools, sağ tık menüsü (yeniden başlatma gerektirir)</small>
              </span>
            </label>
          </div>
          
          <!-- Restart Required Notice -->
          <div v-if="restartRequired" class="restart-notice">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 12a9 9 0 1 1-9-9c2.52 0 4.93 1 6.74 2.74L21 8"/>
              <path d="M21 3v5h-5"/>
            </svg>
            <span>Değişikliklerin etkili olması için uygulamayı yeniden başlatın</span>
          </div>
          
          <div v-if="settings.developerMode" class="developer-info">
            <h4>Aktif Geliştirici Özellikleri:</h4>
            <ul>
              <li>✓ F12 ile DevTools açma</li>
              <li>✓ Sağ tık ile Inspect Element</li>
              <li>✓ Konsol debug logları</li>
              <li>✓ API çağrı detayları</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, h, onMounted } from 'vue'
import { useSettings } from '../composables/useSettings'
import { api } from '../api'
import AutocompleteSelect from './AutocompleteSelect.vue'

defineEmits(['close'])

const { settings, changelog, saveSettings, setTheme } = useSettings()

const activeTab = ref('general')
// Derive current version and status from changelog (first entry)
const currentVersion = changelog && changelog.length > 0 ? changelog[0].version : 'N/A'
const currentStatus = changelog && changelog.length > 0 ? (Array.isArray(changelog[0].status) ? changelog[0].status : [changelog[0].status]) : []
const currentDate = changelog && changelog.length > 0 ? changelog[0].date : ''
const restartRequired = ref(false)

// Load developer mode from backend on mount
onMounted(async () => {
  try {
    const result = await api.getDeveloperMode()
    if (result && typeof result.enabled === 'boolean') {
      settings.value.developerMode = result.enabled
    }
  } catch (e) {
    console.error('Failed to load developer mode:', e)
  }
})

// Handle developer mode change
async function handleDeveloperModeChange() {
  saveSettings()
  try {
    const result = await api.setDeveloperMode(settings.value.developerMode)
    if (result && result.restart_required) {
      restartRequired.value = true
    }
  } catch (e) {
    console.error('Failed to save developer mode:', e)
  }
}

// Options for autocomplete
const itemsPerPageOptions = [
  { label: '10', value: 10 },
  { label: '25', value: 25 },
  { label: '50', value: 50 },
  { label: '100', value: 100 }
]

const unitOptions = [
  { label: 'Adet', value: 'adet' },
  { label: 'Litre', value: 'litre' },
  { label: 'Kutu', value: 'kutu' },
  { label: 'Paket', value: 'paket' }
]

// Status label helper
function getStatusLabel(status) {
  const labels = {
    'release': 'Son Sürüm',
    'stable': 'Kararlı',
    'unstable': 'Geliştirme',
    'latest': 'En Son',
    'compatible': 'Uyumlu',
    'incompatible': 'Uyumsuz',
    'first': 'İlk Sürüm',
    'pre-release': 'Ön Sürüm'
  }
  return labels[status] || status
}

// Tab configuration with inline SVG icons
const tabs = [
  { 
    id: 'general', 
    label: 'Genel',
    icon: {
      render: () => h('svg', { 
        xmlns: 'http://www.w3.org/2000/svg', 
        width: 18, 
        height: 18, 
        viewBox: '0 0 24 24', 
        fill: 'none', 
        stroke: 'currentColor', 
        'stroke-width': 2 
      }, [
        h('path', { d: 'M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z' }),
        h('circle', { cx: 12, cy: 12, r: 3 })
      ])
    }
  },
  { 
    id: 'appearance', 
    label: 'Görünüm',
    icon: {
      render: () => h('svg', { 
        xmlns: 'http://www.w3.org/2000/svg', 
        width: 18, 
        height: 18, 
        viewBox: '0 0 24 24', 
        fill: 'none', 
        stroke: 'currentColor', 
        'stroke-width': 2 
      }, [
        h('circle', { cx: 12, cy: 12, r: 4 }),
        h('path', { d: 'M12 2v2' }),
        h('path', { d: 'M12 20v2' }),
        h('path', { d: 'm4.93 4.93 1.41 1.41' }),
        h('path', { d: 'm17.66 17.66 1.41 1.41' }),
        h('path', { d: 'M2 12h2' }),
        h('path', { d: 'M20 12h2' }),
        h('path', { d: 'm6.34 17.66-1.41 1.41' }),
        h('path', { d: 'm19.07 4.93-1.41 1.41' })
      ])
    }
  },
  { 
    id: 'about', 
    label: 'Hakkında',
    icon: {
      render: () => h('svg', { 
        xmlns: 'http://www.w3.org/2000/svg', 
        width: 18, 
        height: 18, 
        viewBox: '0 0 24 24', 
        fill: 'none', 
        stroke: 'currentColor', 
        'stroke-width': 2 
      }, [
        h('circle', { cx: 12, cy: 12, r: 10 }),
        h('path', { d: 'M12 16v-4' }),
        h('path', { d: 'M12 8h.01' })
      ])
    }
  },
  { 
    id: 'developer', 
    label: 'Geliştirici',
    icon: {
      render: () => h('svg', { 
        xmlns: 'http://www.w3.org/2000/svg', 
        width: 18, 
        height: 18, 
        viewBox: '0 0 24 24', 
        fill: 'none', 
        stroke: 'currentColor', 
        'stroke-width': 2 
      }, [
        h('path', { d: 'm18 16 4-4-4-4' }),
        h('path', { d: 'm6 8-4 4 4 4' }),
        h('path', { d: 'm14.5 4-5 16' })
      ])
    }
  }
]
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 1rem;
}

.settings-modal {
  background: var(--bg-card);
  border-radius: 1rem;
  width: 100%;
  max-width: 600px;
  max-height: 85vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
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
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.modal-header h2 svg {
  color: var(--accent-color);
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
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.settings-tabs {
  display: flex;
  gap: 0.25rem;
  padding: 0.75rem 1.5rem;
  border-bottom: 1px solid var(--border-color);
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 0.375rem;
  background: transparent;
  color: var(--text-muted);
  font-size: 0.875rem;
  font-weight: 500;
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

.modal-body {
  flex: 1;
  overflow-y: auto;
  padding: 1.5rem;
}

.settings-section {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.setting-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.setting-label {
  font-weight: 600;
  color: var(--text-primary);
}

.setting-desc {
  font-size: 0.875rem;
  color: var(--text-muted);
  margin: 0;
}

.setting-toggle {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
  cursor: pointer;
}

.setting-toggle input {
  display: none;
}

.toggle-slider {
  position: relative;
  width: 44px;
  height: 24px;
  background: var(--border-color);
  border-radius: 9999px;
  transition: background 0.2s;
  flex-shrink: 0;
}

.toggle-slider::after {
  content: '';
  position: absolute;
  top: 2px;
  left: 2px;
  width: 20px;
  height: 20px;
  background: white;
  border-radius: 50%;
  transition: transform 0.2s;
}

.setting-toggle input:checked + .toggle-slider {
  background: var(--accent-color);
}

.setting-toggle input:checked + .toggle-slider::after {
  transform: translateX(20px);
}

.toggle-label {
  display: flex;
  flex-direction: column;
  gap: 0.125rem;
}

.toggle-label strong {
  color: var(--text-primary);
}

.toggle-label small {
  color: var(--text-muted);
  font-size: 0.8125rem;
}

.theme-options {
  display: flex;
  gap: 0.75rem;
}

.theme-option {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  padding: 1rem;
  border: 2px solid var(--border-color);
  border-radius: 0.75rem;
  background: var(--bg-secondary);
  color: var(--text-primary);
  cursor: pointer;
  transition: all 0.15s;
}

.theme-option:hover {
  border-color: var(--accent-color);
}

.theme-option.active {
  border-color: var(--accent-color);
  background: var(--accent-color);
  color: white;
}

.theme-icon {
  font-size: 1.5rem;
}

/* About Tab Styles */
.about-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding-bottom: 1.5rem;
  border-bottom: 1px solid var(--border-color);
}

.app-icon {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.app-icon img {
  width: 48px;
  height: 48px;
  object-fit: contain;
}

.app-info h3 {
  margin: 0 0 0.25rem;
  font-size: 1.25rem;
  color: var(--text-primary);
}

.app-version {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.version-badge {
  background: var(--success-bg);
  color: var(--success-color);
  padding: 0.25rem 0.5rem;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  font-weight: 600;
}

.version-label {
  padding: 0.2rem 0.5rem;
  border-radius: 0.25rem;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
}

.version-label.release {
  background: linear-gradient(135deg, #8b5cf6, #7c3aed);
  color: white;
}

.version-label.stable {
  background: linear-gradient(135deg, #10b981, #059669);
  color: white;
}

.version-label.unstable {
  background: linear-gradient(135deg, #ef4444, #dc2626);
  color: white;
}

.version-label.latest {
  background: linear-gradient(135deg, #3b82f6, #2563eb);
  color: white;
}

.version-label.compatible {
  background: linear-gradient(135deg, #14b8a6, #0d9488);
  color: white;
}

.version-label.incompatible {
  background: linear-gradient(135deg, #f97316, #ea580c);
  color: white;
}

.version-label.first {
  background: linear-gradient(135deg, #6b7280, #4b5563);
  color: white;
}

.version-label.pre-release {
  background: linear-gradient(135deg, #f59e0b, #d97706);
  color: white;
}

.app-subtitle {
  color: var(--text-muted);
  font-size: 0.875rem;
  margin-bottom: 0.25rem;
}

.about-links {
  display: flex;
  gap: 1rem;
  padding: 1rem 0;
}

.about-link {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: var(--accent-color);
  text-decoration: none;
  font-weight: 500;
}

.about-link:hover {
  text-decoration: underline;
}

.changelog {
  margin-top: 1rem;
}

.changelog h4 {
  color: var(--text-primary);
  margin: 0 0 1rem;
  font-size: 1rem;
}

.changelog-item {
  padding: 1rem;
  background: var(--bg-secondary);
  border-radius: 0.5rem;
  margin-bottom: 0.75rem;
}

.changelog-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.75rem;
}

.changelog-version {
  background: var(--accent-color);
  color: white;
  padding: 0.25rem 0.5rem;
  border-radius: 0.25rem;
  font-size: 0.75rem;
  font-weight: 600;
}

.status-badge {
  padding: 0.2rem 0.5rem;
  border-radius: 0.25rem;
  font-size: 0.65rem;
  font-weight: 600;
  text-transform: uppercase;
}

.status-badge.release {
  background: linear-gradient(135deg, #8b5cf6, #7c3aed);
  color: white;
}

.status-badge.stable {
  background: linear-gradient(135deg, #10b981, #059669);
  color: white;
}

.status-badge.unstable {
  background: linear-gradient(135deg, #ef4444, #dc2626);
  color: white;
}

.status-badge.latest {
  background: linear-gradient(135deg, #3b82f6, #2563eb);
  color: white;
}

.status-badge.compatible {
  background: linear-gradient(135deg, #14b8a6, #0d9488);
  color: white;
}

.status-badge.incompatible {
  background: linear-gradient(135deg, #f97316, #ea580c);
  color: white;
}

.status-badge.first {
  background: linear-gradient(135deg, #6b7280, #4b5563);
  color: white;
}

.status-badge.pre-release {
  background: linear-gradient(135deg, #f59e0b, #d97706);
  color: white;
}

.changelog-date {
  color: var(--text-muted);
  font-size: 0.75rem;
  margin-left: auto;
}

.changelog-changes {
  margin: 0;
  padding-left: 1.25rem;
  color: var(--text-muted);
  font-size: 0.875rem;
}

.changelog-changes li {
  margin-bottom: 0.25rem;
}

.about-footer {
  margin-top: 1.5rem;
  padding-top: 1rem;
  border-top: 1px solid var(--border-color);
  text-align: center;
  color: var(--text-muted);
  font-size: 0.875rem;
}

.about-footer p {
  margin: 0.25rem 0;
}

/* Developer Tab Styles */
.developer-warning {
  display: flex;
  gap: 0.75rem;
  padding: 1rem;
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.1), rgba(217, 119, 6, 0.1));
  border: 1px solid rgba(245, 158, 11, 0.3);
  border-radius: 0.75rem;
  margin-bottom: 1.5rem;
}

.developer-warning svg {
  color: #f59e0b;
  flex-shrink: 0;
  margin-top: 0.125rem;
}

.developer-warning strong {
  color: #f59e0b;
  display: block;
  margin-bottom: 0.25rem;
}

.developer-warning p {
  margin: 0;
  color: var(--text-muted);
  font-size: 0.813rem;
  line-height: 1.4;
}

.developer-info {
  margin-top: 1.5rem;
  padding: 1rem;
  background: rgba(16, 185, 129, 0.1);
  border: 1px solid rgba(16, 185, 129, 0.3);
  border-radius: 0.75rem;
}

.developer-info h4 {
  margin: 0 0 0.75rem 0;
  color: #10b981;
  font-size: 0.875rem;
}

.developer-info ul {
  margin: 0;
  padding-left: 0;
  list-style: none;
}

.developer-info li {
  color: var(--text-secondary);
  font-size: 0.813rem;
  padding: 0.25rem 0;
}

.restart-notice {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-top: 1rem;
  padding: 0.875rem 1rem;
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.15), rgba(37, 99, 235, 0.15));
  border: 1px solid rgba(59, 130, 246, 0.4);
  border-radius: 0.75rem;
  color: #60a5fa;
  font-size: 0.875rem;
  animation: pulse-border 2s ease-in-out infinite;
}

.restart-notice svg {
  flex-shrink: 0;
  animation: spin 2s linear infinite;
}

@keyframes pulse-border {
  0%, 100% { border-color: rgba(59, 130, 246, 0.4); }
  50% { border-color: rgba(59, 130, 246, 0.8); }
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>
