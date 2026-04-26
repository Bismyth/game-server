import { useRoomStore } from '@/stores/room'
import { useRouter } from 'vue-router'

export const gameTypes = ['liarsdice', 'skull'] as const

export type GameTypes = (typeof gameTypes)[number]

interface GameInfo {
  displayName: string
  icon: string
}

export const gameTypeLabels: { [k in GameTypes]: GameInfo } = {
  liarsdice: {displayName: 'Liars Dice', icon: 'bi:dice-2-fill'},
  skull: {displayName: 'Skull', icon: 'mdi:skull'},
}

export const handleLobbyBack = () => {
  const room = useRoomStore()
  const router = useRouter()
  router.replace({ name: 'room', params: { id: room.data.id } })
}
