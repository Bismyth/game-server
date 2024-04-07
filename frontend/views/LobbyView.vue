<script setup lang="ts">
import api from '@/api'
import LobbyOptions from '@/components/LobbyOptions.vue'
import GameOptionsForm from '@/components/games/GameOptionsForm.vue'
import { gameTypeLabels } from '@/game'
import { useLobbyStore } from '@/stores/lobby'
import FullLogo from '@/components/FullLogo.vue'
import DarkModeToggle from '@/components/DarkModeToggle.vue'
import LobbyUsers from '@/components/LobbyUsers.vue'
import ErrorStore from '@/components/ErrorStore.vue'
import IconButton from '@/components/IconButton.vue'
import { computed, ref } from 'vue'
import GameOptionsInfo from '@/components/games/GameOptionsInfo.vue'
import { useUserStore } from '@/stores/user'

const LINK_COPIED_TIMEOUT = 2 //seconds

const lobby = useLobbyStore()

const leave = () => {
  lobby.leave()
}

const create = () => {
  api.game.newGame(lobby.id)
}

const user = useUserStore()

const isHost = computed(() => lobby.users[user.data.id]?.host ?? false)

const showLinkCopiedText = ref(false)

const shareLink = () => {
  navigator.clipboard.writeText(`${window.location.origin}/?id=${lobby.id}`)
  showLinkCopiedText.value = true

  setTimeout(() => {
    showLinkCopiedText.value = false
  }, LINK_COPIED_TIMEOUT * 1000)
}
</script>

<template>
  <main class="centerize">
    <ErrorStore />
    <div class="box container b-primary" v-if="lobby.ready">
      <div class="logo-header mb-6">
        <div class="outer">
          <IconButton icon="fa6-solid:arrow-left" label="Leave" @click="leave" />
        </div>
        <FullLogo :width="300" />
        <div class="outer">
          <DarkModeToggle />
        </div>
      </div>
      <div class="body-wrapper">
        <div class="box body is-1 mb-0">
          <h1 class="title is-4">Players</h1>
          <LobbyUsers :users="lobby.users" />
        </div>
        <div class="box body is-5">
          <div class="is-flex mb-4">
            <h1 class="title is-4 mb-0">Lobby Info</h1>
            <div class="buttons ml-auto" v-if="isHost">
              <LobbyOptions />
            </div>
          </div>
          <div class="is-flex mb-5 is-flex-direction-column">
            <div class="is-flex mb-4">
              <div class="is-size-5">
                <span class="has-text-weight-semibold">Game Type: </span>
                {{ lobby.gameType === '' ? 'Not Selected' : gameTypeLabels[lobby.gameType] }}
              </div>
              <div class="buttons ml-auto" v-if="isHost">
                <GameOptionsForm />
              </div>
            </div>
            <div v-if="lobby.gameType !== ''">
              <h1 class="title is-5">Game Options</h1>
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
                @click="create"
                label="Start"
                icon="fa6-solid:play"
                color="primary"
                :disabled="lobby.gameType === ''"
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
.box.body.is-5 {
  flex-grow: 5;
  min-width: 400px;
}

.box.body.is-1 {
  flex-grow: 1;
  min-width: 300px;
}

.body-wrapper {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
}

.info-header {
  display: flex;
}
</style>
