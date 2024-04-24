/* eslint-disable prefer-const */
import api from '@/api'
import { reactive, ref } from 'vue'
import { z } from 'zod'

const publicStateSchema = z.object({
  tilesPlaced: z.record(z.string().uuid(), z.number().int()),
  tilesRevealed: z.record(z.string().uuid(), z.array(z.boolean())),
  bid: z.number().int(),
  passed: z.array(z.string().uuid()),
  points: z.record(z.string().uuid(), z.number().int()),
  flipper: z.string().uuid(),
  gameOver: z.boolean(),
  turnOrder: z.array(z.string()),
  turn: z.string().uuid(),
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

const gameData = reactive<{
  publicState: publicStateT | undefined
  privateState: privateStateT | undefined
  isTurn: boolean
  currentOptions: string[]
}>({
  publicState: undefined,
  privateState: undefined,
  isTurn: false,
  currentOptions: [],
})

const showCall = ref(false)
const showGameOver = ref(false)

const ready = () => {
  api.game.ready()
}

const takeAction = (option: string, data?: any) => {
  api.game.action(option, data)
}

const place = (tile: boolean) => {
  takeAction('place', { data: { tile } })
}

const bid = (amount: number) => {}

const handleState = (data: unknown) => {
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
    gameData.publicState = result.data.public
  }
  if (result.data.private) {
    gameData.privateState = result.data.private
  }
}

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

export default {
  create,
  ready,
  place,
  gameData,
  showCall,
  showGameOver,
}
