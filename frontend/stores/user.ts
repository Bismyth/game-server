import { ref } from 'vue'
import { defineStore } from 'pinia'

import { gameTypes } from '@/game'
import { z } from 'zod'

export const userSchema = z.object({
  id: z.string().uuid(),
  name: z.string(),
  lobbyId: z.string().uuid(),
  gameId: z.string().uuid(),
  gameType: z.union([z.enum(gameTypes), z.literal('')]),
})

export type UserData = z.infer<typeof userSchema>

const nilUUID = '00000000-0000-0000-0000-000000000000'

export const useUserStore = defineStore('user', () => {
  //todo: proxy that errors if not set?
  const data = ref<UserData>({
    id: nilUUID,
    name: 'INVALID',
    lobbyId: nilUUID,
    gameId: nilUUID,
    gameType: '',
  })

  return { data }
})
