<template>
  <div class="developer-panel">
    <!-- Konsol Bölümü -->
    <div class="console-section">
      <div class="console-header">
        <h3>
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="4 17 10 11 4 5"/>
            <line x1="12" y1="19" x2="20" y2="19"/>
          </svg>
          {{ t('developerConsole.title') }}
        </h3>
        <div class="console-actions">
          <button @click="refreshBackendLogs" class="action-btn" :title="t('common.refresh')">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 12a9 9 0 0 0-9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"/>
              <path d="M3 3v5h5"/>
            </svg>
          </button>
          <button @click="clearAllLogs" class="action-btn" :title="t('common.clear')">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/>
            </svg>
          </button>
          <button @click="togglePause" class="action-btn" :class="{ active: paused }" :title="paused ? t('developerConsole.resume') : t('developerConsole.pause')">
            <svg v-if="paused" xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polygon points="5 3 19 12 5 21 5 3"/>
            </svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="6" y="4" width="4" height="16"/><rect x="14" y="4" width="4" height="16"/>
            </svg>
          </button>
          <select v-model="sourceFilter" class="filter-select">
            <option value="all">{{ t('developerConsole.all') }}</option>
            <option value="System">{{ t('developerConsole.system') }}</option>
            <option value="frontend">{{ t('developerConsole.frontend') }}</option>
          </select>
          <select v-model="logLevel" class="filter-select">
            <option value="all">{{ t('developerConsole.level') }}</option>
            <option value="error">{{ t('developerConsole.error') }}</option>
            <option value="warn">{{ t('developerConsole.warn') }}</option>
            <option value="info">{{ t('developerConsole.info') }}</option>
            <option value="debug">{{ t('developerConsole.debug') }}</option>
          </select>
        </div>
      </div>
      
      <div class="console-body" ref="consoleBody">
        <div v-if="filteredLogs.length === 0" class="empty-console">
          <p>{{ t('developerConsole.emptyConsole') }}</p>
        </div>
        <div 
          v-for="(log, i) in filteredLogs" 
          :key="log.id || i" 
          :class="['log-entry', log.type, { 'backend': log.isBackend }]"
        >
          <span class="log-time">{{ formatTime(log.timestamp) }}</span>
          <span class="log-source" :class="log.source?.toLowerCase()">{{ log.source || 'FE' }}</span>
          <span class="log-type-badge" :class="log.type">{{ log.type.charAt(0).toUpperCase() }}</span>
          <span class="log-message">{{ getLogMessage(log) }}</span>
        </div>
      </div>
      
      <div class="console-footer">
        <span>{{ filteredLogs.length }}/{{ logs.length }} {{ t('developerConsole.records') }}</span>
        <span :class="['status', { paused }]">{{ paused ? '⏸️' : '🔴' }}</span>
      </div>
    </div>
    
    <!-- Sistem Bilgileri Bölümü -->
    <div class="system-section">
      <div class="system-header">
        <h3>
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="3"/><path d="M12 1v6m0 6v10"/><path d="m5.6 5.6 4.3 4.3m4.3 4.3 4.2 4.2"/><path d="M1 12h6m6 0h10"/><path d="m5.6 18.4 4.3-4.3m4.3-4.3 4.2-4.2"/>
          </svg>
          {{ t('developerConsole.systemInfo') }}
        </h3>
        <button @click="refreshSystemInfo(true)" class="action-btn" :title="t('common.refresh')">
          <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 12a9 9 0 0 0-9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"/><path d="M3 3v5h5"/>
          </svg>
        </button>
      </div>
      
      <div class="system-grid" v-if="systemInfo">
        <div class="info-card">
          <span class="info-label">{{ t('developerConsole.application') }}</span>
          <span class="info-value">{{ systemInfo.app_version }}</span>
        </div>
        <div class="info-card">
          <span class="info-label">Vue</span>
          <span class="info-value">{{ vueVersion }}</span>
        </div>
        <div class="info-card">
          <span class="info-label">{{ t('developerConsole.processor') }}</span>
          <span class="info-value">{{ systemInfo.cpu_model }}</span>
        </div>
        <div class="info-card">
          <span class="info-label">Go</span>
          <span class="info-value">{{ systemInfo.go_version }}</span>
        </div>
        <div class="info-card">
          <span class="info-label">Bleve</span>
          <span class="info-value">{{ systemInfo.bleve_version }}</span>
        </div>
        <div class="info-card">
          <span class="info-label">WebView</span>
          <span class="info-value">{{ systemInfo.webview_version }}</span>
        </div>
        <div class="info-card">
          <span class="info-label">{{ t('developerConsole.platform') }}</span>
          <span class="info-value">{{ systemInfo.os }}/{{ systemInfo.arch }}</span>
        </div>

        <div class="info-card">
          <span class="info-label">{{ t('developerConsole.uptime') }}</span>
          <span class="info-value">{{ uptimeDisplay }}</span>
        </div>
      </div>
      <div v-else class="system-loading">{{ t('common.loading') }}</div>
      
      <!-- Aksiyon Butonları -->
      <div class="system-actions">
        <button @click="handleExportData" class="sys-btn" :disabled="loading.export">
          <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/>
          </svg>
          {{ loading.export ? t('common.loading') : t('common.export') }}
        </button>
        <button @click="handleCopyInfo" class="sys-btn">
          <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
          </svg>
          {{ t('developerConsole.copySystemInfo') }}
        </button>
        <button @click="handleClearCache" class="sys-btn warning" :disabled="loading.cache">
          <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/>
          </svg>
          {{ loading.cache ? t('developerConsole.clearing') : t('developerConsole.clearCache') }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick, version as vueVersion } from 'vue'
