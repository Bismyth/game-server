<script setup lang="ts">
import IconButton from '@/components/IconButton.vue'
import { ref } from 'vue'

const showModal = ref(false)

const handleOpen = () => {
  showModal.value = !showModal.value
}

const handleClose = () => {
  showModal.value = false
}
</script>

<template>
  <IconButton icon="fa6-solid:book" label="Rules" @click="handleOpen" />

  <Teleport to="#slideout-right">
    <div :class="{ 'is-active': showModal }">Test</div>
  </Teleport>
</template>
