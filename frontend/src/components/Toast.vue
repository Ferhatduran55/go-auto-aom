<template>
  <Teleport to="body">
    <Transition name="toast">
      <div v-if="toastVisible" :class="['toast', toastType]">
        <span v-if="toastType === 'error'">❌</span>
        <span v-else-if="toastType === 'success'">✅</span>
        <span v-else-if="toastType === 'warning'">⚠️</span>
        <span v-else>ℹ️</span>
        {{ toastMessage }}
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { useToast } from '@/composables/useToast'

const { toastMessage, toastType, toastVisible } = useToast()
</script>

<style scoped>
.toast {
  @apply fixed bottom-5 left-1/2 -translate-x-1/2 px-8 py-4 rounded-xl font-bold text-white z-[2000] shadow-xl;
}
.toast.success {
  @apply bg-success;
}
.toast.error {
  @apply bg-danger;
}
.toast.warning {
  @apply bg-warning;
}
.toast.info {
  @apply bg-info;
}
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}
.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateX(-50%) translateY(20px);
}
</style>
