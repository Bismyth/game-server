<script setup lang="ts">
import ld from '@/game/liarsdice'
import { computed, onMounted, watch } from 'vue'
import IconButton from '@/components/IconButton.vue'
import DiceHand from './DiceHand.vue'
import DiceCube from './DiceCube.vue'
import BidForm from './BidForm.vue'
import ErrorStore from '@/components/ErrorStore.vue'
import { useRoomStore } from '@/stores/room'
import { Icon } from '@iconify/vue'
import CallModal from './CallModal.vue'
import EndGame from './EndGame.vue'
import { useRouter } from 'vue-router'
import RulesPage from './RulesPage.vue'
import RoomName from '@/components/RoomName.vue'

const room = useRoomStore()

const router = useRouter()

onMounted(async () => {
  ld.create()
})

watch(
  () => room.data.inGame,
  (v, ov) => {
    if (v && !ov) {
      ld.create()
    }
  },
)

const outPlayers = computed(() => {
  const ids: string[] = []
  for (const k of room.users.order) {
    if (!ld.gameData.publicState?.turnOrder.includes(k)) {
      ids.push(k)
    }
  }
  return ids
})

const isIn = (id: string) => {
  return ld.gameData.publicState?.turnOrder.includes(id)
}

const handleBid = (bid: string) => {
  ld.bid(bid)
}

const handleCall = () => {
  ld.call()
}

const leave = () => {
  room.leave()
}

const totalDice = computed(() => {
  if (ld.gameData.publicState?.diceAmounts) {
    return Object.values(ld.gameData.publicState.diceAmounts).reduce((a, b) => a + b)
  }
  return 0
})

const showCall = () => {
  ld.showCall.value = true
}

const handleCallScreenClose = () => {
  ld.showCall.value = false
  ld.closeCallScreen()
}

const handleGameOverClose = () => {
  ld.showGameOver.value = false
}

const handleLobbyBack = () => {
  router.replace({ name: 'room', params: { id: room.data.id } })
}
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
            v-if="ld.gameData.publicState?.gameOver"
          />
          <IconButton icon="fa6-solid:arrow-left" label="Leave" @click="leave" v-else />
        </div>
        <div class="title-box">
          <div class="is-flex">
            <h1 class="title mr-3">Liars Dice</h1>
            <div>
              <RulesPage />
            </div>
          </div>
        </div>
        <div class="outer">
          <IconButton
            icon="fa6-solid:clock-rotate-left"
            @click="showCall"
            label="Previous Round"
            v-if="(ld.gameData.publicState?.previousRound.round ?? 0) != 0"
          />
        </div>
      </div>
      <div class="body-wrapper">
        <div class="box mb-0 is-1">
          <h4 class="title is-5">Players - Dice Left: {{ totalDice }}</h4>
          <div class="mb-4">
            <div
              v-for="id in ld.gameData.publicState?.turnOrder"
              :key="id"
              class="is-flex is-justify-content-space-between"
            >
              <div>
                <span class="icon-text">
                  <span class="icon">
                    <Icon
                      icon="fa6-solid:arrow-right"
                      v-if="id == ld.gameData.publicState?.playerTurn"
                    />
                  </span>
                  <RoomName :id="id" kick />
                </span>
              </div>
              <div>
                <span>{{ ld.gameData.publicState?.diceAmounts[id] }} Dice</span>
              </div>
            </div>
          </div>
          <h4 class="title is-5">Out</h4>
          <div>
            <div v-for="id in outPlayers" :key="id">
              <RoomName :id="id" kick />
            </div>
          </div>
        </div>
        <div class="box is-5">
          <div class="mb-4" v-if="isIn(room.data.userId)">
            <h4 class="title is-4 mb-2">Hand</h4>
            <DiceHand
              :true-value="ld.gameData.privateState?.dice ?? []"
              :shuffle="ld.rollHand.value"
            />
          </div>
          <div v-else class="mb-4">
            <span>You are now out</span>
          </div>
          <div>
            <h4 class="title is-4 mb-2">Highest Bid</h4>
            <div class="is-flex is-align-items-center" v-if="ld.gameData.publicState?.highestBid">
              <span class="is-size-3 mr-5">{{ ld.gameData.publicState?.bidAmount }}x</span>
              <DiceCube :value="ld.gameData.publicState?.bidFace ?? 0" />
            </div>
            <div v-else>None...</div>
          </div>
          <div v-if="ld.gameData.isTurn">
            <hr />
            <div class="is-size-5 mb-3">Its your turn! Please input a move:</div>
            <div class="is-flex">
              <BidForm @bid="handleBid" :current-bid="ld.gameData.publicState?.highestBid ?? ''" />
              <div>
                <span class="is-size-5 mx-5">or</span>
                <button class="button is-primary" @click="handleCall">Call</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </main>
  <CallModal
    :show="ld.showCall.value"
    @close="handleCallScreenClose"
    :data="ld.gameData.publicState?.previousRound"
  />
  <EndGame
    :show="ld.showGameOver.value"
    :last-player="room.users.names[ld.gameData.publicState?.turnOrder[0] ?? '']"
    @close="handleGameOverClose"
    @back="handleLobbyBack"
  />
</template>
