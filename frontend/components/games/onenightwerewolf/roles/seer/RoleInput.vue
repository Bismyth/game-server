<script setup lang="ts">
import { seerAction } from '@/game/onenightwerewolf'
import { ref } from 'vue'
import PlayerDropdown from '../../PlayerDropdown.vue'

const props = defineProps<{
  id: string
}>()

const player = ref('')

const steal = () => {
  seerAction(props.id, player.value, undefined)
}
</script>

<template>
  <div>
    <PlayerDropdown v-model="player" />
    <button class="button" @click="steal">Scry</button>
  </div>
</template>
