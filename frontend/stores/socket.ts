import { defineStore } from 'pinia'
import { ref } from 'vue'

import id from '@/api/id'
import name from '@/api/name'
import { useErrorStore } from './error'

export const useSocketStore = defineStore('socket', () => {
  const conn = new WebSocket('/ws')

  const active = ref(false)

  const errorStore = useErrorStore()

  const send = (type: string, payload: any) => {
    conn.send(
      JSON.stringify({
        type,
        data: payload
      })
    )
  }

  conn.onopen = (evt) => {
    send('id', {
      id: id.getLocalId()
    })
    active.value = true
  }

  conn.onclose = function (evt) {
    active.value = false

    errorStore.add({
      type: 'danger',
      message: 'No active websocket connection, please reload the page',
      noexpire: true
    })
    errorStore.add({
      type: 'warning',
      message: 'Test second',
      noexpire: true
    })
  }

  conn.onmessage = function (evt) {
    const msg = JSON.parse(evt.data)
    console.log(msg)

    switch (msg.type) {
      case 'id':
        id.handlePacket(msg.data)
        break
      case 'name':
        name.handlePacket(msg.data)
        break
      default:
        console.error('Unknown message type')
        break
    }
  }

  return { conn, send, active }
})
