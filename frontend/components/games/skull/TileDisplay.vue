<script setup lang="ts">
import { Icon } from '@iconify/vue'

defineProps<{
  skull?: boolean
  show: boolean
}>()
</script>

<template>
  <div class="game-tile">
    <div v-show="show" class="tile-icon">
      <Icon :icon="`mdi:${skull ? 'skull' : 'flower-poppy'}`" height="100%" width="100%" />
    </div>
  </div>
</template>

<style>
.game-tile {
  border: purple 4px solid;
  width: 4rem;
  height: 4rem;
  border-radius: 2rem;
  display: flex;
  justify-content: center;
  align-items: center;
}

.tile-icon {
  width: 100%;
  height: 100%;
  padding: 0.5rem;
}
</style>
