<template>
  <div class="autocomplete-select" ref="wrapperRef">
    <div class="input-container" @click="toggleDropdown">
      <svg v-if="showSearchIcon" class="search-icon" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="11" cy="11" r="8"></circle>
        <path d="m21 21-4.35-4.35"></path>
      </svg>
      <input 
        type="text" 
        v-model="searchQuery"
        @input="onInput"
        @focus="onFocus"
        @keydown.down.prevent="navigateDown"
        @keydown.up.prevent="navigateUp"
        @keydown.enter.prevent="selectHighlighted"
        @keydown.escape="closeDropdown"
        @keydown.tab="closeDropdown"
        class="form-input"
        :class="{ 'with-icon': showSearchIcon, 'with-clear': showClearButton }"
        :placeholder="placeholder"
        :disabled="disabled"
        ref="inputRef"
      />
      <button 
        v-if="showClearButton && selectedItem" 
        @click.stop="clearSelection" 
        class="clear-btn"
        type="button"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M18 6 6 18"></path>
          <path d="m6 6 12 12"></path>
        </svg>
      </button>
      <svg class="chevron-icon" :class="{ 'rotate': isOpen }" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="m6 9 6 6 6-6"/>
      </svg>
    </div>
    
    <!-- Selected Item Display (optional) -->
    <div v-if="selectedItem && showSelectedBadge" class="selected-badge">
      <slot name="selected" :item="selectedItem">
        <span class="badge-text">{{ getDisplayText(selectedItem) }}</span>
      </slot>
    </div>
    
    <!-- Dropdown -->
    <Teleport to="body">
      <div 
        v-if="isOpen" 
        class="autocomplete-dropdown"
        :style="dropdownStyle"
        ref="dropdownRef"
      >
        <div v-if="filteredItems.length === 0" class="no-results">
          {{ noResultsText }}
        </div>
        <div 
          v-else
          v-for="(item, index) in filteredItems" 
          :key="getItemValue(item)"
          @click="selectItem(item)"
          @mouseenter="highlightedIndex = index"
          :class="['dropdown-item', { 
            'highlighted': index === highlightedIndex,
            'selected': isItemSelected(item),
            'disabled': isItemDisabled(item)
          }]"
        >
          <slot name="option" :item="item" :index="index">
            <div class="item-content">
              <span class="item-label">{{ getDisplayText(item) }}</span>
              <span v-if="getItemMeta(item)" class="item-meta">{{ getItemMeta(item) }}</span>
            </div>
          </slot>
          <svg v-if="isItemSelected(item)" class="check-icon" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="20 6 9 17 4 12"></polyline>
          </svg>
        </div>
        
        <!-- Add New Option -->
        <div v-if="allowCreate && searchQuery && !hasExactMatch" class="dropdown-item create-item" @click="createItem">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 5v14"></path>
            <path d="M5 12h14"></path>
          </svg>
          <span>"{{ searchQuery }}" ekle</span>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'

const props = defineProps({
  modelValue: {
    type: [String, Number, Object],
    default: null
  },
  items: {
    type: Array,
    default: () => []
  },
  labelKey: {
    type: String,
    default: 'label'
  },
  valueKey: {
    type: String,
    default: 'value'
  },
  metaKey: {
    type: String,
    default: null
  },
  placeholder: {
    type: String,
    default: 'Seçin veya arayın...'
  },
  noResultsText: {
    type: String,
    default: 'Sonuç bulunamadı'
  },
  disabled: {
    type: Boolean,
    default: false
  },
  showSearchIcon: {
    type: Boolean,
    default: true
  },
  showClearButton: {
    type: Boolean,
    default: true
  },
  showSelectedBadge: {
    type: Boolean,
    default: false
  },
  allowCreate: {
    type: Boolean,
    default: false
  },
  filterFn: {
    type: Function,
    default: null
  },
  disabledFn: {
    type: Function,
    default: null
  }
})

const emit = defineEmits(['update:modelValue', 'change', 'create', 'search'])

const wrapperRef = ref(null)
const inputRef = ref(null)
const dropdownRef = ref(null)
const isOpen = ref(false)
const searchQuery = ref('')
const highlightedIndex = ref(0)
const dropdownStyle = ref({})

// Get selected item from value
const selectedItem = computed(() => {
  if (props.modelValue === null || props.modelValue === undefined || props.modelValue === '') {
    return null
  }
  
  // If modelValue is an object, return it directly
  if (typeof props.modelValue === 'object') {
    return props.modelValue
  }
  
  // Find item by value
  return props.items.find(item => getItemValue(item) === props.modelValue) || null
})

