<template>
  <div class="notes-container">
    <!-- Sol Panel: Not Listesi -->
    <div class="notes-sidebar card">
      <div class="notes-header">
        <h2 class="text-lg font-semibold">{{ t('notes.title') }}</h2>
        <button @click="handleNewNote" class="btn btn-primary btn-sm">
          {{ t('notes.newNote') }}
        </button>
      </div>
      
      <!-- Arama -->
      <div class="notes-search">
        <input 
          v-model="searchInput"
          @input="handleSearch"
          type="text" 
          class="form-input"
          :placeholder="t('notes.searchPlaceholder')"
        />
      </div>
      
      <!-- Not Listesi -->
      <div class="notes-list">
        <div v-if="loading" class="notes-loading">
          <div class="loading-spinner"></div>
          <span>{{ t('notes.loading') }}</span>
        </div>
        
        <div v-else-if="notes.length === 0" class="notes-empty">
          <span class="text-muted">{{ t('notes.empty') }}</span>
        </div>
        
        <div 
          v-else
          v-for="note in notes" 
          :key="note.id"
          @click="handleSelectNote(note)"
          :class="['note-item', { active: currentNote?.id === note.id }]"
        >
          <div class="note-item-title">{{ note.title || t('notes.untitled') }}</div>
          <div class="note-item-meta">
            <span v-if="note.customer_name" class="note-customer">{{ note.customer_name }}</span>
            <span class="note-date">{{ formatDate(note.updated_at) }}</span>
          </div>
          <div class="note-item-preview">{{ getPreview(note.content) }}</div>
        </div>
      </div>
    </div>
    
    <!-- Sağ Panel: Not Düzenleme -->
    <div class="notes-editor card">
      <template v-if="currentNote">
        <!-- Üst Bar -->
        <div class="editor-header">
          <input 
            v-model="currentNote.title"
            @input="markModified"
            type="text" 
            class="title-input"
            :placeholder="t('notes.titlePlaceholder')"
          />
          <div class="editor-actions">
            <button 
              @click="handleFormat" 
              class="btn btn-secondary btn-sm"
              :title="t('notes.formatTooltip')"
            >
              ⚡ {{ t('notes.format') }}
            </button>
            <button 
              @click="handleSave" 
              class="btn btn-primary btn-sm"
              :disabled="!hasUnsavedChanges"
            >
              💾 {{ t('notes.save') }}
            </button>
            <button 
              @click="handleDelete" 
              class="btn btn-danger btn-sm"
              :disabled="!currentNote.id"
            >
              🗑️ {{ t('notes.delete') }}
            </button>
          </div>
        </div>
        
        <!-- Müşteri Adı -->
        <div class="customer-row">
          <label>{{ t('notes.customer') }}</label>
          <input 
            v-model="currentNote.customer_name"
            @input="markModified"
            type="text" 
            class="form-input customer-input"
            :placeholder="t('notes.customerPlaceholder')"
          />
        </div>
        
        <!-- İçerik Editörü -->
        <div class="editor-content">
          <textarea 
            v-model="currentNote.content"
            @input="markModified"
            class="note-textarea"
            :placeholder="t('notes.contentPlaceholder')"
          ></textarea>
        </div>
        
        <!-- Durum Çubuğu -->
        <div class="editor-footer">
          <span v-if="hasUnsavedChanges" class="unsaved-indicator">
            ● {{ t('notes.unsavedChanges') }}
          </span>
          <span v-else class="saved-indicator">
            ✓ {{ t('notes.saved') }}
          </span>
          <span class="char-count">{{ t('notes.charCount', { count: currentNote.content?.length || 0 }) }}</span>
        </div>
      </template>
      
      <template v-else>
        <div class="editor-empty">
          <div class="empty-icon">📝</div>
          <div class="empty-text">{{ t('notes.selectOrCreate') }}</div>
          <button @click="handleNewNote" class="btn btn-primary">
            {{ t('notes.createNew') }}
          </button>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useNotes } from '@/composables/useNotes'
import { useToast } from '@/composables/useToast'
import { useI18n } from '@/i18n'

const { t, locale } = useI18n()
const { showToast } = useToast()
const {
  notes,
  currentNote,
  loading,
  hasUnsavedChanges,
  loadNotes,
  loadNote,
  newNote,
  saveNote,
  deleteNote,
  formatCurrentNote,
  search
} = useNotes()

const searchInput = ref('')
let searchTimeout = null

// Başlangıçta notları yükle
onMounted(async () => {
  await loadNotes()
})

// Arama (debounce ile)
function handleSearch() {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    search(searchInput.value)
  }, 300)
}

// Yeni not
function handleNewNote() {
  newNote()
}

// Not seç
async function handleSelectNote(note) {
  if (hasUnsavedChanges.value) {
    if (!confirm(t('notes.unsavedConfirm'))) {
      return
    }
  }
  await loadNote(note.id)
}

// Kaydet
async function handleSave() {
  const result = await saveNote()
  if (result.success) {
    showToast(t('notes.noteSaved'), 'success')
  } else {
    showToast(t('notes.saveError', { error: result.error || t('notes.unknownError') }), 'error')
  }
}

