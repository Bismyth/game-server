<script setup lang="ts">
import type { GameTypes } from '@/game'
import { defineAsyncComponent, watch } from 'vue'

const props = defineProps<{
  gameType: GameTypes
}>()

const getRulesPage = (gameType: GameTypes) =>
  defineAsyncComponent({
    // the loader function
    loader: () => {
      return import(`./${gameType}/RulesPage.vue`)
    },
    timeout: 3000,
  })

let Rules = getRulesPage(props.gameType)

watch(
  () => props.gameType,
  (nv) => {
    Rules = getRulesPage(nv)
  },
)
</script>

<template>
  <Rules />
</template>
