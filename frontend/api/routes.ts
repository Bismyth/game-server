import { errorHandle, kickMessage } from './error'
import { z } from 'zod'
import { IPacketType } from './packetTypes'
import { handleRoomChange, handleRoomUserChange } from './room'
import { useErrorStore } from '@/stores/error'
import { handleGameAction, handleGameEvent, handleGameState } from './game'

const routeMap: Record<IPacketType, (data: unknown) => void> = {
  [IPacketType.RoomInfo]: handleRoomChange,
  [IPacketType.RoomUserChange]: handleRoomUserChange,
  [IPacketType.Error]: errorHandle,
  [IPacketType.GameState]: handleGameState,
  [IPacketType.GameAction]: handleGameAction,
  [IPacketType.GameEvent]: handleGameEvent,
  [IPacketType.GameError]: errorHandle,
  [IPacketType.RoomKick]: kickMessage,
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
