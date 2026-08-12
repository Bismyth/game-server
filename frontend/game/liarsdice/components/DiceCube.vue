<script setup lang="ts">
import { Icon } from '@iconify/vue'

withDefaults(
  defineProps<{
    value: number
    height?: number
    fill?: string
  }>(),
  {
    height: 3,
  },
)

const dice = [1, 2, 3, 4, 5, 6]
</script>

<template>
  <div :style="{ height: `${height}rem` }" :class="`has-text-${fill}`">
    <Icon
      v-for="i in dice"
      :key="i"
      :icon="`bi:dice-${i}-fill`"
      v-show="value === i"
      class="dice-icon"
    />
  </div>
</template>

<style>
.dice-icon {
  height: inherit;
}
</style>
