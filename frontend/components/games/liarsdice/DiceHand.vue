<script setup lang="ts">
import { ref } from 'vue'

import DiceHandTabs from './DiceHandTabs.vue'
import DiceHandValues from './DiceHandValues.vue'

defineProps<{
  trueValue: number[]
  shuffle?: boolean
  highlight?: number
}>()

const tabIndex = ref(0)
</script>

<template>
  <DiceHandTabs @change-tab="(v) => (tabIndex = v)" />
  <DiceHandValues
    :true-value="trueValue"
    :shuffle="shuffle"
    :highlight="highlight"
    :tab-index="tabIndex"
  />
</template>
