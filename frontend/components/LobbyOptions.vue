<script setup lang="ts">
import { gameTypeLabels, type GameTypes } from '@/game'
import ModalWrap from './ModalWrap.vue'
import api from '@/api-utils/index.ts'
import { ref } from 'vue'
import IconButton from './IconButton.vue'
import { Icon } from '@iconify/vue'

const selectGame = (v: GameTypes) => {
  api.room.change({ gameType: v })
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
      <div class="boxes is-gap-1">
        <div class="box s-box" v-for="(value, key) in gameTypeLabels" :key="key">
          <span class="p-1" :style="{ fontSize: '64px' }">
            <Icon :icon="value.icon" />
          </span>
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
  width: 8rem;
}
.boxes {
  display: flex;
}
</style>
