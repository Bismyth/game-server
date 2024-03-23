import { defineStore } from 'pinia'
import { ref } from 'vue'

import id from '@/api/id'
import name from '@/api/name'

export const useSocketStore = defineStore('socket', () => {

  const conn = new WebSocket('/ws')

  const active = ref(false)

  const send = (type: string, payload: any) => {
    conn.send(JSON.stringify({
      type,
      data: payload
    }))
  }

  conn.onopen = (evt) => {
    send("id", {
      id: id.getLocalId()
    })
    active.value = true
  }

  conn.onclose = function (evt) {
    active.value = false
  };

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
        console.log("Unknown message type")
        break
    }
  }

 

  return { conn, send, active }
})
