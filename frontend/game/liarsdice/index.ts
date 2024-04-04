/* eslint-disable prefer-const */
import api from '@/api'
import { useUserStore } from '@/stores/user'
import { reactive } from 'vue'
import { z } from 'zod'

const publicStateSchema = z.object({
  playerTurn: z.string().uuid(),
  diceAmounts: z.record(z.string(), z.number().int()),
  highestBid: z.string(),
})

type publicStateT = z.infer<typeof publicStateSchema>

const privateStateScehma = z.object({
  dice: z.array(z.number().int().min(1).max(6)),
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

  ready()
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

const ready = () => {
  api.game.ready()
}

const takeAction = (option: string, data: any) => {
  if (!gameData.isTurn) {
    return
  }
  api.game.action(option, data)
}

const bid = (bid: string) => {
  takeAction('bid', { bid })
}

const call = () => {
  takeAction('call', undefined)
}

const handleState = (data: unknown) => {
  const result = stateSchema.safeParse(data)
  if (!result.success) {
    //todo: better error
    console.log(result.error.format())
    console.error('bad data')
    return
  }

  if (result.data.public) {
    gameData.publicState = result.data.public
  }
  if (result.data.private) {
    gameData.privateState = result.data.private
  }

  const user = useUserStore()

  if (result.data.public) {
    gameData.isTurn = result.data.public.playerTurn === user.data.id
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
}
const handleEvent = (data: unknown) => {
  console.log(data)
}

export default { create, bid, call, gameData }
