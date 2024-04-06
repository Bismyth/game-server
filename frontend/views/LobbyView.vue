<script setup lang="ts">
import api from '@/api'
import LobbyOptions from '@/components/LobbyOptions.vue'
import LobbyGameOptions from '@/components/games/LobbyGameOptions.vue'
import { gameTypeLabels } from '@/game'
import { useLobbyStore } from '@/stores/lobby'

const lobby = useLobbyStore()

const leave = () => {
  lobby.leave()
}

const create = () => {
  api.game.newGame(lobby.id)
}

const shareLink = () => {
  navigator.clipboard.writeText(`${window.location.origin}/?id=${lobby.id}`)
}
</script>

<template>
  <div class="container" v-if="lobby.ready">
    <h1 class="title">Lobby</h1>
    <button class="button" @click="leave">Leave</button>
    <button class="button" @click="create">Create</button>
    <button class="button" @click="shareLink">Share Link</button>
    <LobbyOptions />
    <LobbyGameOptions />
    <h1>Game Type: {{ lobby.gameType === '' ? 'Not Set' : gameTypeLabels[lobby.gameType] }}</h1>

    <h4 class="title is-size-4 my-4">Users</h4>
    <ul>
      <li v-for="(key, val) in lobby.users" :key="val">{{ key }}</li>
    </ul>
  </div>
  <div v-else>Loading...</div>
</template>
