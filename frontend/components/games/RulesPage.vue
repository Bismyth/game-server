<script setup lang="ts">
import type { GameTypes } from '@/game'
import { defineAsyncComponent, watch, ref } from 'vue'
import IconButton from '../IconButton.vue'

const props = defineProps<{
  gameType: GameTypes
}>()

const getRulesPage = (gameType: GameTypes) =>
  defineAsyncComponent({
    // the loader function
    loader: () => {
      return import(`./${gameType}/RulesPage.vue`)
    },
    timeout: 3000,
  })

let Rules = getRulesPage(props.gameType)

watch(
  () => props.gameType,
  (nv) => {
    Rules = getRulesPage(nv)
  },
)

const showModal = ref(false)

const handleOpen = () => {
  showModal.value = true
}

const handleClose = () => {
  showModal.value = false
}
</script>

<template>
  <IconButton icon="fa6-solid:book" label="Rules" @click="handleOpen" />

  <Teleport to="body">
    <ModalWrap :shown="showModal" title="Rules" @close="handleClose">
      <template #body>
        <Rules />
      </template>
    </ModalWrap>
  </Teleport>
</template>
