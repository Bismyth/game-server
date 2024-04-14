<script setup lang="ts">
defineProps<{
  shown: boolean
  title: string
}>()

const emit = defineEmits<{
  close: []
  submit: []
}>()

const close = (e: Event) => {
  e.preventDefault()
  emit('close')
}
</script>

<template>
  <div class="modal" :class="{ 'is-active': shown }">
    <div class="modal-background" @click="close"></div>
    <div class="modal-card">
      <header class="modal-card-head">
        <slot name="header">
          <p class="modal-card-title">{{ title }}</p>
          <button class="delete" aria-label="close" @click="close"></button>
        </slot>
      </header>
      <section class="modal-card-body">
        <slot name="body"></slot>
      </section>
      <footer class="modal-card-foot">
        <slot name="footer">
          <div class="buttons">
            <button class="button" @click="close">Close</button>
          </div>
        </slot>
      </footer>
    </div>
  </div>
</template>
