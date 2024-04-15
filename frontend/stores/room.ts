import api from '@/api'
import type { RoomUsers, RoomInfo } from '@/api/room'
import type { GameTypes } from '@/game'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import router from '@/router'
import { isNilUUID } from '@/api/main'

const nilUUID = '00000000-0000-0000-0000-000000000000'

export const useRoomStore = defineStore('room', () => {
  const id = ref<string>(nilUUID)
  const userId = ref<string>(nilUUID)
  const gameType = ref<GameTypes | ''>('')
  const inGame = ref(false)
  const users = ref<RoomUsers>({
    order: [],
    names: {},
  })
  const ready = ref(false)
  const options = ref<any>()
  const host = ref<string>(nilUUID)
  const isHost = computed(() => host.value === userId.value)

  id.value = router.currentRoute.value.params.id.toString()
  const token = localStorage.getItem(`room:${id.value}`) ?? ''
  try {
    userId.value = JSON.parse(atob(token.split('.')[1] ?? ''))?.userId
  } catch (e) {
    console.log('failed to parse token')
  }

  const wsProtocol = location.protocol === 'http:' ? 'ws:' : 'wss:'

  const conn = new WebSocket(`${wsProtocol}//${document.location.host}/ws`)
  conn.onopen = () => {
    conn.send(token)
  }

  const send = (payload: any) => {
    conn.send(JSON.stringify(payload))
  }

  api.setSendMessage(send)

  conn.onmessage = (evt) => api.handleIncomingMessage(evt.data)

  const handleRoomUserChange = (iUsers: RoomUsers) => {
    users.value = iUsers
  }
  api.room.roomUserChangeCB.fn = handleRoomUserChange

  const handleLobbyChange = (d: RoomInfo) => {
    if (d.inGame !== null) {
      inGame.value = d.inGame
    }

    if (users.value.order.length === 0) {
      api.room.users()
    }

    // if (router.currentRoute.value.name !== 'game') {
    //   if (inGame.value) {
    //     router.replace({ name: 'game' })
    //   } else {
    //     router.replace({ name: 'lobby' })
    //   }
    // }

    if (d.gameType) {
      if (gameType.value !== d.gameType) {
        options.value = undefined
        gameType.value = d.gameType
      }
    }

    if (!isNilUUID(d.host)) {
      host.value = d.host
    }

    if (d.gameOptions) {
      options.value = JSON.parse(d.gameOptions)
    }

    ready.value = true
  }
  api.room.roomChangeCB.fn = handleLobbyChange

  const clear = () => {
    id.value = nilUUID
    gameType.value = ''
    users.value = {
      order: [],
      names: {},
    }
    inGame.value = false
    ready.value = false
  }

  const leave = () => {
    clear()
    api.room.leave()
  }

  return { id, users, gameType, inGame, ready, options, clear, leave, isHost, host, userId }
})
