<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import DiceCube from './DiceCube.vue'
import { Icon } from '@iconify/vue'

const props = defineProps<{
  trueValue: number[]
  shuffle?: boolean
}>()

const tabs = ['Dice', 'Value']
const tabIndex = ref(0)

const handleTabChange = (e: Event, i: number) => {
  e.preventDefault()
  tabIndex.value = i
}

const isIndex = (i: number) => {
  return i === tabIndex.value
}

const rawValues = computed(() => {
  const faces = [1, 2, 3, 4, 5, 6]
  return faces.map((i) => props.trueValue.filter((v) => v === i).length)
})

const shuffleDelay = 100

const diceHeight = computed(() => {
  if (props.trueValue.length < 5) {
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
  <div class="tabs mb-1">
    <ul>
      <li v-for="(t, i) in tabs" :key="i" :class="{ 'is-active': i === tabIndex }">
        <a @click="(e) => handleTabChange(e, i)">{{ t }}</a>
      </li>
    </ul>
  </div>
  <div v-show="tabIndex === 0">
    <div class="dice-hand">
      <div v-if="shuffle">Rolling...</div>
      <DiceCube
        v-for="(number, i) in shuffle ? randomValues : trueValue"
        :key="i"
        :value="number"
        :height="diceHeight"
      />
    </div>
  </div>
  <div v-show="isIndex(1)">
    <div class="dice-hand">
      <div v-if="shuffle">Rolling...</div>
      <div v-else class="is-flex is-flex-direction-column">
        <div class="icon-text mb-1">
          <span class="icon"><Icon icon="bi:dice-1-fill" height="100%" /></span>
          <span>- {{ rawValues[0] }}</span>
        </div>
        <div class="icon-text mb-1" v-for="(a, f) in rawValues.slice(1)" :key="f">
          <span class="icon"><Icon :icon="`bi:dice-${f + 2}-fill`" height="100%" /></span>
          <span>- {{ a }} ({{ a + rawValues[0] }})</span>
        </div>
      </div>
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
