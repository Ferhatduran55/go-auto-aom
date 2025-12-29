<template>
  <div class="critical-stock-badge" v-if="criticalCount > 0">
    <button 
      @click="dropdownOpen = !dropdownOpen" 
      class="badge-btn"
      :class="{ active: dropdownOpen }"
    >
      <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"></path>
        <path d="M12 9v4"></path>
        <path d="M12 17h.01"></path>
      </svg>
      <span class="badge-count">{{ criticalCount }}</span>
    </button>
    
    <!-- Dropdown -->
    <div v-if="dropdownOpen" class="dropdown-overlay" @click="dropdownOpen = false"></div>
    <div v-if="dropdownOpen" class="dropdown-menu">
      <div class="dropdown-header">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"></path>
          <path d="M12 9v4"></path>
          <path d="M12 17h.01"></path>
        </svg>
        Kritik Stok Uyarısı
      </div>
      
      <div class="dropdown-body">
        <div 
          v-for="product in criticalProducts" 
          :key="product.id" 
          class="critical-item"
          @click="productSelected(product)"
        >
          <div class="item-info">
            <span class="item-name">{{ product.name }}</span>
            <span class="item-category" v-if="product.category">{{ product.category }}</span>
          </div>
          <div class="item-stock">
            <span class="stock-value">{{ formatQuantity(product.stock_quantity, product.unit) }}</span>
            <span class="stock-unit">{{ product.unit || 'adet' }}</span>
          </div>
          <button @click.stop="openStockIn(product)" class="item-action" title="Stok Ekle">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 5v14"></path>
              <path d="M5 12h14"></path>
            </svg>
          </button>
        </div>
      </div>
      
      <div class="dropdown-footer">
        <button @click="showAll" class="btn btn-sm btn-secondary">
          Tüm Kritik Stokları Gör
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useStock } from '../composables/useStock'

const emit = defineEmits(['stock-in', 'product-selected', 'show-all'])

const { 
  criticalStockProducts: criticalProducts, 
  criticalCount, 
  loadCriticalStockProducts,
  formatQuantity 
} = useStock()

const dropdownOpen = ref(false)

// Periodic check interval
let checkInterval = null

const openStockIn = (product) => {
  dropdownOpen.value = false
  emit('stock-in', product)
}

const productSelected = (product) => {
  dropdownOpen.value = false
  emit('product-selected', product)
}

const showAll = () => {
  dropdownOpen.value = false
  emit('show-all')
}

// Close on ESC key
const handleKeydown = (e) => {
  if (e.key === 'Escape') {
    dropdownOpen.value = false
  }
}

onMounted(() => {
  // Initial load
  loadCriticalStockProducts()
  
  // Check every 30 seconds
  checkInterval = setInterval(() => {
    loadCriticalStockProducts()
  }, 30000)
  
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  if (checkInterval) {
    clearInterval(checkInterval)
  }
  document.removeEventListener('keydown', handleKeydown)
})
</script>

<style scoped>
.critical-stock-badge {
  position: relative;
}

.badge-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 0.75rem;
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 0.5rem;
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s;
}

[data-theme="dark"] .badge-btn {
  background: rgba(248, 113, 113, 0.15);
  border-color: rgba(248, 113, 113, 0.3);
  color: #f87171;
}

.badge-btn:hover, .badge-btn.active {
  background: rgba(239, 68, 68, 0.25);
  border-color: rgba(239, 68, 68, 0.5);
}

[data-theme="dark"] .badge-btn:hover, 
[data-theme="dark"] .badge-btn.active {
  background: rgba(248, 113, 113, 0.25);
  border-color: rgba(248, 113, 113, 0.5);
}

.badge-count {
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 1.25rem;
  height: 1.25rem;
  padding: 0 0.375rem;
  background: #ef4444;
  color: white;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 600;
}

[data-theme="dark"] .badge-count {
  background: #f87171;
}

.dropdown-overlay {
  position: fixed;
  inset: 0;
  z-index: 999;
}

.dropdown-menu {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 0.5rem;
  width: 320px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 0.75rem;
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.2), 0 8px 10px -6px rgba(0, 0, 0, 0.15);
  z-index: 1000;
  overflow: hidden;
}

.dropdown-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1rem;
  background: rgba(239, 68, 68, 0.1);
  border-bottom: 1px solid var(--border-color);
  font-weight: 600;
  font-size: 0.875rem;
  color: #ef4444;
}

[data-theme="dark"] .dropdown-header {
  background: rgba(248, 113, 113, 0.1);
  color: #f87171;
}

.dropdown-body {
  max-height: 300px;
  overflow-y: auto;
}

.critical-item {
  display: grid;
  grid-template-columns: 1fr auto auto;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1rem;
  border-bottom: 1px solid var(--border-color);
  cursor: pointer;
  transition: background 0.15s;
}

.critical-item:hover {
  background: var(--bg-secondary);
}

.critical-item:last-child {
  border-bottom: none;
}

.item-info {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.item-name {
  font-weight: 500;
  color: var(--text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.item-category {
  font-size: 0.75rem;
  color: var(--text-muted);
}

.item-stock {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.stock-value {
  font-weight: 600;
  color: #ef4444;
  font-size: 1rem;
}

[data-theme="dark"] .stock-value {
  color: #f87171;
}

.stock-unit {
  font-size: 0.7rem;
  color: var(--text-muted);
}

.item-action {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 1.75rem;
  height: 1.75rem;
  border: none;
  border-radius: 0.375rem;
  background: rgba(34, 197, 94, 0.15);
  color: #22c55e;
  cursor: pointer;
  transition: all 0.15s;
}

[data-theme="dark"] .item-action {
  background: rgba(74, 222, 128, 0.15);
  color: #4ade80;
}

.item-action:hover {
  background: rgba(34, 197, 94, 0.3);
}

[data-theme="dark"] .item-action:hover {
  background: rgba(74, 222, 128, 0.3);
}

.dropdown-footer {
  padding: 0.75rem 1rem;
  border-top: 1px solid var(--border-color);
  background: var(--bg-secondary);
}

.dropdown-footer button {
  width: 100%;
}
</style>
