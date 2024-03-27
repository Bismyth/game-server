import { defineStore } from 'pinia'
import { ref } from 'vue'

import api from '@/api'
import { useErrorStore } from './error'

export const useSocketStore = defineStore('socket', () => {
  const conn = new WebSocket('/ws', api.user.getLocalId() ?? '')

  const active = ref(false)

  const errorStore = useErrorStore()

  const send = (payload: any) => {
    conn.send(JSON.stringify(payload))
  }

  api.setSendMessage(send)

  conn.onopen = (evt) => {
    active.value = true
  }

  conn.onclose = () => {
    active.value = false

    errorStore.add({
      type: 'danger',
      message: 'No active websocket connection, please reload the page',
      noexpire: true,
    })
  }

  conn.onmessage = function (evt) {
    api.handleIncomingMessage(evt.data)
  }

  return { conn, send, active }
})
