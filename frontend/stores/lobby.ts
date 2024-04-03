import api from '@/api'
import type { LobbyUsers } from '@/api/lobby'
import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export const useLobbyStore = defineStore('lobby', () => {
  const users = ref<LobbyUsers>({})

  const handleLobbyUserChange = (iUsers: LobbyUsers) => {
    users.value = iUsers
  }
  api.lobby.setOnLobbyChange(handleLobbyUserChange)

  const id = ref<string>('')

  watch(id, (newId) => {
    if (newId !== '' && !api.isNilUUID(newId)) {
      api.lobby.users(newId)
    }
  })

  return { id, users }
})
