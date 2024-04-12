<script setup lang="ts">
import game from '@/game/onenightwerewolf'
import { useLobbyStore } from '@/stores/lobby'
import { onMounted } from 'vue'

const lobby = useLobbyStore()

onMounted(() => {
  game.create(lobby.id)
})
</script>

<template>
  <h1>Game</h1>
  <pre>{{ game.data.public }}</pre>
  <pre>{{ game.data.private }}</pre>
</template>
