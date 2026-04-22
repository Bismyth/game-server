<script setup lang="ts">
import ErrorStore from '@/components/ErrorStore.vue'
import skull from '@/game/skull'
import { useRoomStore } from '@/stores/room'
import { onMounted, ref, watch } from 'vue'
import TileDisplay from './TileDisplay.vue'
import { handleLobbyBack } from '@/game'
import RulesPage from './RulesPage.vue'
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

const place = () => {
  skull.place(tilePlace.value)
}

const bid = ref('')

const tilePlace = ref(false)

const makeBid = () => {
  skull.bid(parseInt(bid.value))
}

const pass = () => {
  skull.pass()
}

const flip = (sP: string) => {
  console.log('flipping', sP)
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
              <RulesPage />
            </div>
          </div>
        </div>
        <div class="outer">
          <!-- <IconButton
            icon="fa6-solid:clock-rotate-left"
            @click="showCall"
            label="Previous Round"
            v-if="(ld.gameData.publicState?.previousRound.round ?? 0) != 0"
          /> -->
        </div>
      </div>
      <div class="body-wrapper">
        <div class="box mb-0 is-1">
          <h4 class="title is-5">Players</h4>
          <div class="mb-4">
            <div
              v-for="id in skull.gameData.publicState?.turnOrder"
              :key="id"
              class="is-flex is-justify-content-space-between"
            >
              <div>
                <span class="icon-text">
                  <span class="icon">
                    <Icon
                      icon="fa6-solid:arrow-right"
                      v-if="id == skull.gameData.publicState?.turn"
                    />
                  </span>
                  <RoomName :id="id" kick />
                </span>
              </div>
            </div>
          </div>
          <!-- <h4 class="title is-5">Out</h4>
          <div>
            <div v-for="id in outPlayers" :key="id">
              <RoomName :id="id" kick />
            </div>
          </div> -->
        </div>
        <div class="box is-5">
          <div v-for="id in skull.gameData.publicState?.turnOrder">
            <RoomName :id="id" />

            <TileHand
              :true-value="skull.gameData.publicState?.tilesRevealed[id]"
              :size="skull.gameData.publicState?.tilesPlaced[id]"
              placed
              :clickable="skull.gameData.publicState?.flipper == room.data.userId"
              @select-hand="flip(id)"
            />
          </div>

          <span>Current Hand</span>
          <TileHand :true-value="skull.currentHand.value" hand @select-tile="hTile" />
          <div class="is-flex mt-4">
            <p class="control">
              <input class="input" v-model="bid" />
            </p>

            <button class="button" @click="makeBid">Bid</button>

            <button class="button" @click="pass">Pass</button>
          </div>
          <button class="button" v-show="skull.hasNextRound.value" @click="skull.nextRound">
            Next Round
          </button>
        </div>
      </div>
    </div>
  </main>
</template>
