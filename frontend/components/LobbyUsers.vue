<script setup lang="ts">
import type { RoomUsers } from '@/api/room'
import { Icon } from '@iconify/vue'
import RoomName from './RoomName.vue'

defineProps<{
  users: RoomUsers
  host: string
}>()
</script>

<template>
  <div>
    <div class="player-box" v-for="id in users.order" :key="id">
      <span class="icon-text">
        <span class="icon">
          <Icon icon="fa6-solid:crown" v-if="host === id" />
        </span>
        <RoomName :id="id" />
      </span>
    </div>
  </div>
</template>
