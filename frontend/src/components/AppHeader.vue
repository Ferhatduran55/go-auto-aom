<template>
  <header class="bg-gradient-to-r from-accent to-purple sticky top-0 z-50">
    <div class="flex justify-between items-center px-6 py-4">
      <div class="flex items-center gap-4">
        <div class="w-12 h-12 bg-white/20 rounded-xl flex items-center justify-center overflow-hidden">
          <img src="/favicon.ico" alt="Logo" class="w-10 h-10 object-contain" />
        </div>
        <div>
          <h1 class="text-xl font-extrabold text-white tracking-tight">AutoManagement</h1>
          <span class="text-sm text-white/80">Oto Yönetim Sistemi</span>
        </div>
      </div>
      <div class="flex gap-3 flex-wrap justify-end items-center">
        <!-- Kritik Stok Badge Slot -->
        <slot name="kritik-stok"></slot>
        
        <button @click="$emit('showCatalog')" class="header-btn">
          📚 Katalog
        </button>
        <button v-if="activeTab === 'orders'" @click="$emit('showHistory')" class="header-btn">
          📊 Müşteri Geçmişi
        </button>
        <button @click="$emit('newOrder')" class="header-btn header-btn-primary">
          🆕 Yeni Sipariş
        </button>
        <button @click="toggleTheme" class="header-btn">
          {{ settings.theme === 'dark' ? '🌙' : '☀️' }} Tema
        </button>
        <button @click="$emit('showSettings')" class="header-btn">
          ⚙️ Ayarlar
        </button>
      </div>
    </div>
  </header>
</template>

<script setup>
import { useSettings } from '../composables/useSettings'

defineProps({
  activeTab: {
    type: String,
    default: 'orders'
  }
})

defineEmits(['showHistory', 'newOrder', 'showCatalog', 'showSettings'])

const { settings, toggleTheme } = useSettings()
</script>

<style scoped>
.header-btn {
  @apply bg-white/15 border border-white/20 text-white px-4 py-2.5 rounded-lg cursor-pointer font-semibold text-sm flex items-center gap-2 transition-all hover:bg-white/25;
}

.header-btn-primary {
  @apply bg-white/30 border-white/40 hover:bg-white/40;
}
</style>
