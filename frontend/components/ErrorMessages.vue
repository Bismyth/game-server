<script setup lang="ts">
import type { ErrorMessage } from '@/stores/error'

defineProps<{ messages: Record<number, ErrorMessage> }>()

const emit = defineEmits<{
  delete: [number]
}>()

const deleteMessage = (e: MouseEvent) => {
  const target = e.target as HTMLButtonElement
  let idTarget = parseInt(target.getAttribute('data-id') ?? '')
  if (Number.isInteger(idTarget)) {
    emit('delete', idTarget)
  }
}
</script>

<template>
  <div class="container mb-3">
    <div
      v-for="(message, id) in messages"
      :key="id"
      class="notification is-light mt-3"
      :class="`is-${message.type}`"
    >
      <button class="delete" :data-id="id" @click="deleteMessage"></button>
      <p>
        {{ message.message }}
      </p>
    </div>
  </div>
</template>
