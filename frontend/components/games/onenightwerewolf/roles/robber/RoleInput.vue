<script setup lang="ts">
import { robberAction } from '@/game/onenightwerewolf'
import { ref } from 'vue'
import PlayerDropdown from '../../PlayerDropdown.vue'

const props = defineProps<{
  id: string
}>()

const player = ref('')

const steal = () => {
  robberAction(props.id, player.value)
}
</script>

<template>
  <div>
    <PlayerDropdown v-model="player" />
    <button class="button" @click="steal">Steal</button>
  </div>
</template>
