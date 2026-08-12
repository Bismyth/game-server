<script setup lang="ts">
import type { Options } from '@/game/skull/options'

defineProps<{
  data?: Options
}>()
</script>

<template>
  <div>Random Discard: {{ data?.discardRandom ? 'Yes' : 'No' }}</div>
</template>