import { api } from '../api'
import { useI18n } from '../i18n'

const { t, locale } = useI18n()

// Konsol state
const logs = ref([])
const paused = ref(false)
const logLevel = ref('all')
const sourceFilter = ref('all')
const consoleBody = ref(null)
const lastBackendLogId = ref(0)
let logIdCounter = 0
let backendPollInterval = null

// Sistem state - Global cache variable (outside component instance)
const cachedSystemInfo = ref(null)
const systemStartTime = ref(null)

const systemInfo = ref(null)
const clientUptime = ref(0)
let uptimeInterval = null

const loading = ref({
  export: false,
  cache: false
})

// Original console methods
const originalConsole = {
  log: console.log,
  info: console.info,
  warn: console.warn,
  error: console.error,
  debug: console.debug
}

// Add log entry
function addLog(type, source, message, isBackend = false, timestamp = null, code = null, params = null) {
  if (paused.value) return
  
  logs.value.push({
    id: ++logIdCounter,
    type,
    source,
    message,
    code,
    params,
    isBackend,
    timestamp: timestamp || new Date()
  })
  
  // Keep only last 1000 logs
  if (logs.value.length > 1000) {
    logs.value = logs.value.slice(-1000)
  }
  
  // Auto scroll
  nextTick(() => {
    if (consoleBody.value) {
      consoleBody.value.scrollTop = consoleBody.value.scrollHeight
    }
  })
}

// Intercept console
function interceptConsole() {
  const createInterceptor = (type) => (...args) => {
    originalConsole[type](...args)
    
    const message = args.map(arg => {
      if (typeof arg === 'object') {
        try { return JSON.stringify(arg, null, 2) } catch { return String(arg) }
      }
      return String(arg)
    }).join(' ')
    
    addLog(type, 'Frontend', message, false)
  }
  
  console.log = createInterceptor('log')
  console.info = createInterceptor('info')
  console.warn = createInterceptor('warn')
  console.error = createInterceptor('error')
  console.debug = createInterceptor('debug')
}

function restoreConsole() {
  Object.assign(console, originalConsole)
}

