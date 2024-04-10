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
  const ready = ref(false)
  const options = ref<any>()

  const handleLobbyUserChange = (iUsers: LobbyUsers) => {
    users.value = iUsers
  }
  api.lobby.lobbyUserChangeCB.fn = handleLobbyUserChange

  const handleLobbyChange = (d: LobbyData) => {
    if (api.isNilUUID(d.id)) {
      router.replace({ name: 'home' })
      return
    }

    if (api.isNilUUID(id.value)) {
      api.lobby.users(d.id)
    }
    id.value = d.id

    if (d.inGame !== null) {
      inGame.value = d.inGame
    }
    if (router.currentRoute.value.name !== 'game') {
      if (inGame.value) {
        router.replace({ name: 'game' })
      } else {
        router.replace({ name: 'lobby' })
      }
    }

    if (d.gameType) {
      if (gameType.value !== d.gameType) {
        options.value = undefined
        gameType.value = d.gameType
      }
    }

    if (d.gameOptions) {
      options.value = JSON.parse(d.gameOptions)
    }

    ready.value = true
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
    ready.value = false
  }

  const leave = () => {
    const oldId = id.value
    clear()
    api.lobby.leave(oldId)
  }

  return { id, users, gameType, inGame, ready, options, getInfo, clear, leave }
})
