<script setup lang="ts">
import ModalWrap from '@/components/ModalWrap.vue'
import RoomName from '@/components/RoomName.vue'
import type { PreviousRound, publicStateT } from '@/game/nothanks/index.ts'
import Hand from './Hand.vue'
import { computed } from 'vue'

const props = defineProps<{
  show: boolean
  previousRound?: PreviousRound
  publicState?: publicStateT
}>()

const emit = defineEmits<{
  close: []
}>()

const handleClose = () => {
  emit('close')
}
const title = computed(() => {
  if (props.publicState?.gameOver) {
    return 'Game Over'
  }
  return 'Previous Round'
})
</script>

<template>
  <ModalWrap :title="title" :shown="show" @close="handleClose">
    <template #body>
      <div v-if="previousRound">
        <p>Round: {{ previousRound.round }}</p>

        <div v-for="(s, p) in previousRound.score" class="mt-3">
          <div class="mb-2">
            <RoomName :id="p" /> - {{ s }}
            <span v-if="publicState?.gameOver">Total: {{ publicState.score[p] }}</span>
          </div>
          <Hand :hand="previousRound.playerCards[p]" />
        </div>
      </div>
      <div v-else>No Data...</div>
    </template>
    <template #footer>
      <div class="buttons">
        <button class="button" @click="handleClose">Close</button>
      </div>
    </template>
  </ModalWrap>
</template>