// Fetch backend logs
async function fetchBackendLogs() {
  if (paused.value) return
  try {
    const backendLogs = await api.getBackendLogsAfterId(lastBackendLogId.value)
    if (backendLogs?.length > 0) {
      backendLogs.forEach(log => {
        const ts = new Date(log.timestamp)
        addLog(log.level, log.source, log.message, true, ts, log.code, log.params)
        if (log.id > lastBackendLogId.value) {
          lastBackendLogId.value = log.id
        }
      })
    }
  } catch (e) { /* silent */ }
}

async function refreshBackendLogs() {
  try {
    const allLogs = await api.getBackendLogs()
    if (allLogs?.length > 0) {
      logs.value = logs.value.filter(l => !l.isBackend)
      allLogs.forEach(log => {
        const ts = new Date(log.timestamp)
        addLog(log.level, log.source, log.message, true, ts, log.code, log.params)
        if (log.id > lastBackendLogId.value) {
          lastBackendLogId.value = log.id
        }
      })
    }
  } catch (e) { originalConsole.error('[DevConsole] Error:', e) }
}

async function clearAllLogs() {
  logs.value = []
  logIdCounter = 0
  lastBackendLogId.value = 0
  try { await api.clearBackendLogs() } catch {}
}

const filteredLogs = computed(() => {
  return logs.value.filter(log => {
    if (logLevel.value !== 'all' && log.type !== logLevel.value) return false
    if (sourceFilter.value === 'frontend' && log.isBackend) return false
    if (sourceFilter.value !== 'all' && sourceFilter.value !== 'frontend' && log.source !== sourceFilter.value) return false
    return true
  })
})

function formatTime(date) {
  if (!(date instanceof Date)) date = new Date(date)
  const loc = (navigator.language) || locale.value || 'en-US'
  try {
    return date.toLocaleTimeString(loc, { hour: '2-digit', minute: '2-digit', second: '2-digit' })
  } catch {
    return date.toLocaleTimeString()
  }
}

function getLogMessage(log) {
  if (log && log.code) {
    const key = 'developerConsole.' + log.code
    const localized = t(key, log.params || {})
    if (localized && localized !== key) return localized
  }
  return log.message
}

function togglePause() { paused.value = !paused.value }

// Sistem fonksiyonları
async function refreshSystemInfo(force = false) {
  // Use cache if available and not forced
  if (!force && cachedSystemInfo.value) {
    systemInfo.value = cachedSystemInfo.value
    // If we have start time, calculate current uptime client-side
    if (systemStartTime.value) {
       updateClientUptime()
    }
    return
  }

  try {
    const info = await api.getSystemInfo()
    systemInfo.value = info
    cachedSystemInfo.value = info
    
    // Calculate start time based on initial uptime
    // startTime = Now - UptimeSeconds
    if (info.uptime_seconds) {
      systemStartTime.value = Date.now() - (info.uptime_seconds * 1000)
      clientUptime.value = info.uptime_seconds
    }
  } catch (e) {
    console.error("System info fetch error", e)
  }
}

function updateClientUptime() {
  if (systemStartTime.value) {
    clientUptime.value = Math.floor((Date.now() - systemStartTime.value) / 1000)
  }
}

// Uptime display uses client-side counter
const uptimeDisplay = computed(() => {
   return formatUptime(clientUptime.value, systemInfo.value?.uptime)
})

function formatUptime(seconds, fallback) {
  if (seconds === undefined || seconds === null) return fallback || '-'
  if (locale.value === 'en') {
    const s = Number(seconds)
    const h = Math.floor(s / 3600)
    const m = Math.floor((s % 3600) / 60)
    const ss = Math.floor(s % 60)
    return `${String(h).padStart(2,'0')}:${String(m).padStart(2,'0')}:${String(ss).padStart(2,'0')}`
  }
  return formatTimeHuman(seconds)
}

function formatTimeHuman(s) {
  const sec = Number(s)
  const h = Math.floor(sec / 3600)
  const m = Math.floor((sec % 3600) / 60)
  const ss = Math.floor(sec % 60)
  if (h > 0) return `${h}sa ${m}dk ${ss}sn`
  if (m > 0) return `${m}dk ${ss}sn`
  return `${ss}sn`
}

