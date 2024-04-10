<script setup lang="ts">
import { computed } from 'vue'
import { Icon } from '@iconify/vue'

const props = defineProps<{ bid: string | [number, number] }>()

const numberValues = computed<[number, number]>(() => {
  if (typeof props.bid === 'string') {
    const numbers = props.bid.split(',').map((v) => parseInt(v))
    return [numbers[0], numbers[1]]
  } else {
    return props.bid
  }
})
</script>

<template>
  <span class="icon-text">
    <span>{{ numberValues[0] }}x</span>
    <span class="icon"><Icon :icon="`bi:dice-${numberValues[1]}-fill`" /></span>
  </span>
</template>
