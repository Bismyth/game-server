import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useUserStore = defineStore('user', () => {
  const data = ref<{
    id?: string
    name?: string
  } | undefined>(undefined)

  return { data }
})
