<script setup lang="ts">
import type { Options } from '@/game/nothanks/options'

defineProps<{
  data?: Options
}>()
</script>

<template>
  <div>Total Rounds: {{ data?.rounds }}</div>
</template>
