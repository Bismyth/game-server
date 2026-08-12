<script setup lang="ts">
import ErrorStore from '@/components/ErrorStore.vue'
import skull from '@/game/skull'
import { useRoomStore } from '@/stores/room'
import { onMounted, ref, watch, computed } from 'vue'
import { handleLobbyBack } from '@/game'
import RulesPage from '@/components/games/RulesPage.vue'
import IconButton from '@/components/IconButton.vue'
import { Icon } from '@iconify/vue'
import RoomName from '@/components/RoomName.vue'
import TileHand from './TileHand.vue'

const room = useRoomStore()

onMounted(async () => {
  skull.create()
})

watch(
  () => room.data.inGame,
  (v, ov) => {
    if (v && !ov) {
      skull.create()
    }
  },
)
const bid = ref('')

const makeBid = () => {
  skull.bid(parseInt(bid.value))
  bid.value = ''
}

const outPlayers = computed(() => {
  const ids: string[] = []
  for (const k of room.users.order) {
    if (!skull.gameData.publicState?.turnOrder.includes(k)) {
      ids.push(k)
    }
  }
  return ids
})

const pass = () => {
  skull.pass()
  bid.value = ''
}

const flip = (sP: string) => {
  skull.flip(sP)
}

const hTile = (v: boolean) => {
  skull.place(v)
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
            v-if="skull.gameData.publicState?.gameOver"
          />
          <IconButton icon="fa6-solid:arrow-left" label="Leave" @click="room.leave" v-else />
        </div>
        <div class="title-box">
          <div class="is-flex">
            <h1 class="title mr-3">Skull</h1>
            <div>
              <RulesPage game-type="skull" />
            </div>
          </div>
        </div>
        <div class="outer"></div>
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
                  <th :style="{ width: '40px' }">Hand</th>
                  <th :style="{ width: '40px' }">Pts</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="id in skull.gameData.publicState?.turnOrder" :key="id">
                  <td>
                    <span class="icon">
                      <Icon
                        icon="fa6-solid:star"
                        v-if="id == skull.gameData.publicState?.flipper"
                      />
                      <Icon
                        icon="fa6-solid:arrow-right"
                        v-else-if="id == skull.gameData.publicState?.turn"
                      />
                      <Icon icon="fa6-solid:x" v-else-if="skull.passed(id)" />
                    </span>
                  </td>
                  <td><RoomName :id="id" kick /></td>
                  <td>{{ skull.gameData.publicState?.tileCount[id] }}</td>
                  <td>{{ skull.gameData.publicState?.points[id] }}</td>
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
          <div v-if="skull.userActions.value.showMessage" class="is-flex mb-3 is-size-4">
            {{ skull.userActions.value.message }}
          </div>
          <div class="mb-4">
            <div v-for="id in skull.gameData.publicState?.turnOrder" class="mb-1">
              <RoomName :id="id" />

              <TileHand
                :true-value="skull.gameData.publicState?.tilesRevealed[id]"
                :size="skull.gameData.publicState?.tilesPlaced[id]"
                placed
                :clickable="
                  skull.gameData.publicState?.flipper == room.data.userId &&
                  skull.canFlip(id) &&
                  !skull.hasNextRound.value
                "
                @select-hand="flip(id)"
              />
            </div>
          </div>

          <span>Current Hand</span>
          <TileHand :true-value="skull.currentHand.value" hand @select-tile="hTile" />
          <div class="is-flex mt-4">
            <div v-if="skull.userActions.value.showBid" class="is-flex">
              <p class="control">
                <input class="input" v-model="bid" />
              </p>

              <button class="button is-primary" @click="makeBid">Bid</button>
            </div>
            <div v-if="skull.userActions.value.showPass">
              <button class="button is-primary" @click="pass">Pass</button>
            </div>
          </div>
          <button class="button" v-if="skull.hasNextRound.value" @click="skull.nextRound">
            Next Round
          </button>
        </div>
      </div>
    </div>
  </main>
</template>
