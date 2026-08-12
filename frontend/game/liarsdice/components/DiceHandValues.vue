<script setup lang="ts">
import DiceCube from './DiceCube.vue'
import { Icon } from '@iconify/vue'
import { computed, onMounted, ref, watch } from 'vue'
import DiceHandTotals from './DiceHandTotals.vue'

const props = defineProps<{
  trueValue: number[]
  shuffle?: boolean
  highlight?: number
  tabIndex: number
}>()

const highlightOnes = computed(() => {
  if (props.highlight) {
    return 'success'
  }
  return undefined
})

const rawValues = computed(() => {
  const faces = [1, 2, 3, 4, 5, 6]
  return faces.map((i) => props.trueValue.filter((v) => v === i).length)
})

const shuffleDelay = 100

const diceHeight = computed(() => {
  if (props.trueValue.length < 8) {
    return 3
  }

  if (props.trueValue.length < 15) {
    return 2
  }

  return 1.5
})

const randomise = () => {
  const diceNum = props.trueValue.length
  randomValues.value = Array.from({ length: diceNum }, () => Math.floor(Math.random() * 6) + 1)

  if (props.shuffle) {
    setTimeout(randomise, shuffleDelay)
  }
}

const randomValues = ref<number[]>([])

onMounted(() => {
  if (props.shuffle) {
    randomise()
  }
})

watch(
  () => props.shuffle,
  (ns) => {
    if (ns) {
      randomise()
    }
  },
)
</script>

<template>
  <div v-show="tabIndex === 0">
    <div class="dice-hand">
      <DiceCube
        v-for="(number, i) in shuffle ? randomValues : trueValue"
        :key="i"
        :value="number"
        :height="diceHeight"
        :fill="number === highlight ? 'success' : number === 1 ? highlightOnes : undefined"
      />
    </div>
    <div v-if="shuffle">Rolling...</div>
  </div>
  <div v-show="tabIndex === 1">
    <div class="dice-hand">
      <div v-if="shuffle">Rolling...</div>
      <DiceHandTotals :highlight="highlight" :values="rawValues" />
    </div>
  </div>
</template>

<style>
.dice-hand {
  display: flex;
  gap: 10px;
  max-width: 800px;
  flex-wrap: wrap;
}
</style>
