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

const action = (option: string, data: any) => {
  const lobby = useLobbyStore()

  sendMessage({
    type: OPacketType.GameAction,
    data: {
      id: lobby.id,
      option,
      data,
    },
  })
}

const ready = () => {
  const lobby = useLobbyStore()

  sendMessage({
    type: OPacketType.GameReady,
    data: {
      id: lobby.id,
    },
  })
}

export default { newGame, action, ready, handleAction, handleEvent, handleState }
