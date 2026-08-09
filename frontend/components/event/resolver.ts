
import { h } from 'vue'
import { useRoomStore } from '@/stores/room' 
import { Icon } from '@iconify/vue'


type Resolver = (innerText: string) => any

export const RESOLVER_MAP: Record<string, Resolver> = {
  player: (id) => {
    const room = useRoomStore()
    return h('strong', { class: 'player-name' }, room.users.names[id])
  },
  icon: (name) => {
    return h(Icon, {icon: name})
  }
}
