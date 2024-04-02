/* eslint-disable prefer-const */
import game, { handleGameState } from '@/api/game'
import { reactive, ref, type Ref } from 'vue'
import { z } from 'zod'

const gameType = 'liarsdice'

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

let state: { id: string; type: 'liarsdice' }

const create = (id: string) => {
  state = {
    id: id,
    type: gameType,
  }

  game.setHandleAction(handleAction)
  game.setHandleEvent(handleEvent)
  game.setHandleState(handleState)

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
  game.ready(state)
}

const takeAction = (option: string, data: any) => {
  if (!gameData.isTurn) {
    return
  }
  game.action(state, option, data)
  gameData.isTurn = false
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
}
const handleAction = (data: unknown) => {
  const result = z.array(z.string()).safeParse(data)
  if (!result.success) {
    //Todo: better error response
    console.error('bad data')
    return
  }
  gameData.isTurn = true
  gameData.currentOptions = result.data
}
const handleEvent = (data: unknown) => {
  console.log(data)
}

export default { create, bid, call, gameData }

const optionsSchema = z.object({
  startingDice: z.number().int().min(1).max(99),
})

export const createGame = (lobbyId: string, options: { startingDice: number }) => {
  game.newGame(gameType, lobbyId, options)
}
