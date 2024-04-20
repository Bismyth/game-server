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
            <strong>Objective: </strong> To reveal roses around the table without hitting a skull
          </p>
          <p>
            <strong>Gameplay</strong>
          </p>

          <ul class="list">
            <li>Each player places a rose or a skull in front of themselves</li>
            <li>Play goes round and each player may place a rose or skull or start the bid</li>
            <li>
              When the bid starts you are bidding on how many roses you can flip over without
              hitting a skull
            </li>
            <li>Play continues but now you can only raise the bid or pass</li>
            <li>
              Once all other players have passed, the person who bid highest starts flipping over
              tiles.
            </li>
            <li>When flipping over tiles you start with all the tiles you have placed first.</li>
            <li>You stop flipping when either you have hit the bid or flipped a skull.</li>
            <li>If you hit a skull you lose a tile, if you didn't you get a point.</li>
            <li>First to two points win.</li>
          </ul>
        </div>
      </template>
    </ModalWrap>
  </Teleport>
</template>
