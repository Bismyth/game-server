/* eslint-disable prefer-const */
import api from '@/api'
import { useUserStore } from '@/stores/user'
import { reactive, ref } from 'vue'
import { z } from 'zod'

const roundInfoSchema = z.object({
  round: z.number().int(),
  highestBid: z.string(),
  hands: z.record(z.string(), z.array(z.number())).nullable(),
  callUser: z.string().uuid(),
  lastBid: z.string().uuid(),
  diceLost: z.string().uuid(),
  leave: z.string().nullable(),
})

const publicStateSchema = z.object({
  playerTurn: z.string().uuid(),
  diceAmounts: z.record(z.string(), z.number().int()),
  highestBid: z.string(),
  turnOrder: z.array(z.string()),
  gameOver: z.boolean(),
  previousRound: roundInfoSchema,
})

type publicStateT = z.infer<typeof publicStateSchema>

export type RoundInfo = z.infer<typeof roundInfoSchema>

const privateStateScehma = z.object({
  dice: z.array(z.number().int().min(1).max(6)),
})

const stateSchema = z.object({
  public: publicStateSchema.nullable(),
  private: privateStateScehma.nullable(),
})

type privateStateT = z.infer<typeof privateStateScehma>

let idStore = ''

const create = (id: string) => {
  idStore = id

  api.game.handleAction.fn = handleAction
  api.game.handleEvent.fn = handleEvent
  api.game.handleState.fn = handleState

  resetValues()

  ready(idStore)
}

const resetValues = () => {
  showCall.value = false
  showGameOver.value = false
  gameData.publicState = undefined
  gameData.privateState = undefined
  gameData.isTurn = false
  gameData.currentOptions = []
}

interface computedPublic {
  bidAmount: number
  bidFace: number
}

const gameData = reactive<{
  publicState: (publicStateT & computedPublic) | undefined
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
const rollHand = ref(false)
const rollHandTime = 3 * 1000 //5 seconds

let triggerRollHand = false

const ready = (lobbyId: string) => {
  api.game.ready(lobbyId)
}

const takeAction = (lobbyId: string, option: string, data: any) => {
  if (!gameData.isTurn) {
    return
  }
  api.game.action(lobbyId, option, data)
}

const bid = (lobbyId: string, bid: string) => {
  takeAction(lobbyId, 'bid', { bid })
}

const call = (lobbyId: string) => {
  takeAction(lobbyId, 'call', undefined)
}

const closeCallScreen = () => {
  if (triggerRollHand) {
    rollHand.value = true
    setTimeout(() => {
      rollHand.value = false
    }, rollHandTime)
    triggerRollHand = false
  }
}

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

  let newRound = false

  if (result.data.public) {
    newRound =
      result.data.public.previousRound.round !== gameData.publicState?.previousRound.round &&
      gameData.publicState !== undefined &&
      !result.data.public.gameOver
    if (newRound) {
      showCall.value = true
    }

    if (result.data.public.highestBid != '' && showCall.value) {
      showCall.value = false
    }

    const [a, f] = splitBid(result.data.public.highestBid)
    gameData.publicState = {
      ...result.data.public,
      bidAmount: a,
      bidFace: f,
    }
  }
  if (result.data.private) {
    if (gameData.privateState !== undefined) {
      triggerRollHand = true
    }
    gameData.privateState = result.data.private
  }

  const user = useUserStore()

  if (result.data.public) {
    gameData.isTurn = result.data.public.playerTurn === user.data.id
  }
}

const splitBid = (bid: string): [number, number] => {
  const bidSplit = bid.split(',')
  if (bidSplit.length !== 2) {
    return [0, 0]
  }

  return [parseInt(bidSplit[0]), parseInt(bidSplit[1])]
}

const handleAction = (data: unknown) => {
  const result = z.array(z.string()).safeParse(data)
  if (!result.success) {
    //Todo: better error response
    console.error('bad data')
    return
  }
  gameData.currentOptions = result.data
}
const handleEvent = (data: unknown) => {
  console.log(data)
}

export default { create, bid, call, gameData, showCall, showGameOver, rollHand, closeCallScreen }
