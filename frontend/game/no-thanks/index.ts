/* eslint-disable prefer-const */
import api from '@/api-utils'
import { useRoomStore } from '@/stores/room'
import { computed, reactive, ref } from 'vue'
import { z } from 'zod'


const publicStateSchema = z.object({
  inPlayCard: z.number().int(),
  tokensOnCard: z.number().int(),
  playerCards: z.record(z.uuid(), z.array(z.number().int())),
  score: z.record(z.uuid(), z.number().int()),
  deckLeft: z.number().int(),
  turnOrder: z.array(z.uuid()),
  currentPlayer: z.uuid(),
  currentRound: z.number().int(),
  totalRounds: z.number().int(),
})

type publicStateT = z.infer<typeof publicStateSchema>

const privateStateScehma = z.object({
    tokens: z.number().int()
})

type privateStateT = z.infer<typeof privateStateScehma>

const stateSchema = z.object({
  public: publicStateSchema.nullable(),
  private: privateStateScehma.nullable(),
})

const create = () => {
  api.game.handleAction.fn = handleAction
  api.game.handleState.fn = handleState

  resetValues()

  api.game.ready()
}

const gameData = reactive<{
  publicState: (publicStateT) | undefined
  privateState: privateStateT | undefined
  isTurn: boolean
  currentOptions: string[]
}>({
  publicState: undefined,
  privateState: undefined,
  isTurn: false,
  currentOptions: [],
})

const resetValues = () => {
  gameData.publicState = undefined
  gameData.privateState = undefined
  gameData.isTurn = false
  gameData.currentOptions = []
}


const takeAction = (option: string, data?: any) => {
  if (!gameData.isTurn) {
    return
  }
  api.game.action(option, data)
}

const pass = () => {
  takeAction('pass')
}

const take = () => {
  takeAction('take')
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

const handleState = (data: unknown) => {
  const result = stateSchema.safeParse(data)
  if (!result.success) {
    console.log(z.treeifyError(result.error))
    console.log("bad data")
    console.log(data)
    console.log(result.error)
    return
  }

  if (result.data.public) {
    gameData.publicState = result.data.public
    const room = useRoomStore()

    if (result.data.public.currentPlayer !== room.data.userId) {
      gameData.isTurn = false
    }
  }
  if (result.data.private) {
    gameData.privateState = result.data.private
  }
}

export default {
    gameData,
    create,
    pass,
    take,
}