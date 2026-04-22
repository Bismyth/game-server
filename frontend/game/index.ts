import { useRoomStore } from '@/stores/room'
import { useRouter } from 'vue-router'

export const gameTypes = ['liarsdice', 'skull'] as const

export type GameTypes = (typeof gameTypes)[number]

export const gameTypeLabels: { [k in GameTypes]: string } = {
  liarsdice: 'Liars Dice',
  skull: 'Skull',
}

export const handleLobbyBack = () => {
  const room = useRoomStore()
  const router = useRouter()
  router.replace({ name: 'room', params: { id: room.data.id } })
}
