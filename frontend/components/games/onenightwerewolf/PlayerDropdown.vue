<script setup lang="ts">
import { useLobbyStore } from '@/stores/lobby'

const lobby = useLobbyStore()

const player = defineModel()
</script>

<template>
  <div class="select">
    <select v-model="player">
      <option v-for="(u, id) in lobby.users" :key="id">{{ u.name }}</option>
    </select>
  </div>
</template>
