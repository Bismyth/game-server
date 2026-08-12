<script setup lang="ts">
import ModalWrap from '@/components/ModalWrap.vue'
import RoomName from '@/components/RoomName.vue'
import Hand from './Hand.vue'
import { useRoomStore } from '@/stores/room.ts'

defineProps<{
  show: boolean
  hands?: Record<string, number[]>
}>()

const emit = defineEmits<{
  close: []
}>()

const handleClose = () => {
  emit('close')
}

const room = useRoomStore()
</script>

<template>
  <ModalWrap title="Current Hands" :shown="show" @close="handleClose">
    <template #body>
      <div v-if="hands">
        <div v-for="(h, id) in hands" class="mt-3">
          <div class="mb-2"><RoomName :id="id" /></div>
          <Hand :hand="h" />
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
