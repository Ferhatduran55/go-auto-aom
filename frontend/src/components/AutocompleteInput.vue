<template>
  <div class="autocomplete-container" ref="containerRef">
    <input
      type="text"
      :value="modelValue"
      :placeholder="placeholder"
      :class="inputClass"
      @input="onInput"
      @focus="onFocus"
      @blur="onBlur"
      @keydown.down.prevent="navigateDown"
      @keydown.up.prevent="navigateUp"
      @keydown.enter.prevent="selectHighlighted"
      @keydown.escape="closeSuggestions"
      ref="inputRef"
    />
    
    <Transition name="fade">
      <div 
        v-if="showSuggestions && (suggestions.length > 0 || isLoading)" 
        class="suggestions-dropdown"
        ref="dropdownRef"
      >
        <div 
          v-if="isLoading && suggestions.length === 0" 
          class="suggestion-loading"
        >
          <span class="loading-spinner"></span>
          Yükleniyor...
        </div>
        
        <div
          v-for="(item, index) in suggestions"
          :key="getItemKey(item, index)"
          :class="['suggestion-item', { highlighted: highlightedIndex === index }]"
          @mousedown.prevent="selectItem(item)"
          @mouseenter="highlightedIndex = index"
        >
          <slot name="item" :item="item" :index="index">
            {{ getDisplayValue(item) }}
          </slot>
        </div>
        
        <!-- Lazy load sentinel -->
        <div 
          v-if="hasMore && !isLoading" 
          ref="sentinelRef" 
          class="load-more-sentinel"
        >
          <span class="loading-spinner small"></span>
        </div>
        
        <div v-if="isLoading && suggestions.length > 0" class="suggestion-loading-more">
          <span class="loading-spinner small"></span>
          Daha fazla yükleniyor...
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  searchFn: {
    type: Function,
    required: true
  },
  displayKey: {
    type: [String, Function],
    default: 'name'
  },
  valueKey: {
    type: String,
    default: null
  },
  itemKey: {
    type: String,
    default: 'id'
  },
  placeholder: {
    type: String,
    default: 'Ara...'
  },
  inputClass: {
    type: String,
    default: 'form-input'
  },
  pageSize: {
    type: Number,
    default: 10
  },
  debounceMs: {
    type: Number,
    default: 300
  },
  minChars: {
    type: Number,
    default: 0
  },
  showOnFocus: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['update:modelValue', 'select'])

// Refs
const containerRef = ref(null)
const inputRef = ref(null)
const dropdownRef = ref(null)
const sentinelRef = ref(null)

// State
const suggestions = ref([])
const showSuggestions = ref(false)
const isLoading = ref(false)
const highlightedIndex = ref(-1)
const currentPage = ref(1)
const hasMore = ref(true)
const currentQuery = ref('')
let debounceTimer = null
let observer = null

// Computed
const getDisplayValue = computed(() => {
  return (item) => {
    if (typeof props.displayKey === 'function') {
      return props.displayKey(item)
    }
    return item[props.displayKey] || item
  }
})

function getItemKey(item, index) {
  if (item && typeof item === 'object' && item[props.itemKey]) {
    return item[props.itemKey]
  }
  return index
}

// Methods
async function search(query, page = 1, append = false) {
  if (query.length < props.minChars && props.minChars > 0) {
    suggestions.value = []
    hasMore.value = false
    return
  }

  isLoading.value = true
  currentQuery.value = query
  
  try {
    const result = await props.searchFn(query, page, props.pageSize)
    
    // Handle different result formats
    let items = []
    let total = 0
    
    if (Array.isArray(result)) {
      items = result
      total = result.length
      hasMore.value = false
    } else if (result && typeof result === 'object') {
      items = result.items || result.products || result.data || []
      total = result.total || result.count || items.length
      hasMore.value = (page * props.pageSize) < total
    }
    
    if (append) {
      suggestions.value = [...suggestions.value, ...items]
    } else {
      suggestions.value = items
    }
    
    currentPage.value = page
    highlightedIndex.value = -1
  } catch (error) {
    console.error('Autocomplete search error:', error)
    suggestions.value = []
    hasMore.value = false
  } finally {
    isLoading.value = false
  }
}

