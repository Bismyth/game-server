<script setup lang="ts">
import type { Options } from '@/game/onenightwerewolf/options'
import { roleDisplay } from '@/game/onenightwerewolf/roles'

defineProps<{
  data?: Options
}>()
</script>

<template>
  <h1>Roles:</h1>
  <ul>
    <li v-for="(r, i) in data?.roles ?? []" :key="i">
      {{ roleDisplay[r].name }}
    </li>
  </ul>
</template>
