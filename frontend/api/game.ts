import { z } from 'zod'
import { CallBackFunc, sendMessage } from './main'
import { OPacketType } from './packetTypes'
import { useErrorStore } from '@/stores/error'
import { useLobbyStore } from '@/stores/lobby'

const handleAction = new CallBackFunc<unknown>()
export const handleGameAction = (data: unknown) => {
  handleAction.run(data)
}

const handleEvent = new CallBackFunc<unknown>()
export const handleGameEvent = (data: unknown) => {
  handleEvent.run(data)
}

const handleState = new CallBackFunc<unknown>()
export const handleGameState = (data: unknown) => {
  handleState.run(data)
}

const newGame = (lobbyId: string) => {
  const es = useErrorStore()

  const result = z.string().uuid().safeParse(lobbyId)
  if (!result.success) {
    es.add({
      type: 'warning',
      message: 'invalid lobbyid to start game',
    })
    return
  }

  sendMessage({
    type: OPacketType.GameNew,
    data: lobbyId,
  })
}

const action = (lobbyId: string, option: string, data: any) => {
  sendMessage({
    type: OPacketType.GameAction,
    data: {
      id: lobbyId,
      option,
      data,
    },
  })
}

const ready = (lobbyId: string) => {
  sendMessage({
    type: OPacketType.GameReady,
    data: {
      id: lobbyId,
    },
  })
}

export default { newGame, action, ready, handleAction, handleEvent, handleState }
