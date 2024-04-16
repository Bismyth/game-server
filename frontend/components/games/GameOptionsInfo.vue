<script setup lang="ts">
import { useRoomStore } from '@/stores/room'
import { defineAsyncComponent, watch } from 'vue'

const room = useRoomStore()

const getOptionsInfo = (gameType: typeof room.data.gameType) =>
  defineAsyncComponent({
    // the loader function
    loader: () => {
      return import(`./${gameType}/OptionsInfo.vue`)
    },
    timeout: 3000,
  })

let OptionsInfo = getOptionsInfo(room.data.gameType)

watch(
  () => room.data.gameType,
  (nv) => {
    OptionsInfo = getOptionsInfo(nv)
  },
)
</script>

<template>
  <div v-if="room.data.gameType">
    <OptionsInfo :data="room.data.options" />
  </div>
</template>