async function handleExportData() {
  loading.value.export = true
  try {
    const data = await api.exportAllData()
    if (data) {
      const blob = new Blob([data], { type: 'application/json' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `automanagement-export-${new Date().toISOString().slice(0,10)}.json`
      a.click()
      URL.revokeObjectURL(url)
      originalConsole.info('[System] ' + t('developerConsole.dataExported'))
    }
  } finally {
    loading.value.export = false
  }
}

async function handleCopyInfo() {
  const info = await api.copySystemInfo()
  if (info) {
    await navigator.clipboard.writeText(info)
    originalConsole.info('[System] ' + t('developerConsole.systemInfoCopied'))
  }
}

async function handleClearCache() {
  loading.value.cache = true
  try {
    const result = await api.clearCache()
    if (result.success) {
      originalConsole.info('[System] ' + t('developerConsole.cacheCleared'))
    } else {
      originalConsole.error('[System] ' + t('developerConsole.cacheClearError', { error: result.error }))
    }
  } finally {
    loading.value.cache = false
  }
}

onMounted(() => {
  interceptConsole()
  refreshBackendLogs()
  refreshSystemInfo() // Will use cache if available
  backendPollInterval = setInterval(fetchBackendLogs, 500)
  
  // Update uptime every second locally
  uptimeInterval = setInterval(updateClientUptime, 1000)
})

onUnmounted(() => {
  restoreConsole()
  if (backendPollInterval) clearInterval(backendPollInterval)
  if (uptimeInterval) clearInterval(uptimeInterval)
})
</script>

<style scoped>
.developer-panel {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  font-family: 'JetBrains Mono', 'Fira Code', 'Consolas', monospace;
}

/* Konsol Bölümü */
.console-section {
  background: #1a1a2e;
  border-radius: 0.5rem;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  height: 350px; /* Sabit yükseklik */
  max-height: 350px;
}

.console-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 0.75rem;
  background: #16213e;
  border-bottom: 1px solid #0f3460;
  flex-shrink: 0;
}

.console-header h3 {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin: 0;
  font-size: 0.8rem;
  font-weight: 600;
  color: #e94560;
}

.info-sub {
  font-size: 0.75rem;
  color: #9aa0a6;
  margin-left: 0.5rem;
}

.console-actions {
  display: flex;
  align-items: center;
  gap: 0.35rem;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 26px;
  height: 26px;
  background: transparent;
  border: 1px solid #0f3460;
  border-radius: 0.25rem;
  color: #8892b0;
  cursor: pointer;
  transition: all 0.15s;
}

.action-btn:hover {
  background: #0f3460;
  color: #e94560;
}

.action-btn.active {
  background: #e94560;
  color: white;
  border-color: #e94560;
}

.filter-select {
  padding: 0.25rem 0.5rem;
  background: #0f3460;
  border: 1px solid #0f3460;
  border-radius: 0.25rem;
  color: #ccd6f6;
  font-size: 0.7rem;
  cursor: pointer;
}

.console-body {
  flex: 1;
  overflow-y: auto;
  padding: 0.35rem;
  font-size: 0.7rem;
  line-height: 1.5;
}

.empty-console {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #4a5568;
  font-size: 0.75rem;
}

.log-entry {
  display: flex;
  gap: 0.5rem;
  padding: 0.15rem 0.35rem;
  border-radius: 0.2rem;
  align-items: flex-start;
}

.log-entry:hover {
  background: rgba(255, 255, 255, 0.03);
}

.log-entry.backend {
  border-left: 2px solid #e94560;
  padding-left: 0.5rem;
}

.log-time {
  color: #4a5568;
  flex-shrink: 0;
  font-size: 0.65rem;
}

.log-source {
  flex-shrink: 0;
  width: 28px;
  font-weight: 600;
  padding: 0 0.15rem;
  border-radius: 0.15rem;
  text-align: center;
  font-size: 0.6rem;
}

.log-source.frontend { background: rgba(100, 255, 218, 0.2); color: #64ffda; }
.log-source.ai { background: rgba(233, 69, 96, 0.2); color: #e94560; }
.log-source.system { background: rgba(100, 108, 255, 0.2); color: #646cff; }

.log-type-badge {
  flex-shrink: 0;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.55rem;
  font-weight: 700;
}

.log-type-badge.info { background: rgba(100, 255, 218, 0.2); color: #64ffda; }
.log-type-badge.warn { background: rgba(255, 217, 61, 0.2); color: #ffd93d; }
.log-type-badge.error { background: rgba(255, 107, 107, 0.2); color: #ff6b6b; }
.log-type-badge.debug { background: rgba(136, 146, 176, 0.2); color: #8892b0; }
.log-type-badge.log { background: rgba(204, 214, 246, 0.1); color: #ccd6f6; }

.log-entry.info { color: #64ffda; }
.log-entry.warn { color: #ffd93d; background: rgba(255, 217, 61, 0.05); }
.log-entry.error { color: #ff6b6b; background: rgba(255, 107, 107, 0.05); }
.log-entry.debug { color: #8892b0; opacity: 0.7; }
.log-entry.log { color: #ccd6f6; }

.log-message {
  white-space: pre-wrap;
  word-break: break-word;
  flex: 1;
}

.console-footer {
  display: flex;
  justify-content: space-between;
  padding: 0.35rem 0.75rem;
  background: #16213e;
  border-top: 1px solid #0f3460;
  font-size: 0.65rem;
  color: #8892b0;
  flex-shrink: 0;
}

.status.paused { color: #ffd93d; }

/* Sistem Bilgileri Bölümü */
.system-section {
  background: #1a1a2e;
  border-radius: 0.5rem;
  padding: 0.75rem;
}

.system-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.75rem;
}

.system-header h3 {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin: 0;
  font-size: 0.8rem;
  font-weight: 600;
  color: #646cff;
}

.system-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 0.5rem;
  margin-bottom: 0.75rem;
}

.info-card {
  background: #16213e;
  border-radius: 0.35rem;
  padding: 0.5rem 0.65rem;
  display: flex;
  flex-direction: column;
  gap: 0.15rem;
}

.info-card.full-width {
  grid-column: span 4;
}

.info-label {
  font-size: 0.6rem;
  color: #8892b0;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.info-value {
  font-size: 0.75rem;
  color: #ccd6f6;
  font-weight: 500;
}

.info-value.status-ok { color: #64ffda; }
.info-value.status-error { color: #ff6b6b; }

.system-loading {
  text-align: center;
  padding: 1rem;
  color: #8892b0;
  font-size: 0.75rem;
}

/* Aksiyon Butonları */
.system-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.sys-btn {
  display: flex;
  align-items: center;
  gap: 0.35rem;
  padding: 0.4rem 0.65rem;
  background: #0f3460;
  border: 1px solid #0f3460;
  border-radius: 0.35rem;
  color: #ccd6f6;
  font-size: 0.7rem;
  cursor: pointer;
  transition: all 0.15s;
}

.sys-btn:hover {
  background: #16213e;
  border-color: #646cff;
  color: #646cff;
}

.sys-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.sys-btn.warning:hover {
  border-color: #ff6b6b;
  color: #ff6b6b;
}

/* Scrollbar */
.console-body::-webkit-scrollbar {
  width: 6px;
}

.console-body::-webkit-scrollbar-track {
  background: #1a1a2e;
}

.console-body::-webkit-scrollbar-thumb {
  background: #0f3460;
  border-radius: 3px;
}

.console-body::-webkit-scrollbar-thumb:hover {
  background: #e94560;
}

/* Responsive */
@media (max-width: 768px) {
  .system-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  .info-card.full-width {
    grid-column: span 2;
  }
}
</style>
