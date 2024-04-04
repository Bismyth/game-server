<script setup lang="ts">
import api from '@/api'
import LobbyOptions from '@/components/LobbyOptions.vue'
import { gameTypeLabels } from '@/game'
import { useLobbyStore } from '@/stores/lobby'
import { defineAsyncComponent, watch } from 'vue'
import NoOptions from '@/components/games/NoOptions.vue'

const lobby = useLobbyStore()

const leave = () => {
  lobby.leave()
}

const create = () => {
  api.game.newGame(lobby.id)
}

const shareLink = () => {
  navigator.clipboard.writeText(`${window.location.origin}/lobby/join/${lobby.id}`)
}

const getGameOptions = (gameType: typeof lobby.gameType) => {
  return defineAsyncComponent({
    // the loader function
    loader: () => {
      if (gameType === '') {
        return import('../components/games/NoOptions.vue')
      }

      return import(`../components/games/${gameType}/OptionsForm.vue`)
    },
    timeout: 3000,
  })
}

let OptionsForm = getGameOptions(lobby.gameType)

watch(
  () => lobby.gameType,
  (newGameType) => {
    OptionsForm = getGameOptions(newGameType)
  },
)
</script>

<template>
  <div class="container">
    <h1 class="title">Lobby</h1>
    <button class="button" @click="leave">Leave</button>
    <button class="button" @click="create">Create</button>
    <button class="button" @click="shareLink">Share Link</button>
    <LobbyOptions />

    <h1>Game Type: {{ lobby.gameType === '' ? 'Not Set' : gameTypeLabels[lobby.gameType] }}</h1>

    <h4 class="title is-size-4 my-4">Game Options</h4>
    <OptionsForm />
    <h4 class="title is-size-4 my-4">Users</h4>
    <ul>
      <li v-for="(key, val) in lobby.users" :key="val">{{ key }}</li>
    </ul>
  </div>
</template>
