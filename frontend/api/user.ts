import { useUserStore, userSchema, type UserData } from '@/stores/user'
import { isNilUUID, parseData, sendMessage } from './main'
import { OPacketType } from './packetTypes'
import { useSocketStore } from '@/stores/socket'
import router from '@/router'
import { useLobbyStore } from '@/stores/lobby'

const idLSKey = 'id'

const handleRouteChange = (user: UserData) => {
  if (!isNilUUID(user.gameId)) {
    router.replace({ name: 'game' })
    return
  }

  if (!isNilUUID(user.lobbyId)) {
    router.replace({ name: 'lobby' })
    return
  }

  router.replace({ name: 'home' })
}

const setLobbyStateId = (id: string) => {
  const lobby = useLobbyStore()
  lobby.id = id
}

export const handleUserInit = (data: unknown) => {
  const parsedData = parseData(data, userSchema)

  localStorage.setItem(idLSKey, parsedData.id ?? 'wrong uuid')

  const user = useUserStore()
  user.data = parsedData

  const socket = useSocketStore()
  socket.active = true

  setLobbyStateId(user.data.lobbyId)
  handleRouteChange(user.data)
}

export const handleUserChange = (data: unknown) => {
  const parsedData = parseData(data, userSchema)

  const user = useUserStore()
  user.data = parsedData

  setLobbyStateId(user.data.lobbyId)
  handleRouteChange(user.data)
}

const getLocalId = () => {
  return localStorage.getItem(idLSKey)
}

const sendNameChange = (name: string) => {
  sendMessage({
    type: OPacketType.UserNameChange,
    data: {
      name,
    },
  })
}

export default { getLocalId, sendNameChange }
