import { ref } from 'vue'
import { defineStore } from 'pinia'

type UserData = {
  id?: string
  name?: string
}

export const useUserStore = defineStore('user', () => {
  const data = ref<UserData | undefined>(undefined)

  return { data }
})
