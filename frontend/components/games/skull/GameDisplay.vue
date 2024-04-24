<script setup lang="ts">
import ErrorStore from '@/components/ErrorStore.vue'
import skull from '@/game/skull'
import { useRoomStore } from '@/stores/room'
import { onMounted, watch } from 'vue'
import TileDisplay from './TileDisplay.vue'

const room = useRoomStore()

onMounted(async () => {
  skull.create()
})

watch(
  () => room.data.inGame,
  (v, ov) => {
    if (v && !ov) {
      skull.create()
    }
  },
)

const place = () => {
  skull.place(false)
}
</script>

<template>
  <ErrorStore />
  <TileDisplay show />

  <h1 class="title">Public</h1>
  <pre>{{ skull.gameData.publicState }}</pre>

  <h1 class="title">Private</h1>
  <pre>{{ skull.gameData.privateState }}</pre>

  <button class="button" @click="place">Place Tile</button>
</template>
