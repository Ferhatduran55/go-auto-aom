<template>
  <header class="bg-gradient-to-r from-accent to-purple sticky top-0 z-50">
    <div class="flex justify-between items-center px-6 py-4">
      <div class="flex items-center gap-4">
        <div class="w-12 h-12 bg-white/20 rounded-xl flex items-center justify-center overflow-hidden">
          <img src="/favicon.ico" alt="Logo" class="w-10 h-10 object-contain" />
        </div>
        <div>
          <h1 class="text-xl font-extrabold text-white tracking-tight">OTO PARÇA</h1>
          <span class="text-sm text-white/80">Sipariş Yönetim Sistemi</span>
        </div>
      </div>
      <div class="flex gap-3 flex-wrap justify-end">
        <button @click="$emit('showCatalog')" class="header-btn">
          📚 Katalog
        </button>
        <button @click="$emit('showHistory')" class="header-btn">
          📊 Müşteri Geçmişi
        </button>
        <button @click="$emit('newOrder')" class="header-btn">
          🆕 Yeni Sipariş
        </button>
        <button @click="toggleTheme" class="header-btn">
          {{ isDark ? '🌙' : '☀️' }} Tema
        </button>
      </div>
    </div>
  </header>
</template>

<script setup>
import { ref, onMounted } from 'vue'

defineEmits(['showHistory', 'newOrder', 'showCatalog'])

const isDark = ref(true)

function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

onMounted(() => {
  const savedTheme = localStorage.getItem('theme') || 'dark'
  isDark.value = savedTheme === 'dark'
  document.documentElement.classList.toggle('dark', isDark.value)
})
</script>

<style scoped>
.header-btn {
  @apply bg-white/15 border border-white/20 text-white px-4 py-2.5 rounded-lg cursor-pointer font-semibold text-sm flex items-center gap-2 transition-all hover:bg-white/25;
}
</style>
