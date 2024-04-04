import api from '@/api'
import type { LobbyUsers, LobbyData } from '@/api/lobby'
import type { GameTypes } from '@/game'
import { defineStore } from 'pinia'
import { ref } from 'vue'
import router from '@/router'

const nilUUID = '00000000-0000-0000-0000-000000000000'

export const useLobbyStore = defineStore('lobby', () => {
  const id = ref<string>(nilUUID)
  const gameType = ref<GameTypes | ''>('')
  const inGame = ref(false)
  const users = ref<LobbyUsers>({})

  const handleLobbyUserChange = (iUsers: LobbyUsers) => {
    users.value = iUsers
  }
  api.lobby.lobbyUserChangeCB.fn = handleLobbyUserChange

  const handleLobbyChange = (d: LobbyData) => {
    if (!api.isNilUUID(d.id)) {
      if (api.isNilUUID(id.value)) {
        api.lobby.users(d.id)
      }
      id.value = d.id

      if (d.inGame) {
        router.replace({ name: 'game' })
      } else {
        router.replace({ name: 'lobby' })
      }
    } else {
      router.replace({ name: 'home' })
    }
    if (!d.inGame && inGame.value) {
      router.replace({ name: 'lobby' })
    }

    if (d.gameType) {
      gameType.value = d.gameType
    }

    if (d.inGame !== inGame.value) {
      inGame.value = d.inGame ?? false
    }
  }

  const getInfo = () => {
    api.lobby.info()
  }
  api.lobby.lobbyChangeCB.fn = handleLobbyChange

  const clear = () => {
    id.value = nilUUID
    gameType.value = ''
    users.value = {}
    inGame.value = false
  }

  const leave = () => {
    api.lobby.leave(id.value)
    clear()
  }

  return { id, users, gameType, inGame, getInfo, clear, leave }
})
