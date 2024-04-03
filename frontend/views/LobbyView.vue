<script setup lang="ts">
import api from '@/api'
import { createGame } from '@/game/liarsdice'
import { useLobbyStore } from '@/stores/lobby'

const lobby = useLobbyStore()

const leave = () => {
  api.lobby.leave(lobby.id)
}

const create = () => {
  createGame(lobby.id, { startingDice: 5 })
}
</script>

<template>
  <div class="container">
    <h1 class="title">Lobby</h1>
    <button class="button" @click="leave">Leave</button>
    <button class="button" @click="create">Create</button>

    <h4 class="title is-size-4 my-4">Users</h4>
    <ul>
      <li v-for="(key, val) in lobby.users" :key="val">{{ key }}</li>
    </ul>
  </div>
</template>