// Filter items based on search
const filteredItems = computed(() => {
  let items = props.items
  
  if (props.filterFn) {
    items = items.filter(item => props.filterFn(item, searchQuery.value))
  } else if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    items = items.filter(item => {
      const label = getDisplayText(item).toLowerCase()
      const meta = getItemMeta(item)?.toLowerCase() || ''
      return label.includes(query) || meta.includes(query)
    })
  }
  
  return items.slice(0, 50) // Limit results
})

// Check if there's an exact match
const hasExactMatch = computed(() => {
  if (!searchQuery.value) return true
  const query = searchQuery.value.toLowerCase()
  return props.items.some(item => getDisplayText(item).toLowerCase() === query)
})

// Helper functions
function getDisplayText(item) {
  if (item === null || item === undefined) return ''
  if (typeof item === 'string') return item
  if (typeof item === 'object') return item[props.labelKey] || item.name || item.label || ''
  return String(item)
}

function getItemValue(item) {
  if (item === null || item === undefined) return null
  if (typeof item === 'string' || typeof item === 'number') return item
  if (typeof item === 'object') return item[props.valueKey] || item.id || item.value
  return item
}

function getItemMeta(item) {
  if (!props.metaKey || typeof item !== 'object') return null
  return item[props.metaKey]
}

function isItemSelected(item) {
  return getItemValue(item) === getItemValue(selectedItem.value)
}

function isItemDisabled(item) {
  if (props.disabledFn) return props.disabledFn(item)
  return item?.disabled === true
}

// Actions
function toggleDropdown(e) {
  if (props.disabled) return
  
  // If clicking on input, let onFocus handle it
  if (e && e.target === inputRef.value) {
    return
  }
  
  if (isOpen.value) {
    closeDropdown()
  } else {
    openDropdown()
  }
}

function openDropdown() {
  if (props.disabled || isOpen.value) return
  isOpen.value = true
  updateDropdownPosition()
  nextTick(() => inputRef.value?.focus())
}

function onFocus() {
  if (!isOpen.value) {
    isOpen.value = true
    updateDropdownPosition()
  }
}

function onInput() {
  isOpen.value = true
  highlightedIndex.value = 0
  emit('search', searchQuery.value)
  updateDropdownPosition()
}

function closeDropdown() {
  isOpen.value = false
  // Restore search query to selected item
  if (selectedItem.value) {
    searchQuery.value = getDisplayText(selectedItem.value)
  } else {
    searchQuery.value = ''
  }
}

function selectItem(item) {
  if (isItemDisabled(item)) return
  
  const value = getItemValue(item)
  searchQuery.value = getDisplayText(item)
  emit('update:modelValue', value)
  emit('change', item)
  isOpen.value = false
}

function clearSelection() {
  searchQuery.value = ''
  emit('update:modelValue', null)
  emit('change', null)
  inputRef.value?.focus()
}

function createItem() {
  const val = searchQuery.value
  emit('create', val)
  // Immediately update v-model so parent gets the value (helps UX when items aren't yet reloaded)
  emit('update:modelValue', val)
  // close dropdown but keep the typed query visible until parent updates items
  isOpen.value = false
}

function navigateDown() {
  if (!isOpen.value) {
    isOpen.value = true
    updateDropdownPosition()
    return
  }
  if (highlightedIndex.value < filteredItems.value.length - 1) {
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

function selectHighlighted() {
  if (filteredItems.value[highlightedIndex.value]) {
    selectItem(filteredItems.value[highlightedIndex.value])
  }
}

function scrollToHighlighted() {
  nextTick(() => {
    const dropdown = dropdownRef.value
    if (!dropdown) return
    const item = dropdown.children[highlightedIndex.value]
    if (item) {
      item.scrollIntoView({ block: 'nearest' })
    }
  })
}

function updateDropdownPosition() {
  nextTick(() => {
    if (!wrapperRef.value) return
    const rect = wrapperRef.value.getBoundingClientRect()
    const spaceBelow = window.innerHeight - rect.bottom
    const spaceAbove = rect.top
    const dropdownHeight = 300
    
    const showAbove = spaceBelow < dropdownHeight && spaceAbove > spaceBelow
    
    dropdownStyle.value = {
      position: 'fixed',
      left: `${rect.left}px`,
      width: `${rect.width}px`,
      zIndex: 9999,
      ...(showAbove
        ? { bottom: `${window.innerHeight - rect.top + 4}px` }
        : { top: `${rect.bottom + 4}px` }
      )
    }
  })
}

// Click outside handler
function handleClickOutside(e) {
  if (wrapperRef.value && !wrapperRef.value.contains(e.target) && 
      dropdownRef.value && !dropdownRef.value.contains(e.target)) {
    closeDropdown()
  }
}

// Update search query when modelValue changes
watch(() => props.modelValue, (newVal) => {
  if (selectedItem.value) {
    searchQuery.value = getDisplayText(selectedItem.value)
  } else if (newVal === null || newVal === undefined || newVal === '') {
    // only clear when model is explicitly emptied
    searchQuery.value = ''
  }
  // otherwise keep user's typed query (e.g., when creating a new item that isn't yet in the items list)
}, { immediate: true })

// Reset highlighted index when filtered items change
watch(filteredItems, () => {
  highlightedIndex.value = 0
})

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  window.addEventListener('scroll', updateDropdownPosition, true)
  window.addEventListener('resize', updateDropdownPosition)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  window.removeEventListener('scroll', updateDropdownPosition, true)
  window.removeEventListener('resize', updateDropdownPosition)
})
</script>

