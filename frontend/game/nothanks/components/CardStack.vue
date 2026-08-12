<script setup lang="ts">
import { computed, ref } from 'vue'
import SingleCard from './SingleCard.vue'

const props = defineProps<{
  value: number
  stack: number
}>()

const boxShadow = computed(() => {
  let count = Math.min(props.stack, 3)
  let shadows = []
  for (let i = 1; i < count; i++) {
    const d = 2 * i
    shadows.push(`${d}px ${d}px 0 -1px var(--bulma-background)`)
    shadows.push(`${d}px ${d}px 0 0 var(--bulma-primary)`)
  }
  return shadows.join(', ')
})

const expanded = ref(false)

const fullList = computed(() => {
  const output: number[] = []
  for (let i = 0; i < props.stack; i++) {
    output.push(props.value + i)
  }
  return output
})

const open = () => {
  if (props.stack > 1) {
    expanded.value = true
  }
}
</script>

<template>
  <SingleCard
    :value="value"
    :style="{ boxShadow }"
    v-if="!expanded"
    @click="open"
    :clickable="stack > 1"
  />
  <div v-else @click="expanded = false" class="cardList">
    <SingleCard v-for="card in fullList" :value="card" clickable />
  </div>
</template>

<style scoped>
.cardList {
  display: flex;
}
</style>
