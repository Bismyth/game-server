<script setup lang="ts">
import type { LobbyUsers } from '@/api/lobby'
import { Icon } from '@iconify/vue'

defineProps<{
  users: LobbyUsers
}>()

// const players = {
//   id:
// }
</script>

<template>
  <div>
    <div class="player-box" v-for="(user, id) in users" :key="id">
      <span class="icon-text">
        <span class="icon">
          <Icon icon="fa6-solid:crown" v-if="user.host" />
        </span>
        <span>{{ user.name }}</span>
      </span>
    </div>
  </div>
</template>

<style>
.player-box {
}
</style>