<style scoped>
.autocomplete-select {
  position: relative;
  width: 100%;
}

.input-container {
  position: relative;
  display: flex;
  align-items: center;
}

.input-container .form-input {
  width: 100%;
  padding-right: 2.5rem;
  cursor: pointer;
}

.input-container .form-input.with-icon {
  padding-left: 2.5rem;
}

.input-container .form-input.with-clear {
  padding-right: 4rem;
}

.search-icon {
  position: absolute;
  left: 0.75rem;
  color: var(--text-muted);
  pointer-events: none;
  z-index: 1;
}

.chevron-icon {
  position: absolute;
  right: 0.75rem;
  color: var(--text-muted);
  pointer-events: none;
  transition: transform 0.2s;
}

.chevron-icon.rotate {
  transform: rotate(180deg);
}

.clear-btn {
  position: absolute;
  right: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 1.25rem;
  height: 1.25rem;
  border: none;
  border-radius: 50%;
  background: var(--bg-tertiary);
  color: var(--text-muted);
  cursor: pointer;
  transition: all 0.15s;
  z-index: 1;
}

.clear-btn:hover {
  background: var(--danger-color);
  color: white;
}

.selected-badge {
  display: flex;
  align-items: center;
  padding: 0.5rem 0.75rem;
  margin-top: 0.5rem;
  background: var(--bg-secondary);
  border: 1px solid var(--accent-color);
  border-radius: 0.5rem;
}

.badge-text {
  font-weight: 500;
  color: var(--text-primary);
}

.autocomplete-dropdown {
  background: var(--bg-card);
  border: 2px solid var(--border-color);
  border-radius: 0.75rem;
  max-height: 300px;
  overflow-y: auto;
  box-shadow: 0 10px 40px -5px rgba(0, 0, 0, 0.3);
}

.dropdown-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.75rem 1rem;
  cursor: pointer;
  border-bottom: 1px solid var(--border-color);
  transition: all 0.1s;
  color: var(--text-primary);
}

.dropdown-item:last-child {
  border-bottom: none;
}

.dropdown-item:hover,
.dropdown-item.highlighted {
  background: var(--accent-color);
  color: white;
}

.dropdown-item.highlighted .item-meta,
.dropdown-item:hover .item-meta {
  color: rgba(255, 255, 255, 0.8);
}

.dropdown-item.selected {
  background: rgba(99, 102, 241, 0.1);
}

.dropdown-item.selected.highlighted,
.dropdown-item.selected:hover {
  background: var(--accent-color);
}

.dropdown-item.disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.dropdown-item.disabled:hover {
  background: transparent;
  color: var(--text-primary);
}

.item-content {
  display: flex;
  flex-direction: column;
  gap: 0.125rem;
  flex: 1;
  min-width: 0;
}

.item-label {
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.item-meta {
  font-size: 0.75rem;
  color: var(--text-muted);
}

.check-icon {
  flex-shrink: 0;
  color: var(--accent-color);
}

.dropdown-item.highlighted .check-icon,
.dropdown-item:hover .check-icon {
  color: white;
}

.no-results {
  padding: 1rem;
  text-align: center;
  color: var(--text-muted);
}

.create-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: var(--accent-color);
  font-weight: 500;
  border-top: 2px solid var(--border-color);
}

.create-item:hover {
  color: white;
}
</style>
