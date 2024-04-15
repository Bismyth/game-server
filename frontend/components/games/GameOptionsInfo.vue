<script setup lang="ts">
import { useRoomStore } from '@/stores/room'
import { defineAsyncComponent, watch } from 'vue'

const room = useRoomStore()

const getOptionsInfo = (gameType: typeof room.gameType) =>
  defineAsyncComponent({
    // the loader function
    loader: () => {
      return import(`./${gameType}/OptionsInfo.vue`)
    },
    timeout: 3000,
  })

let OptionsInfo = getOptionsInfo(room.gameType)

watch(
  () => room.gameType,
  (nv) => {
    OptionsInfo = getOptionsInfo(nv)
  },
)
</script>

<template>
  <div v-if="room.gameType">
    <OptionsInfo :data="room.options" />
  </div>
</template>
