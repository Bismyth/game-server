<script setup lang="ts">
import { useErrorStore } from '@/stores/error'
import ErrorMessages from './ErrorMessages.vue'

const errorState = useErrorStore()

const messages = errorState.messageGetter

const onDelete = (n: number) => {
  errorState.deleteMessage(n)
}
</script>

<template>
  <ErrorMessages :messages="messages" @delete="onDelete" />
</template>
