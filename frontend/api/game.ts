import { z } from 'zod'
import { sendMessage } from './main'
import { OPacketType } from './packetTypes'
import { useErrorStore } from '@/stores/error'

const undefinedFn = (data: unknown) => {
  console.error('not linked')
}

let ha = undefinedFn
export const handleGameAction = (data: unknown) => {
  ha(data)
}

let he = undefinedFn
export const handleGameEvent = (data: unknown) => {
  he(data)
}

let hs = undefinedFn
export const handleGameState = (data: unknown) => {
  hs(data)
}

const setHandleAction = (fn: typeof handleGameAction) => {
  ha = fn
}

const setHandleEvent = (fn: typeof handleGameAction) => {
  he = fn
}

const setHandleState = (fn: typeof handleGameAction) => {
  hs = fn
}

const newGame = (gameType: string, lobbyId: string, options: any) => {
  const validId = z.string().uuid().parse(lobbyId)

  sendMessage({
    type: OPacketType.GameNew,
    data: {
      type: gameType,
      lobbyId: validId,
      options,
    },
  })
}

const sharedStateSchema = z.object({
  id: z.string().uuid(),
  type: z.string(),
})

const action = (sharedState: { id: string; type: string }, option: string, data: any) => {
  const es = useErrorStore()

  const result = sharedStateSchema.safeParse(sharedState)
  if (!result.success) {
    es.add({
      type: 'warning',
      message: 'could not parse shared state',
    })
    return
  }

  sendMessage({
    type: OPacketType.GameAction,
    data: {
      ...result.data,
      option,
      data,
    },
  })
}

const ready = (sharedState: { id: string; type: string }) => {
  const es = useErrorStore()

  const result = sharedStateSchema.safeParse(sharedState)
  if (!result.success) {
    es.add({
      type: 'warning',
      message: 'could not parse shared state',
    })
    return
  }

  sendMessage({
    type: OPacketType.GameReady,
    data: result.data,
  })
}

export default { newGame, action, ready, setHandleAction, setHandleEvent, setHandleState }
