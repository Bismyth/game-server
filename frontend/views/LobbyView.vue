<script setup lang="ts">
import api from '@/api'
import { useSocketStore } from '@/stores/socket'
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const lobbyId = computed(() => route.params.id.toString())

const users = ref<string[]>([])

const handleLobbyUserChange = (iUsers: string[]) => {
  users.value = iUsers
}

const leave = () => {
  api.lobby.leave(lobbyId.value)
  router.push({ name: 'home' })
}

onMounted(async () => {
  const socket = useSocketStore()
  await socket.isActive

  api.lobby.setOnLobbyChange(handleLobbyUserChange)
  api.lobby.users(lobbyId.value)
})

onUnmounted(() => {
  api.lobby.clearOnLobbyChange()
})
</script>

<template>
  <div class="container">
    <h1 class="title">Lobby</h1>
    <button class="button" @click="leave">Leave</button>

    <h4 class="title is-size-4 my-4">Users</h4>
    <ul>
      <li v-for="(user, i) in users" :key="i">{{ user }}</li>
    </ul>
  </div>
</template>
