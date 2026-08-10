import { useRoomStore } from '@/stores/room'
import router from '@/router'

export const gameTypes = ['liarsdice', 'skull', 'nothanks'] as const

export type GameTypes = (typeof gameTypes)[number]

interface GameInfo {
  displayName: string
  icon: string
}

export const gameTypeLabels: { [k in GameTypes]: GameInfo } = {
  liarsdice: {displayName: 'Liars Dice', icon: 'mdi:dice-2'},
  skull: {displayName: 'Skull (Beta)', icon: 'mdi:skull'},
  nothanks: {displayName: "No Thanks!", icon: 'mdi:skull'},
}

export const handleLobbyBack = () => {
  const room = useRoomStore()
  router.replace({ name: 'room', params: { id: room.data.id } })
}
