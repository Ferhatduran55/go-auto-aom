<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-container">
      <!-- Header -->
      <div class="modal-header">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 bg-gradient-to-r from-accent to-purple-500 rounded-xl flex items-center justify-center text-white">
            📍
          </div>
          <h3 class="text-lg font-bold">Durum Güncelle</h3>
        </div>
        <button @click="$emit('close')" class="text-2xl hover:opacity-70 transition-opacity">✕</button>
      </div>

      <!-- Body -->
      <div class="modal-body">
        <div class="form-group mb-4">
          <label class="form-label">Yeni Durum</label>
          <select v-model="selectedStatus" class="form-input">
            <option value="">Durum Seçin...</option>
            <option 
              v-for="(info, key) in ORDER_STATUSES" 
              :key="key" 
              :value="key"
            >
              {{ info.emoji }} {{ info.label }}
            </option>
          </select>
        </div>

        <div class="form-group">
          <label class="form-label">Not (opsiyonel)</label>
          <textarea 
            v-model="statusNote"
            rows="3"
            placeholder="Durum hakkında not ekleyin..."
            class="form-input resize-none"
          ></textarea>
        </div>
      </div>

      <!-- Footer -->
      <div class="modal-footer">
        <button @click="$emit('close')" class="btn btn-secondary">
          İptal
        </button>
        <button 
          @click="handleSave" 
          :disabled="!selectedStatus || saving"
          class="btn btn-primary"
        >
          <span v-if="saving">Kaydediliyor...</span>
          <span v-else>💾 Kaydet</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useWhatsAppOrders, ORDER_STATUSES } from '@/composables/useWhatsAppOrders'

const props = defineProps({
  orderId: {
    type: String,
    required: true
  }
})

const emit = defineEmits(['close', 'saved'])

const { addStatus } = useWhatsAppOrders()

const selectedStatus = ref('')
const statusNote = ref('')
const saving = ref(false)

async function handleSave() {
  if (!selectedStatus.value || !props.orderId) return

  saving.value = true
  try {
    const result = await addStatus(props.orderId, selectedStatus.value, statusNote.value)
    if (result.success) {
      emit('saved')
    } else {
      console.error('Durum eklenemedi:', result.error)
    }
  } catch (e) {
    console.error('Hata:', e)
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(2px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-container {
  background: var(--bg-card);
  border-radius: 1rem;
  width: 100%;
  max-width: 480px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.25rem 1.5rem;
  border-bottom: 1px solid var(--border-color);
}

.modal-body {
  padding: 1.5rem;
  flex: 1;
  overflow-y: auto;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  padding: 1rem 1.5rem;
  border-top: 1px solid var(--border-color);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-label {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--text-muted);
}
</style>
