<script setup lang="ts">
import type { Options } from '@/game/liarsdice/options'

defineProps<{
  data?: Options
}>()
</script>

<template>
  <div>Starting Dice: {{ data?.startingDice }}</div>
</template>
