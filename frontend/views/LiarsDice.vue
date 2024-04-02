<script setup lang="ts">
import ld from '@/game/liarsdice'
import { useSocketStore } from '@/stores/socket'
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const gameId = computed(() => route.params.id.toString())

onMounted(async () => {
  const socket = useSocketStore()
  await socket.isActive

  ld.create(gameId.value)
})

const currentBid = ref('')

const bid = () => {
  ld.bid(currentBid.value)
}

const call = () => {
  ld.call()
}
</script>

<template>
  <div class="container">
    <h1 class="title">Liars Dice</h1>
    <h2 class="title is-size-4 my-4">Public</h2>
    <pre>{{ ld.gameData.publicState }}</pre>

    <h2 class="title is-size-4 my-4">Private</h2>
    <pre>{{ ld.gameData.privateState }}</pre>

    <div v-show="ld.gameData.isTurn">
      <input v-model="currentBid" class="input" />
      <button class="button" @click="bid">Bid</button>

      <button class="button" @click="call">Call</button>
    </div>
  </div>
</template>
