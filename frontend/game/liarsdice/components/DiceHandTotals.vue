<script setup lang="ts">
import { Icon } from '@iconify/vue'

defineProps<{ values: number[]; highlight?: number }>()
</script>

<template>
  <div class="is-flex is-flex-direction-column">
    <div class="icon-text mb-1">
      <span class="icon"><Icon icon="bi:dice-1-fill" height="100%" /></span>
      <span>- {{ values[0] }}</span>
    </div>
    <div class="icon-text mb-1" v-for="(a, f) in values.slice(1)" :key="f">
      <span class="icon" :class="{ 'has-text-success': f + 2 === highlight }"
        ><Icon :icon="`bi:dice-${f + 2}-fill`" height="100%"
      /></span>
      <span>- {{ a }} ({{ a + values[0] }})</span>
    </div>
  </div>
</template>
