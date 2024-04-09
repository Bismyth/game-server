<script setup lang="ts">
import ModalWrap from '@/components/ModalWrap.vue'

defineProps<{
  show: boolean
}>()

const emit = defineEmits<{
  close: []
}>()

const handleClose = () => {
  emit('close')
}
</script>

<template>
  <ModalWrap title="Call Result" :shown="show" @close="handleClose">
    <template #body> This is the call screen </template>
    <template #footer>
      <div class="buttons">
        <button class="button" @click="handleClose">Continue</button>
      </div>
    </template>
  </ModalWrap>
</template>
