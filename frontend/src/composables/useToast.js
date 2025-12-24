import { ref } from 'vue'

const toastMessage = ref('')
const toastType = ref('success')
const toastVisible = ref(false)

let toastTimeout = null

function showToast(message, type = 'success') {
  toastMessage.value = message
  toastType.value = type
  toastVisible.value = true
  
  if (toastTimeout) clearTimeout(toastTimeout)
  toastTimeout = setTimeout(() => {
    toastVisible.value = false
  }, 3000)
}

export function useToast() {
  return {
    toastMessage,
    toastType,
    toastVisible,
    showToast,
  }
}
