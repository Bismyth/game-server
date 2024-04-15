import { CallBackFunc, parseData, sendMessage } from './main'
import { OPacketType } from './packetTypes'
import { z } from 'zod'
import { gameTypes } from '@/game'

const leave = () => {
  sendMessage({
    type: OPacketType.LeaveRoom,
    data: undefined,
  })
}

const users = () => {
  sendMessage({
    type: OPacketType.RoomUsers,
    data: undefined,
  })
}

const roomSchema = z.object({
  maxPlayers: z.number().int().nullable(),
  host: z.string().uuid(),
  name: z.string().nullable(),
  gameType: z.union([z.enum(gameTypes), z.literal('')]).nullable(),
  inGame: z.boolean().nullable(),
  gameOptions: z.string().nullable(),
})

export type RoomInfo = z.infer<typeof roomSchema>
const roomChangeCB = new CallBackFunc<RoomInfo>()

export const handleRoomChange = (data: unknown) => {
  const parsedData = parseData(data, roomSchema)
  roomChangeCB.run(parsedData)
}

const roomUserSchema = z.object({
  order: z.array(z.string().uuid()),
  names: z.record(z.string().uuid(), z.string()),
})

export type RoomUsers = z.infer<typeof roomUserSchema>
const roomUserChangeCB = new CallBackFunc<RoomUsers>()

export const handleRoomUserChange = (data: unknown) => {
  const parsedData = parseData(data, roomUserSchema)
  roomUserChangeCB.run(parsedData)
}

export const roomDataSchema = z.object({
  gameType: z.union([z.literal(''), z.enum(gameTypes)]),
})

type RoomDataMessage = z.infer<typeof roomDataSchema>

const change = (data: RoomDataMessage) => {
  sendMessage({
    type: OPacketType.RoomChange,
    data,
  })
}

const options = (data: Record<string, any>) => {
  sendMessage({
    type: OPacketType.RoomOptions,
    data,
  })
}

export default {
  leave,
  users,
  change,
  options,
  roomUserChangeCB,
  roomChangeCB,
}
