import { ref, computed } from 'vue'
import { api } from '@/api'

// Reaktif state
const notes = ref([])
const currentNote = ref(null)
const loading = ref(false)
const searchTerm = ref('')

// Computed
const filteredNotes = computed(() => {
  if (!searchTerm.value) return notes.value
  
  const term = searchTerm.value.toLowerCase()
  return notes.value.filter(note => 
    note.title?.toLowerCase().includes(term) ||
    note.content?.toLowerCase().includes(term) ||
    note.customer_name?.toLowerCase().includes(term)
  )
})

const hasUnsavedChanges = computed(() => {
  if (!currentNote.value) return false
  // Yeni not veya değişiklik varsa
  return currentNote.value.id === '' || currentNote.value._modified
})

// Notları yükle
async function loadNotes(search = '') {
  loading.value = true
  try {
    notes.value = await api.loadNotes(search)
  } catch (e) {
    console.error('Notlar yüklenemedi:', e)
    notes.value = []
  } finally {
    loading.value = false
  }
}

// Not yükle (ID ile)
async function loadNote(id) {
  try {
    const note = await api.loadNoteById(id)
    if (!note.error) {
      currentNote.value = { ...note, _modified: false }
    }
    return note
  } catch (e) {
    console.error('Not yüklenemedi:', e)
    return { error: e.message }
  }
}

// Yeni not oluştur
function newNote() {
  currentNote.value = {
    id: '',
    title: '',
    content: '',
    customer_name: '',
    _modified: false
  }
}

// Notu kaydet
async function saveNote() {
  if (!currentNote.value) return { error: 'Not yok' }
  
  try {
    const result = await api.saveNote({
      id: currentNote.value.id,
      title: currentNote.value.title || 'Başlıksız Not',
      content: currentNote.value.content,
      customer_name: currentNote.value.customer_name
    })
    
    if (result.success) {
      currentNote.value.id = result.id
      currentNote.value._modified = false
      await loadNotes() // Listeyi yenile
    }
    
    return result
  } catch (e) {
    console.error('Not kaydedilemedi:', e)
    return { error: e.message }
  }
}

// Notu sil
async function deleteNote(id) {
  try {
    const result = await api.deleteNote(id)
    if (result.success) {
      // Eğer silinen not açık ise kapat
      if (currentNote.value?.id === id) {
        currentNote.value = null
      }
      await loadNotes() // Listeyi yenile
    }
    return result
  } catch (e) {
    console.error('Not silinemedi:', e)
    return { error: e.message }
  }
}

// Metni formatla
async function formatText(text) {
  try {
    const result = await api.autoFormatText(text)
    return result.formatted || text
  } catch (e) {
    console.error('Formatlama hatası:', e)
    return text
  }
}

// Mevcut notun içeriğini formatla
async function formatCurrentNote() {
  if (!currentNote.value) return
  
  const formatted = await formatText(currentNote.value.content)
  currentNote.value.content = formatted
  currentNote.value._modified = true
}

// Not içeriğini güncelle
function updateContent(content) {
  if (currentNote.value) {
    currentNote.value.content = content
    currentNote.value._modified = true
  }
}

// Not başlığını güncelle
function updateTitle(title) {
  if (currentNote.value) {
    currentNote.value.title = title
    currentNote.value._modified = true
  }
}

// Müşteri adını güncelle
function updateCustomerName(name) {
  if (currentNote.value) {
    currentNote.value.customer_name = name
    currentNote.value._modified = true
  }
}

// Notu kapat
function closeNote() {
  currentNote.value = null
}

// Arama
async function search(term) {
  searchTerm.value = term
  await loadNotes(term)
}

// Export
export function useNotes() {
  return {
    // State
    notes,
    currentNote,
    loading,
    searchTerm,
    
    // Computed
    filteredNotes,
    hasUnsavedChanges,
    
    // Actions
    loadNotes,
    loadNote,
    newNote,
    saveNote,
    deleteNote,
    formatText,
    formatCurrentNote,
    updateContent,
    updateTitle,
    updateCustomerName,
    closeNote,
    search
  }
}
