import { handleUserChange, handleUserInit } from './user'
import error from './error'
import { IPacketType } from './packetTypes'
import { z } from 'zod'

const notImplemented = (_: unknown) => {
  console.error('packet type not implemented')
}

const routeMap: Record<IPacketType, (data: unknown) => void> = {
  [IPacketType.UserInit]: handleUserInit,
  [IPacketType.UserChange]: handleUserChange,
  [IPacketType.Chat]: notImplemented,
  [IPacketType.LobbyChange]: notImplemented,
  [IPacketType.Error]: error.handle,
}

const incomingSchema = z.object({
  type: z.nativeEnum(IPacketType),
  data: z.unknown(),
})

export const handleIncomingMessage = (message: string) => {
  const msg = JSON.parse(message)

  const result = incomingSchema.safeParse(msg)

  if (!result.success) {
    // TODO: better error
    console.log(result.error.format())
    return
  }

  routeMap[result.data.type](result.data.data)
}
