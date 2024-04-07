import { useErrorStore } from '@/stores/error'
import { CallBackFunc, parseData, sendMessage, validateUUID } from './main'
import { OPacketType } from './packetTypes'
import { z } from 'zod'
import { gameTypes } from '@/game'

const create = (name?: string) => {
  sendMessage({
    type: OPacketType.CreateLobby,
    data: {
      name,
    },
  })
}

const join = (id: string, name?: string) => {
  const es = useErrorStore()
  if (!validateUUID(id)) {
    es.add({
      type: 'warning',
      message: 'Found invalid id while trying to join lobby',
    })
    return
  }

  sendMessage({
    type: OPacketType.JoinLobby,
    data: {
      id,
      name,
    },
  })
}

const leave = (id: string) => {
  const es = useErrorStore()

  if (!validateUUID(id)) {
    es.add({
      type: 'warning',
      message: 'Found invalid id while trying to leave lobby',
    })
    return
  }

  sendMessage({
    type: OPacketType.LeaveLobby,
    data: id,
  })
}

const users = (id: string) => {
  const es = useErrorStore()

  if (!validateUUID(id)) {
    es.add({
      type: 'warning',
      message: 'Found invalid id while trying to get users',
    })
    return
  }

  sendMessage({
    type: OPacketType.LobbyUsers,
    data: id,
  })
}

const lobbySchema = z.object({
  id: z.string().uuid(),
  maxPlayers: z.number().int().nullable(),
  minPlayers: z.number().int().nullable(),
  name: z.string().nullable(),
  gameType: z.union([z.enum(gameTypes), z.literal('')]).nullable(),
  inGame: z.boolean().nullable(),
  gameOptions: z.string().nullable(),
})

export type LobbyData = z.infer<typeof lobbySchema>
const lobbyChangeCB = new CallBackFunc<LobbyData>()

export const handleLobbyChange = (data: unknown) => {
  const parsedData = parseData(data, lobbySchema)
  lobbyChangeCB.run(parsedData)
}

const lobbyUserSchema = z.record(z.string().uuid(), z.string())

export type LobbyUsers = z.infer<typeof lobbyUserSchema>
const lobbyUserChangeCB = new CallBackFunc<LobbyUsers>()

export const handleLobbyUserChange = (data: unknown) => {
  const parsedData = parseData(data, lobbyUserSchema)
  lobbyUserChangeCB.run(parsedData)
}

const info = () => {
  sendMessage({
    type: OPacketType.LobbyInfo,
    data: undefined,
  })
}

export const lobbyDataSchema = z.object({
  id: z.string().uuid(),
  gameType: z.union([z.literal(''), z.enum(gameTypes)]),
})

type LobbyDataMessage = z.infer<typeof lobbyDataSchema>

const change = (data: LobbyDataMessage) => {
  sendMessage({
    type: OPacketType.LobbyChange,
    data,
  })
}

const options = (id: string, data: Record<string, any>) => {
  sendMessage({
    type: OPacketType.LobbyOptions,
    data: {
      id,
      data,
    },
  })
}

export default {
  create,
  join,
  leave,
  users,
  info,
  change,
  options,
  lobbyUserChangeCB,
  lobbyChangeCB,
}
