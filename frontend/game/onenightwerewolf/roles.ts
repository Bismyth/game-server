export const roles = [
  'werewolf',
  'villager',
  'robber',
  'drunk',
  'troublemaker',
  'mason',
  'minion',
  'insomniac',
  'seer',
] as const

export type Role = (typeof roles)[number]

type RoleInfo = {
  name: string
}

export const roleDisplay: Record<Role, RoleInfo> = {
  werewolf: {
    name: 'Werewolf',
  },
  villager: {
    name: 'Villager',
  },
  robber: {
    name: 'Robber',
  },
  drunk: {
    name: 'Drunk',
  },
  troublemaker: {
    name: 'Troublemaker',
  },
  mason: {
    name: 'Mason',
  },
  minion: {
    name: 'Minion',
  },
  insomniac: {
    name: 'Insomniac',
  },
  seer: {
    name: 'Seer',
  },
}
