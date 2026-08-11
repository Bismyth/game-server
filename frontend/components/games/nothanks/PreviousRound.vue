<script setup lang="ts">
import ModalWrap from '@/components/ModalWrap.vue'
import RoomName from '@/components/RoomName.vue'
import type { PreviousRound } from '@/game/no-thanks'
import Hand from './Hand.vue'

defineProps<{
  show: boolean
  previousRound?: PreviousRound
}>()

const emit = defineEmits<{
  close: []
}>()

const handleClose = () => {
  emit('close')
}
</script>

<template>
  <ModalWrap title="Previous Round" :shown="show" @close="handleClose">
    <template #body>
      <div v-if="previousRound">
        <p>Round: {{ previousRound.round }}</p>

        <div v-for="(s, p) in previousRound.score" class="mt-3">
          <span class="mb-2"><RoomName :id="p" /> - {{ s }}</span>
          <Hand :hand="previousRound.playerCards[p]" />
        </div>
        <!-- <pre>{{ previousRound }}</pre> -->
      </div>
      <div v-else>No Data...</div>

      <!-- <div v-for=""></div> -->
    </template>
    <template #footer>
      <div class="buttons">
        <button class="button" @click="handleClose">Close</button>
      </div>
    </template>
  </ModalWrap>
</template>
