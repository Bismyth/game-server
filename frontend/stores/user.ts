import { ref } from 'vue'
import { defineStore } from 'pinia'

import { z } from 'zod'

export const userSchema = z.object({
  id: z.string().uuid(),
  name: z.string(),
  lobbies: z.array(z.string().uuid()).max(1),
  token: z.string().nullable(),
})

export type UserData = z.infer<typeof userSchema>

const nilUUID = '00000000-0000-0000-0000-000000000000'

export const useUserStore = defineStore('user', () => {
  //todo: proxy that errors if not set?
  const data = ref<UserData>({
    id: nilUUID,
    name: 'INVALID',
    lobbies: [],
    token: null,
  })

  return { data }
})
