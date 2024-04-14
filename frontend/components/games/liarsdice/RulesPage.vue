<script setup lang="ts">
import IconButton from '@/components/IconButton.vue'
import ModalWrap from '@/components/ModalWrap.vue'
import { ref } from 'vue'

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
        <div class="content">
          <p>
            <strong>Objective: </strong> To outbid opponents by guessing the quantity of dice
            showing a specific face value, while bluffing about your own dice.
          </p>
          <p>
            <strong>Gameplay</strong>
          </p>

          <ul class="list">
            <li>
              Players take turns bidding on the total number of dice showing a specific face value
              among all players' dice.
            </li>
            <li>Bids must increase the quantity or face value, or both.</li>
            <li>A player can call out the previous bid if they believe it is incorrect.</li>
            <li>
              When called out, all dice are revealed, and if the bid was correct, the challenger
              loses a die. If the bid was incorrect, the bidder loses a die.
            </li>
          </ul>
        </div>
      </template>
    </ModalWrap>
  </Teleport>
</template>
