<script setup lang="ts">
import GameError from '@/components/games/GameError.vue'
import GameLoading from '@/components/games/GameLoading.vue'
import { useLobbyStore } from '@/stores/lobby'
import { defineAsyncComponent } from 'vue'

const lobby = useLobbyStore()

const GameRender = defineAsyncComponent({
  // the loader function
  loader: () => {
    if (lobby.gameType === '') {
      throw Error('no game type set')
    }

    return import(`../components/games/${lobby.gameType}/GameDisplay.vue`)
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
</script>

<template>
  <div v-if="!lobby.inGame">Loading...</div>
  <GameRender v-else />
</template>
