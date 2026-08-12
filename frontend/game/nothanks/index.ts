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
  gameOver: z.boolean().optional(),
})

export type publicStateT = z.infer<typeof publicStateSchema>

const privateStateScehma = z.object({
    tokens: z.number().int()
})

type privateStateT = z.infer<typeof privateStateScehma>

const stateSchema = z.object({
  type: z.literal("state"),
  public: publicStateSchema.nullable(),
  private: privateStateScehma.nullable(),
})

type stateT = z.infer<typeof stateSchema>

const previousRoundSchema = z.object({
  type: z.literal("previous"),
  score: z.record(z.uuid(), z.number().int()),
  playerCards: z.record(z.uuid(), z.array(z.number().int())),
  removed: z.array(z.number().int()),
  round: z.number().int(),
})

export type PreviousRound = z.infer<typeof previousRoundSchema>


const messageSchema = z.discriminatedUnion("type", [
  stateSchema,
  previousRoundSchema
])

const create = () => {
  api.game.handleAction.fn = handleAction
  api.game.handleState.fn = handleMessage

  resetValues()

  api.game.ready()
}

const previousRound = ref<PreviousRound | undefined>(undefined)

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

const showPr = ref(false)

const resetValues = () => {
  gameData.publicState = undefined
  gameData.privateState = undefined
  gameData.isTurn = false
  gameData.currentOptions = []
  showPr.value = false
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

const handleMessage = (data: unknown) => {
  const result = messageSchema.safeParse(data)
  if (!result.success) {
    console.log(z.treeifyError(result.error))
    console.log("bad data")
    console.log(data)
    console.log(result.error)
    return
  }

  switch (result.data.type) {
    case 'state':
      handleState(result.data)
      break;
    case 'previous':
      previousRound.value = result.data
      if (gameData.publicState?.currentRound == result.data.round) {
        showPr.value = true
      }
      break;
  }
}

const handleState = (v: stateT) => {
    if (v.public) {
    gameData.publicState = v.public
    const room = useRoomStore()

    if (v.public.currentPlayer !== room.data.userId) {
      gameData.isTurn = false
    }
    if (v.public.gameOver) {
      showPr.value = true
    }
  }
  if (v.private) {
    gameData.privateState = v.private
  }
}

export default {
    gameData,
    create,
    pass,
    take,
    previousRound,
    showPr
}