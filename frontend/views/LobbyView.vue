<script setup lang="ts">
import api from '@/api'
import LobbyOptions from '@/components/LobbyOptions.vue'
import GameOptionsForm from '@/components/games/GameOptionsForm.vue'
import { gameTypeLabels } from '@/game'
import { useRoomStore } from '@/stores/room'
import FullLogo from '@/components/FullLogo.vue'
import DarkModeToggle from '@/components/DarkModeToggle.vue'
import LobbyUsers from '@/components/LobbyUsers.vue'
import ErrorStore from '@/components/ErrorStore.vue'
import IconButton from '@/components/IconButton.vue'
import { ref } from 'vue'
import GameOptionsInfo from '@/components/games/GameOptionsInfo.vue'
import RulesPage from '@/components/games/RulesPage.vue'

const LINK_COPIED_TIMEOUT = 2 //seconds

const room = useRoomStore()

const leave = () => {
  room.leave()
}

const start = () => {
  api.game.startGame()
}

const showLinkCopiedText = ref(false)

const shareLink = () => {
  navigator.clipboard.writeText(`${window.location.origin}/?id=${room.id}`)
  showLinkCopiedText.value = true

  setTimeout(() => {
    showLinkCopiedText.value = false
  }, LINK_COPIED_TIMEOUT * 1000)
}
</script>

<template>
  <main class="centerize">
    <ErrorStore />
    <div class="box container b-primary" v-if="room.ready">
      <div class="logo-header mb-4">
        <div class="outer">
          <IconButton icon="fa6-solid:arrow-left" label="Leave" @click="leave" />
        </div>
        <div>
          <FullLogo :width="300" class="title-box" />
        </div>
        <div class="outer">
          <DarkModeToggle />
        </div>
      </div>
      <div class="body-wrapper">
        <div class="box is-1 mb-0">
          <h1 class="title is-4">Players</h1>
          <LobbyUsers :users="room.users" :host="room.host" />
        </div>
        <div class="box is-5">
          <div class="is-flex mb-4">
            <h1 class="title is-4 mb-0">Lobby Info</h1>
            <div class="buttons ml-auto" v-if="room.isHost">
              <LobbyOptions />
            </div>
          </div>
          <div class="is-flex mb-5 is-flex-direction-column">
            <div class="is-flex mb-2">
              <div class="is-size-5">
                <div class="mb-3">
                  <span class="has-text-weight-semibold">Game Type: </span>
                  {{ room.gameType === '' ? 'Not Selected' : gameTypeLabels[room.gameType] }}
                </div>
                <div v-if="room.gameType !== ''">
                  <RulesPage :game-type="room.gameType" />
                </div>
              </div>
            </div>
            <div v-if="room.gameType !== ''">
              <hr />
              <div class="is-flex mb-4">
                <h1 class="title is-4 mb-0">Game Options</h1>
                <div class="buttons ml-auto" v-if="room.isHost">
                  <GameOptionsForm />
                </div>
              </div>
              <GameOptionsInfo />
            </div>
          </div>
          <hr />
          <div class="field is-grouped">
            <div class="control">
              <IconButton @click="shareLink" label="Invite" icon="fa6-solid:link" />
              <p class="help" v-show="showLinkCopiedText">Link Copied!</p>
            </div>
            <div class="control">
              <IconButton
                @click="start"
                label="Start"
                icon="fa6-solid:play"
                color="primary"
                :disabled="room.gameType === ''"
                v-if="room.isHost"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
    <div v-else>Loading...</div>
  </main>
</template>

<style>
.info-header {
  display: flex;
}
</style>
