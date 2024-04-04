import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

export type ErrorMessage = {
  message: string
  type: 'warning' | 'danger'
  noexpire?: boolean
}

const ERROR_TIMEOUT = 30 * 1000

export const useErrorStore = defineStore('error', () => {
  const messages = ref<Record<number, ErrorMessage>>({})

  let currentErrorId = 0

  const add = (e: ErrorMessage) => {
    const newId = currentErrorId++

    messages.value[newId] = e

    if (!e.noexpire) {
      setTimeout(() => {
        deleteMessage(newId)
      }, ERROR_TIMEOUT)
    }
  }

  const deleteMessage = (id: number) => {
    delete messages.value[id]
  }

  const messageGetter = computed(() => messages)

  return { add, messageGetter, deleteMessage }
})
