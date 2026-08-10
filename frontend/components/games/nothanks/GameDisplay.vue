<script setup lang="ts">
import { useRoomStore } from '@/stores/room'
import { onMounted, watch } from 'vue'
import nothanks from '@/game/no-thanks'
import ErrorStore from '@/components/ErrorStore.vue'

const room = useRoomStore()

onMounted(async () => {
  nothanks.create()
})

watch(
  () => room.data.inGame,
  (v, ov) => {
    if (v && !ov) {
      nothanks.create()
    }
  },
)
</script>

<template>
  <div class="p-4">
    <ErrorStore />
    <span v-if="nothanks.gameData.publicState?.currentPlayer == room.data.userId">IS TURN!</span>
    <p>Public</p>
    <pre class="mb-4">{{ nothanks.gameData.publicState }}</pre>
    <p>Private</p>
    <pre class="mb-4">{{ nothanks.gameData.privateState }}</pre>

    <button @click="nothanks.pass()" class="button is-primary mr-2">Pass</button>
    <button @click="nothanks.take()" class="button is-primary">Take</button>
  </div>
</template>
