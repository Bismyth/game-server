import { defineStore } from 'pinia'
import { ref } from 'vue'

import api from '@/api'
import { useErrorStore } from './error'

export const useSocketStore = defineStore('socket', () => {
  const conn = new WebSocket(`ws://${document.location.host}/ws`)

  const userReady = ref(false)

  const errorStore = useErrorStore()

  const send = (payload: any) => {
    if (!userReady.value) {
      console.error('tried to send to inactive connection', payload)
      return
    }
    conn.send(JSON.stringify(payload))
  }

  api.setSendMessage(send)

  conn.onopen = () => {
    conn.send(api.user.getLocalId() ?? '')
  }

  conn.onclose = () => {
    errorStore.add({
      type: 'danger',
      message: 'No active websocket connection, please reload the page',
      noexpire: true,
    })
  }

  conn.onmessage = function (evt) {
    api.handleIncomingMessage(evt.data)
  }

  return { conn, send, userReady }
})
