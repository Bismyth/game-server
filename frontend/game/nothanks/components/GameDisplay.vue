<script setup lang="ts">
import { useRoomStore } from '@/stores/room'
import { computed, onMounted, ref, watch } from 'vue'
import nt from '@/game/nothanks/index.ts'
import ErrorStore from '@/components/ErrorStore.vue'
import { handleLobbyBack } from '@/game'
import IconButton from '@/components/IconButton.vue'
import { Icon } from '@iconify/vue'
import RoomName from '@/components/RoomName.vue'
import SingleCard from './SingleCard.vue'
import Hand from './Hand.vue'
import PreviousRound from './PreviousRound.vue'
import AllHands from './AllHands.vue'
import RulesPage from '@/components/games/RulesPage.vue'

const room = useRoomStore()

onMounted(async () => {
  nt.create()
})

watch(
  () => room.data.inGame,
  (v, ov) => {
    if (v && !ov) {
      nt.create()
    }
  },
)
const leave = () => {
  room.leave()
}

const handlePrClose = () => {
  nt.showPr.value = false
}

const showAllHands = ref(false)
const handleAHClose = () => {
  showAllHands.value = false
}

const outPlayers = computed(() => {
  const ids: string[] = []
  for (const k of room.users.order) {
    if (!nt.gameData.publicState?.turnOrder.includes(k)) {
      ids.push(k)
    }
  }
  return ids
})
</script>

<template>
  <main class="centerize">
    <ErrorStore />
    <div class="box b-primary container">
      <div class="logo-header mb-6">
        <div class="outer">
          <IconButton
            icon="fa6-solid:arrow-left"
            label="Back to Lobby"
            @click="handleLobbyBack"
            v-if="nt.gameData.publicState?.gameOver"
          />
          <IconButton icon="fa6-solid:arrow-left" label="Leave" @click="leave" v-else />
        </div>
        <div class="title-box">
          <div class="is-flex">
            <h1 class="title mr-3">No Thanks!</h1>
            <div>
              <RulesPage game-type="nothanks" />
            </div>
          </div>
        </div>
        <div class="outer">
          <IconButton
            icon="fa6-solid:clock-rotate-left"
            label="Previous Round"
            v-if="nt.previousRound.value !== undefined"
            @click="nt.showPr.value = true"
          />
        </div>
      </div>
      <div class="body-wrapper">
        <div class="box mb-0 is-1">
          <h4 class="title is-5">Players</h4>
          <div class="mb-4">
            <table class="table is-fullwidth">
              <thead>
                <tr>
                  <th :style="{ width: '2px' }"></th>
                  <th>Name</th>
                  <th :style="{ width: '40px' }">Pts</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="id in nt.gameData.publicState?.turnOrder" :key="id">
                  <td>
                    <span class="icon">
                      <Icon
                        icon="fa6-solid:arrow-right"
                        v-if="id == nt.gameData.publicState?.currentPlayer"
                      />
                    </span>
                  </td>
                  <td><RoomName :id="id" kick /></td>
                  <td>{{ nt.gameData.publicState?.score[id] }}</td>
                </tr>
              </tbody>
            </table>
          </div>
          <h4 class="title is-5">Out</h4>
          <div>
            <div v-for="id in outPlayers" :key="id">
              <RoomName :id="id" kick />
            </div>
          </div>
        </div>
        <div class="box is-5">
          <div v-if="nt.gameData.publicState">
            <p class="mb-1 is-size-4">
              <strong>Round:</strong> {{ nt.gameData.publicState.currentRound }}/{{
                nt.gameData.publicState?.totalRounds
              }}
            </p>
            <p class="mb-2 is-size-5">Cards Remaining: {{ nt.gameData.publicState.deckLeft }}</p>
            <div class="is-flex is-align-items-center is-gap-2">
              <SingleCard :value="nt.gameData.publicState.inPlayCard" />
              <span>Tokens: {{ nt.gameData.publicState.tokensOnCard }}</span>
            </div>

            <div v-if="nt.gameData.privateState" class="mb-2">
              <p>Tokens Left: {{ nt.gameData.privateState.tokens }}</p>
            </div>
            <div v-if="room.data.userId == nt.gameData.publicState.currentPlayer">
              <p class="mb-1">It's your turn:</p>
              <div class="is-flex is-gap-2 mb-4">
                <button class="button is-primary" @click="nt.pass">Pass</button>
                <button class="button is-primary" @click="nt.take">Take</button>
              </div>
            </div>

            <div class="mb-4">
              <p>Current Hand:</p>
              <Hand :hand="nt.gameData.publicState.playerCards[room.data.userId] ?? []" />
            </div>
            <div>
              <button class="button" @click="showAllHands = true">Show Other Hands</button>
            </div>
          </div>
          <div v-else><span>Loading...</span></div>
        </div>
      </div>
    </div>
  </main>
  <PreviousRound
    :show="nt.showPr.value"
    :previous-round="nt.previousRound.value"
    :public-state="nt.gameData.publicState"
    @close="handlePrClose"
  />
  <AllHands
    :show="showAllHands"
    :hands="nt.gameData.publicState?.playerCards"
    @close="handleAHClose"
  />
</template>
