import api from '@/api'
import type { RoomUsers, RoomInfo } from '@/api/room'
import type { GameTypes } from '@/game'
import { defineStore } from 'pinia'
import { computed, reactive, ref } from 'vue'
import router from '@/router'
import { isNilUUID } from '@/api/main'
import { useErrorStore } from './error'

const nilUUID = '00000000-0000-0000-0000-000000000000'

type RoomData = {
  id: string
  userId: string
  gameType: GameTypes | ''
  options: any
  inGame: boolean
  host: string
}

const emptyRoom: RoomData = {
  id: nilUUID,
  userId: nilUUID,
  gameType: '',
  inGame: false,
  options: undefined,
  host: nilUUID,
}

const emptyUsers = {
  order: [],
  names: {},
}

export const useRoomStore = defineStore('room', () => {
  const data = reactive<RoomData>(emptyRoom)
  const users = ref<RoomUsers>(emptyUsers)
  const ready = ref(false)

  const isHost = computed(() => data.host === data.userId)

  const es = useErrorStore()

  const setupConnection = () => {
    if (!ready.value) {
      data.id = router.currentRoute.value.params.id.toString()
      const token = localStorage.getItem(`room:${data.id}`) ?? ''
      try {
        data.userId = JSON.parse(atob(token.split('.')[1] ?? ''))?.userId
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

      conn.onclose = (ev) => {
        const leaveCode = [3001, 1006]

        if (ev.code === 1006) {
          es.add({
            type: 'danger',
            message: 'Websocket close abnormal',
          })
        }

        if (leaveCode.includes(ev.code)) {
          handleLeave()
          return
        }
        es.add({
          type: 'danger',
          message: 'No active websocket, please reload the page',
        })
      }
    }
  }

  const handleRoomUserChange = (iUsers: RoomUsers) => {
    users.value = iUsers
  }
  api.room.roomUserChangeCB.fn = handleRoomUserChange

  const handleLobbyChange = (d: RoomInfo) => {
    if (d.inGame !== null) {
      data.inGame = d.inGame
    }

    if (users.value.order.length === 0) {
      api.room.users()
    }

    if (router.currentRoute.value.name === 'room') {
      if (data.inGame) {
        router.replace({ name: 'game', params: { id: data.id } })
      }
    }

    if (d.gameType !== null) {
      if (data.gameType !== d.gameType) {
        data.options = undefined
        data.gameType = d.gameType
      }
    }

    if (!isNilUUID(d.host)) {
      data.host = d.host
    }

    if (d.gameOptions) {
      data.options = JSON.parse(d.gameOptions)
    } else {
      data.options = undefined
    }

    ready.value = true
  }
  api.room.roomChangeCB.fn = handleLobbyChange

  const clear = () => {
    console.log('clearing room')
    Object.assign(data, emptyRoom)
    users.value = emptyUsers
    ready.value = false
  }

  const handleLeave = () => {
    clear()
    router.replace({ name: 'home' })
  }

  const leave = () => {
    api.room.leave()
    handleLeave()
  }

  return { data, users, ready, clear, leave, isHost, setupConnection }
})
