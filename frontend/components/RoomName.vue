<script setup lang="ts">
import api from '@/api'
import { useRoomStore } from '@/stores/room'
import { Icon } from '@iconify/vue'
import { computed } from 'vue'

const props = defineProps<{
  id: string
  kick?: boolean
}>()

const room = useRoomStore()

const kickAction = () => {
  api.room.kick(props.id)
}

const showKick = computed(() => {
  return room.isHost && props.id !== room.data.userId
})
</script>

<template>
  <span class="icon-text" v-if="kick">
    <span v-if="room.ready" :class="{ 'has-text-weight-bold': id === room.data.userId }">
      {{ room.users.names[id] }}
    </span>

    <a class="icon btn" v-tooltip.right="'Kick'" @click="kickAction" v-if="showKick">
      <Icon icon="fa6-solid:arrow-right-from-bracket" />
    </a>
  </span>
  <span v-else>
    <span v-if="room.ready" :class="{ 'has-text-weight-bold': id === room.data.userId }">
      {{ room.users.names[id] }}
    </span>
  </span>
</template>