// Sil
async function handleDelete() {
  if (!currentNote.value?.id) return
  
  if (!confirm(t('notes.deleteConfirm'))) {
    return
  }
  
  const result = await deleteNote(currentNote.value.id)
  if (result.success) {
    showToast(t('notes.noteDeleted'), 'success')
  } else {
    showToast(t('notes.deleteError', { error: result.error || t('notes.unknownError') }), 'error')
  }
}

// Formatla
async function handleFormat() {
  await formatCurrentNote()
  showToast(t('notes.formatted'), 'success')
}

// Değişiklik işaretle
function markModified() {
  if (currentNote.value) {
    currentNote.value._modified = true
  }
}

// Tarih formatlama
function formatDate(dateStr) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  const now = new Date()
  const diffDays = Math.floor((now - date) / (1000 * 60 * 60 * 24))
  const loc = locale.value === 'tr' ? 'tr-TR' : 'en-US'
  
  if (diffDays === 0) {
    return date.toLocaleTimeString(loc, { hour: '2-digit', minute: '2-digit' })
  } else if (diffDays === 1) {
    return t('notes.yesterday')
  } else if (diffDays < 7) {
    return date.toLocaleDateString(loc, { weekday: 'short' })
  } else {
    return date.toLocaleDateString(loc, { day: 'numeric', month: 'short' })
  }
}

// Önizleme metni
function getPreview(content) {
  if (!content) return ''
  const preview = content.replace(/\n/g, ' ').substring(0, 80)
  return preview + (content.length > 80 ? '...' : '')
}
</script>

<style scoped>
.notes-container {
  display: grid;
  grid-template-columns: 320px 1fr;
  gap: 1.5rem;
  height: calc(100vh - 200px);
  min-height: 500px;
}

/* Sidebar */
.notes-sidebar {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.notes-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  border-bottom: 1px solid var(--border-color);
}

.notes-search {
  padding: 0.75rem 1rem;
  border-bottom: 1px solid var(--border-color);
}

.notes-list {
  flex: 1;
  overflow-y: auto;
  padding: 0.5rem;
}

.notes-loading,
.notes-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  gap: 0.5rem;
}

.loading-spinner {
  width: 24px;
  height: 24px;
  border: 2px solid var(--border-color);
  border-top-color: var(--accent);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.note-item {
  padding: 0.75rem;
  border-radius: 0.5rem;
  cursor: pointer;
  transition: background-color 0.15s;
  margin-bottom: 0.25rem;
}

.note-item:hover {
  background: var(--bg-secondary);
}

.note-item.active {
  background: var(--accent);
  color: white;
}

.note-item.active .note-item-meta,
.note-item.active .note-item-preview {
  color: rgba(255, 255, 255, 0.8);
}

.note-item-title {
  font-weight: 600;
  font-size: 0.9rem;
  margin-bottom: 0.25rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.note-item-meta {
  display: flex;
  gap: 0.5rem;
  font-size: 0.75rem;
  color: var(--text-muted);
  margin-bottom: 0.25rem;
}

.note-customer {
  background: var(--bg-secondary);
  padding: 0 0.25rem;
  border-radius: 0.25rem;
}

.note-item.active .note-customer {
  background: rgba(255, 255, 255, 0.2);
}

.note-item-preview {
  font-size: 0.8rem;
  color: var(--text-muted);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* Editor */
.notes-editor {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.editor-header {
  display: flex;
  gap: 1rem;
  align-items: center;
  padding: 1rem;
  border-bottom: 1px solid var(--border-color);
}

.title-input {
  flex: 1;
  background: transparent;
  border: none;
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-primary);
  outline: none;
}

.title-input::placeholder {
  color: var(--text-muted);
}

.editor-actions {
  display: flex;
  gap: 0.5rem;
}

.customer-row {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1rem;
  border-bottom: 1px solid var(--border-color);
}

.customer-row label {
  font-size: 0.85rem;
  color: var(--text-muted);
  white-space: nowrap;
}

.customer-input {
  flex: 1;
  padding: 0.5rem 0.75rem !important;
}

.editor-content {
  flex: 1;
  padding: 1rem;
  overflow: hidden;
}

.note-textarea {
  width: 100%;
  height: 100%;
  background: var(--bg-secondary);
  color: var(--text-primary);
  border: 1px solid var(--border-color);
  border-radius: 0.5rem;
  padding: 1rem;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 0.9rem;
  line-height: 1.6;
  resize: none;
  outline: none;
  white-space: pre-wrap;
  tab-size: 4;
}

.note-textarea:focus {
  border-color: var(--accent);
}

.note-textarea::placeholder {
  color: var(--text-muted);
  font-family: inherit;
}

.editor-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 1rem;
  border-top: 1px solid var(--border-color);
  font-size: 0.8rem;
}

.unsaved-indicator {
  color: var(--warning, #f59e0b);
}

.saved-indicator {
  color: var(--success, #22c55e);
}

.char-count {
  color: var(--text-muted);
}

.editor-empty {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  color: var(--text-muted);
}

.empty-icon {
  font-size: 4rem;
  opacity: 0.5;
}

.empty-text {
  font-size: 1.1rem;
}

/* Button Styles */
.btn-sm {
  padding: 0.4rem 0.75rem;
  font-size: 0.85rem;
}

.btn-danger {
  background: #ef4444;
  color: white;
}

.btn-danger:hover {
  background: #dc2626;
}

.btn-danger:disabled {
  background: #9ca3af;
  cursor: not-allowed;
}
</style>
