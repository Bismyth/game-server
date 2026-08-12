<script setup lang="ts">
import EventLog from '@/components/event/EventLog.vue'
import GameError from '@/components/games/GameError.vue'
import GameLoading from '@/components/games/GameLoading.vue'
import { useRoomStore } from '@/stores/room'
import { defineAsyncComponent, onMounted } from 'vue'

const room = useRoomStore()

const GameRender = defineAsyncComponent({
  // the loader function
  loader: () => {
    if (room.data.gameType === '') {
      throw Error('no game type set')
    }

    return import(`@/game/${room.data.gameType}/components/GameDisplay.vue`)
  },

  // A component to use while the async component is loading
  loadingComponent: GameLoading,
  // Delay before showing the loading component. Default: 200ms.
  delay: 200,

  // A component to use if the load fails
  errorComponent: GameError,
  // The error component will be displayed if a timeout is
  // provided and exceeded. Default: Infinity.
  timeout: 3000,
})

onMounted(() => {
  room.setupConnection()
})
</script>

<template>
  <div v-if="!room.ready">Loading...</div>
  <div v-else>
    <GameRender />
    <EventLog />
  </div>
</template>
