<script setup lang="ts">
import ModalWrap from '@/components/ModalWrap.vue'

defineProps<{
  show: boolean
  lastPlayer: string
}>()

const emit = defineEmits<{
  close: []
  back: []
}>()

const handleClose = () => {
  emit('close')
}

const handleBack = () => {
  emit('back')
}
</script>

<template>
  <ModalWrap title="Game Over" :shown="show" @close="handleClose">
    <template #body> {{ lastPlayer }} has won! </template>
    <template #footer>
      <div class="buttons">
        <button class="button" @click="handleClose">Show Board</button>
        <button class="button is-primary" @click="handleBack">Back to Lobby</button>
      </div>
    </template>
  </ModalWrap>
</template>
