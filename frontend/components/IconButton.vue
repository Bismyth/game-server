<script setup lang="ts">
import { Icon } from '@iconify/vue'
import type { IconifyIcon } from '@iconify/vue'

defineProps<{
  icon: string | IconifyIcon
  label: string
  disabled?: boolean
  color?: string
}>()

const emit = defineEmits<{
  click: [MouseEvent]
}>()

const click = (e: MouseEvent) => {
  emit('click', e)
}
</script>

<template>
  <button class="button" @click="click" :disabled="disabled" :class="color ? `is-${color}` : ''">
    <span class="icon">
      <Icon :icon="icon" />
    </span>
    <span>{{ label }}</span>
  </button>
</template>
