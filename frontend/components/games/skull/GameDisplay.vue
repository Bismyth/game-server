<script setup lang="ts">
import ErrorStore from '@/components/ErrorStore.vue'
import skull from '@/game/skull'
import { useRoomStore } from '@/stores/room'
import { onMounted, ref, watch } from 'vue'
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
  skull.place(tilePlace.value)
}

const bid = ref("");

const tilePlace = ref(false)


const makeBid = () => {
  skull.bid(parseInt(bid.value))
}


const pass = () => {
  skull.pass()
}

const selectedPlayer = ref('')


const flip = () => {
  skull.flip(selectedPlayer.value)
}

</script>

<template>
  <ErrorStore />
  <TileDisplay show />

  <h1 class="title">Public</h1>
  <pre>{{ skull.gameData.publicState }}</pre>

  <h1 class="title">Private</h1>
  <pre>{{ skull.gameData.privateState }}</pre>
  

  <h1 v-if="skull.gameData.isTurn">YOUR TURN</h1>
  
  <select v-model="tilePlace">
    <option :value="true">True</option>
    <option :value="false">False</option>
  </select>
  <button class="button" @click="place">Place Tile</button>


  <input v-model="bid">
  <button class="button" @click="makeBid">Bid</button>

  <button class="button" @click="pass">Pass</button>




  <select v-model="selectedPlayer">
    <option v-for="id in skull.gameData.publicState?.turnOrder" :key="id" :value="id">{{ room.users.names[id] }}</option>
  </select>
  <button class="button" @click="flip">Flip</button>


</template>
