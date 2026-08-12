<script setup lang="ts">
import { computed } from 'vue'
import CardStack from './CardStack.vue'

const props = defineProps<{
  hand: number[]
}>()

const stacks = computed(() => {
  const sorted = [...props.hand].sort((a, b) => a - b)
  let pv = -1
  const output: [number, number][] = []
  for (const c of sorted) {
    if (c - 1 != pv) {
      output.push([c, 1])
    } else {
      output[output.length - 1][1] += 1
    }
    pv = c
  }
  return output
})
</script>

<template>
  <div class="is-flex is-gap-1">
    <CardStack v-for="s in stacks" :value="s[0]" :stack="s[1]" />
  </div>
</template>
