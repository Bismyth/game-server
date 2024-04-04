<script setup lang="ts">
import api from '@/api'

const joinLobby = () => {
  const id = prompt('Enter Lobby ID:')
  if (id) {
    api.lobby.join(id)
  }
}

const createLobby = () => {
  api.lobby.create()
}
</script>

<template>
  <main class="container">
    <h1 class="title mb-2">Home</h1>

    <div class="mb-2"><button class="button" @click="createLobby">Create Lobby</button></div>
    <div><button class="button" @click="joinLobby">Join Lobby</button></div>
  </main>
</template>
