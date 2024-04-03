import { handleUserChange, handleUserInit } from './user'
import error from './error'
import { z } from 'zod'
import { IPacketType } from './packetTypes'
import { handleLobbyChange } from './lobby'
import { useErrorStore } from '@/stores/error'
import { handleGameAction, handleGameEvent, handleGameState } from './game'

const notImplemented = () => {
  console.error('packet type not implemented')
}

const routeMap: Record<IPacketType, (data: unknown) => void> = {
  [IPacketType.UserInit]: handleUserInit,
  [IPacketType.UserChange]: handleUserChange,
  [IPacketType.Chat]: notImplemented,
  [IPacketType.LobbyChange]: handleLobbyChange,
  [IPacketType.Error]: error.handle,
  [IPacketType.GameState]: handleGameState,
  [IPacketType.GameAction]: handleGameAction,
  [IPacketType.GameEvent]: handleGameEvent,
}

const incomingSchema = z.object({
  type: z.nativeEnum(IPacketType),
  data: z.unknown(),
})

export const handleIncomingMessage = (message: string) => {
  const msg = JSON.parse(message)

  const result = incomingSchema.safeParse(msg)

  const es = useErrorStore()

  if (!result.success) {
    // TODO: better error
    es.add({
      type: 'warning',
      message: 'Unknown data received from server',
    })

    console.log(result.error.format())
    return
  }

  routeMap[result.data.type](result.data.data)
}
