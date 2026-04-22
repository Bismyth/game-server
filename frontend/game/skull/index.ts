/* eslint-disable prefer-const */
import api from '@/api'
import { useRoomStore } from '@/stores/room'
import { computed, reactive, ref } from 'vue'
import { z } from 'zod'

const publicStateSchema = z.object({
  tilesPlaced: z.record(z.uuid(), z.number().int()),
  tilesRevealed: z.record(z.uuid(), z.array(z.boolean())),
  bid: z.number().int(),
  passed: z.array(z.uuid()),
  points: z.record(z.uuid(), z.number().int()),
  flipper: z.uuid(),
  gameOver: z.boolean(),
  turnOrder: z.array(z.string()),
  turn: z.uuid(),
  round: z.number().int(),
})

type publicStateT = z.infer<typeof publicStateSchema>

const privateStateScehma = z.object({
  tiles: z.array(z.boolean()),
  tilesPlaced: z.array(z.boolean()),
})

const stateSchema = z.object({
  public: publicStateSchema.nullable(),
  private: privateStateScehma.nullable(),
})

type privateStateT = z.infer<typeof privateStateScehma>

type fullStateT = z.infer<typeof stateSchema>

const create = () => {
  api.game.handleAction.fn = handleAction
  api.game.handleEvent.fn = handleEvent
  api.game.handleState.fn = handleState

  resetValues()

  ready()
}

const resetValues = () => {
  showGameOver.value = false
  gameData.publicState = undefined
  gameData.privateState = undefined
  gameData.currentOptions = []
}

interface GameData {
  publicState: publicStateT | undefined
  privateState: privateStateT | undefined
  isTurn: boolean
  currentOptions: string[]
}

const gameData = reactive<GameData>({
  publicState: undefined,
  privateState: undefined,
  isTurn: false,
  currentOptions: [],
})

const newRoundData = ref<fullStateT | undefined>()

const showCall = ref(false)
const showGameOver = ref(false)

const ready = () => {
  api.game.ready()
}

const takeAction = (option: string, data?: any) => {
  api.game.action(option, data)
}

const place = (tile: boolean) => {
  takeAction('place', { tile })
}

const bid = (amount: number) => {
  takeAction('bid', { bid: amount })
}

const pass = () => {
  takeAction('pass')
}

const flip = (player: string) => {
  takeAction('flip', { player })
}

const handleState = (data: unknown) => {
  const room = useRoomStore()

  const result = stateSchema.safeParse(data)
  if (!result.success) {
    //todo: better error
    console.log(result.error.format())
    console.error('bad data')
    return
  }

  if (result.data.public?.gameOver && !showGameOver.value) {
    showGameOver.value = true
  }

  if (result.data.public) {
    if (result.data.public.round > 1 && result.data.public.round !== gameData.publicState?.round) {
      newRoundData.value = result.data
    } else {
      gameData.publicState = result.data.public
    }
  }
  if (result.data.private) {
    gameData.privateState = result.data.private
  }

  gameData.isTurn = result.data.public?.turn == room.data.userId
}

const nextRound = () => {
  console.log(hasNextRound.value)
  console.log(newRoundData.value)
  if (newRoundData.value === undefined) {
    return
  }

  if (newRoundData.value.public) {
    gameData.publicState = newRoundData.value.public
  }
  if (newRoundData.value.private) {
    gameData.privateState = newRoundData.value.private
  }

  newRoundData.value = undefined
}

const hasNextRound = computed(() => newRoundData.value !== undefined)

const handleAction = (data: unknown) => {
  const result = z.array(z.string()).safeParse(data)
  if (!result.success) {
    //Todo: better error response
    console.error('bad data')
    return
  }
  gameData.currentOptions = result.data
  gameData.isTurn = true
}

const handleEvent = (data: unknown) => {
  console.log(data)
}

const currentHand = computed(() => {
  const totalTiles = gameData.privateState?.tiles ?? []
  const placedTiles = gameData.privateState?.tilesPlaced ?? []

  const hasSkull = totalTiles.includes(true)
  const placedSkull = placedTiles.includes(true)

  const cHand = []

  if (hasSkull && !placedSkull) {
    cHand.push(true)
  }
  const restSize = totalTiles.length - placedTiles.length - cHand.length
  for (let x = 0; x < restSize; x++) {
    cHand.push(false)
  }

  console.log(totalTiles, placedTiles, cHand)

  return cHand
})

export default {
  create,
  ready,
  place,
  bid,
  pass,
  flip,
  gameData,
  showCall,
  showGameOver,
  currentHand,
  nextRound,
  hasNextRound,
}
