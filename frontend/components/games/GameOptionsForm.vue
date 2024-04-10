<script setup lang="ts">
import api from '@/api'
import { useLobbyStore } from '@/stores/lobby'
import { defineAsyncComponent, ref, watch } from 'vue'
import ModalWrap from '../ModalWrap.vue'
import { gameTypeLabels } from '@/game'
import IconButton from '../IconButton.vue'

const lobby = useLobbyStore()

const getGameOptions = (gameType: typeof lobby.gameType) =>
  defineAsyncComponent({
    // the loader function
    loader: () => {
      return import(`./${gameType}/OptionsForm.vue`)
    },
    timeout: 3000,
  })

let GameOptions = getGameOptions(lobby.gameType)

watch(
  () => lobby.gameType,
  (nv) => {
    GameOptions = getGameOptions(nv)
  },
)

// TODO: maybe type this?
let optionsFormRef = ref<any>(null)

const handleFormSubmit = (data: any) => {
  api.lobby.options(lobby.id, data)
  showGameOptions.value = false
}

const showGameOptions = ref(false)

const onClose = () => {
  showGameOptions.value = false
}

const onOpen = () => {
  optionsFormRef.value?.init()
  showGameOptions.value = true
}
</script>

<template>
  <div v-if="lobby.gameType">
    <IconButton @click="onOpen" icon="fa6-solid:screwdriver-wrench" label="Options" />
    <ModalWrap
      :shown="showGameOptions"
      :title="`${gameTypeLabels[lobby.gameType]} Options`"
      @close="onClose"
    >
      <template #body>
        <GameOptions ref="optionsFormRef" @submit="handleFormSubmit" />
      </template>
      <template #footer>
        <div class="buttons">
          <button class="button is-link" @click="optionsFormRef?.submit()">Submit</button>
          <button class="button" @click="onClose">Cancel</button>
        </div>
      </template>
    </ModalWrap>
  </div>
</template>
