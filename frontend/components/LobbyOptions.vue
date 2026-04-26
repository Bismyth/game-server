<script setup lang="ts">
import { roomDataSchema } from '@/api/room'
import { gameTypeLabels, type GameTypes } from '@/game'
import { useRoomStore } from '@/stores/room'
import FormWrap from './FormWrap.vue'
import ModalWrap from './ModalWrap.vue'
import type { z } from 'zod'
import api from '@/api'
import { ref, watch } from 'vue'
import type { ComponentExposed } from 'vue-component-type-helpers'
import IconButton from './IconButton.vue'
import { Icon } from '@iconify/vue'

const room = useRoomStore()

type FormResult = z.infer<typeof roomDataSchema>

const submit = (data: FormResult) => {
  api.room.change(data)
  showOptions.value = false
}


const selectGame = (v: GameTypes) => {
  api.room.change({gameType: v})
  showOptions.value = false
}


const showOptions = ref(false)
const onClose = () => {
  showOptions.value = false
}

const openModal = () => {
  showOptions.value = true
}

</script>

<template>
  <IconButton @click="openModal" icon="fa6-solid:pencil" label="Edit" />
  <ModalWrap :shown="showOptions" title="Lobby Options" @close="onClose">
    <template #body>
      <div class="boxes">
        <div class="box s-box" v-for="(value, key) in gameTypeLabels" :key="key">
          <Icon :icon="value.icon" />
          <span>{{ value.displayName }}</span>
          <button @click="selectGame(key)" class="button">Select</button>
        </div>
      </div>
    </template>
    <template #footer>
      <div class="buttons">
        <button class="button" @click="onClose">Cancel</button>
      </div>
    </template>
  </ModalWrap>
</template>

<style>
.s-box {
  display: flex;
  flex-direction: column;
  text-align: center;
  margin-bottom: 0 !important;
}
.boxes {
  display: flex;
}
</style>