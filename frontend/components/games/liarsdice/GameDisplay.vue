<script setup lang="ts">
import ld from '@/game/liarsdice'
import { computed, onMounted, ref } from 'vue'
import IconButton from '@/components/IconButton.vue'
import DiceHand from './DiceHand.vue'
import DiceCube from './DiceCube.vue'
import BidForm from './BidForm.vue'
import ErrorStore from '@/components/ErrorStore.vue'
import { useLobbyStore } from '@/stores/lobby'
import { Icon } from '@iconify/vue'
import CallModal from './CallModal.vue'
import { useUserStore } from '@/stores/user'

const lobby = useLobbyStore()
const user = useUserStore()

onMounted(async () => {
  console.log('REMOUNTED', lobby.id)
  ld.create(lobby.id)
})

const outPlayers = computed(() => {
  const ids: string[] = []
  for (const k in lobby.users) {
    if (!ld.gameData.publicState?.turnOrder.includes(k)) {
      ids.push(k)
    }
  }
  return ids
})

const handleBid = (bid: string) => {
  ld.bid(lobby.id, bid)
}

const handleCall = () => {
  ld.call(lobby.id)
}

const leave = () => {
  lobby.leave()
}

const totalDice = computed(() => {
  if (ld.gameData.publicState?.diceAmounts) {
    return Object.values(ld.gameData.publicState.diceAmounts).reduce((a, b) => a + b)
  }
  return 0
})

const handleCallScreenClose = () => {
  ld.showCall.value = false
}
</script>

<template>
  <main class="centerize">
    <ErrorStore />
    <div class="box b-primary container">
      <div class="logo-header mb-6">
        <div class="outer">
          <IconButton icon="fa6-solid:arrow-left" label="Leave" @click="leave" />
        </div>
        <div><h1 class="title">Liars Dice</h1></div>
        <div class="outer"></div>
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
                  <span :class="{ 'has-text-weight-bold': id === user.data.id }">
                    {{ lobby.users[id]?.name }}
                  </span>
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
              <span> ({{ lobby.users[id].name }}) </span>
            </div>
          </div>
        </div>
        <div class="box is-5">
          <div class="mb-4">
            <h4 class="title is-4 mb-2">Hand</h4>
            <DiceHand :true-value="ld.gameData.privateState?.dice ?? []" />
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
  <CallModal :show="ld.showCall.value" @close="handleCallScreenClose" />
</template>