function loadMore() {
  if (!isLoading.value && hasMore.value) {
    search(currentQuery.value, currentPage.value + 1, true)
  }
}

function onInput(e) {
  const value = e.target.value
  emit('update:modelValue', value)
  
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    currentPage.value = 1
    hasMore.value = true
    search(value, 1, false)
  }, props.debounceMs)
}

function onFocus() {
  if (props.showOnFocus) {
    showSuggestions.value = true
    if (suggestions.value.length === 0) {
      search(props.modelValue || '', 1, false)
    }
  }
}

function onBlur() {
  // Delay to allow click on suggestion
  setTimeout(() => {
    showSuggestions.value = false
    highlightedIndex.value = -1
  }, 200)
}

function navigateDown() {
  if (!showSuggestions.value) {
    showSuggestions.value = true
    return
  }
  if (highlightedIndex.value < suggestions.value.length - 1) {
    highlightedIndex.value++
    scrollToHighlighted()
  }
}

function navigateUp() {
  if (highlightedIndex.value > 0) {
    highlightedIndex.value--
    scrollToHighlighted()
  }
}

function scrollToHighlighted() {
  nextTick(() => {
    if (dropdownRef.value) {
      const highlighted = dropdownRef.value.querySelector('.highlighted')
      if (highlighted) {
        highlighted.scrollIntoView({ block: 'nearest' })
      }
    }
  })
}

function selectHighlighted() {
  if (highlightedIndex.value >= 0 && highlightedIndex.value < suggestions.value.length) {
    selectItem(suggestions.value[highlightedIndex.value])
  }
}

function selectItem(item) {
  const displayValue = getDisplayValue.value(item)
  const emitValue = props.valueKey ? item[props.valueKey] : displayValue
  
  emit('update:modelValue', emitValue)
  emit('select', item)
  
  showSuggestions.value = false
  highlightedIndex.value = -1
}

function closeSuggestions() {
  showSuggestions.value = false
  highlightedIndex.value = -1
}

// Setup IntersectionObserver for lazy loading
function setupObserver() {
  if (observer) {
    observer.disconnect()
  }
  
  observer = new IntersectionObserver(
    (entries) => {
      entries.forEach(entry => {
        if (entry.isIntersecting) {
          loadMore()
        }
      })
    },
    {
      root: dropdownRef.value,
      rootMargin: '50px',
      threshold: 0.1
    }
  )
}

watch(sentinelRef, (newVal) => {
  if (newVal && observer) {
    observer.observe(newVal)
  }
}, { immediate: true })

watch(showSuggestions, (visible) => {
  if (visible) {
    nextTick(() => {
      setupObserver()
      if (sentinelRef.value) {
        observer.observe(sentinelRef.value)
      }
    })
  }
})

onMounted(() => {
  setupObserver()
})

onUnmounted(() => {
  if (observer) {
    observer.disconnect()
  }
  clearTimeout(debounceTimer)
})

// Expose methods
defineExpose({
  focus: () => inputRef.value?.focus(),
  clear: () => {
    emit('update:modelValue', '')
    suggestions.value = []
  }
})
</script>

<style scoped>
.autocomplete-container {
  position: relative;
  width: 100%;
}

.suggestions-dropdown {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  margin-top: 4px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 0.5rem;
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.2);
  max-height: 280px;
  overflow-y: auto;
  z-index: 100;
}

.suggestion-item {
  padding: 0.75rem 1rem;
  cursor: pointer;
  transition: background 0.15s;
  color: var(--text-primary);
  border-bottom: 1px solid var(--border-color);
}

.suggestion-item:last-child {
  border-bottom: none;
}

.suggestion-item:hover,
.suggestion-item.highlighted {
  background: var(--bg-tertiary);
}

.suggestion-loading,
.suggestion-loading-more {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 1rem;
  color: var(--text-muted);
  font-size: 0.875rem;
}

.load-more-sentinel {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0.5rem;
  height: 40px;
}

.loading-spinner {
  width: 20px;
  height: 20px;
  border: 2px solid var(--border-color);
  border-top-color: var(--accent-color);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

.loading-spinner.small {
  width: 16px;
  height: 16px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Transition */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s, transform 0.15s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>
