<script setup lang="ts">
import { useLobbyStore } from '@/stores/lobby'
import { defineAsyncComponent, watch } from 'vue'

const lobby = useLobbyStore()

const getOptionsInfo = (gameType: typeof lobby.gameType) =>
  defineAsyncComponent({
    // the loader function
    loader: () => {
      return import(`./${gameType}/OptionsInfo.vue`)
    },
    timeout: 3000,
  })

let OptionsInfo = getOptionsInfo(lobby.gameType)

watch(
  () => lobby.gameType,
  (nv) => {
    OptionsInfo = getOptionsInfo(nv)
  },
)
</script>

<template>
  <div v-if="lobby.gameType">
    <OptionsInfo :data="lobby.options" />
  </div>
</template>
