import { useErrorStore } from '@/stores/error'
import { parseData, sendMessage, validateUUID } from './main'
import { OPacketType } from './packetTypes'
import { z } from 'zod'
import { useRouter } from 'vue-router'
import router from '@/router'

const create = () => {
  sendMessage({
    type: OPacketType.CreateLobby,
    data: null,
  })
}

const join = (id: string) => {
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
    data: id,
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

const lobbyUserSchema = z.array(z.string())

export const handleLobbyChange = (data: unknown) => {
  const parsedData = parseData(data, lobbyUserSchema)
  onLobbyChange(parsedData)
}

const lobbyIdSchema = z.string().uuid()

export const handleLobbyJoin = (data: unknown) => {
  const parsedData = parseData(data, lobbyIdSchema)
  router.push({ name: 'lobby', params: { id: parsedData } })
}

const unsetLobbyChange = (users: string[]) => {
  console.error('no method defined')
}

let onLobbyChange = unsetLobbyChange

const setOnLobbyChange = (func: typeof onLobbyChange) => {
  onLobbyChange = func
}

const clearOnLobbyChange = () => {
  onLobbyChange = unsetLobbyChange
}

export default { create, join, leave, users, setOnLobbyChange, clearOnLobbyChange }