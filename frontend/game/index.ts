export const gameTypes = ['liarsdice', 'onenightwerewolf'] as const

export type GameTypes = (typeof gameTypes)[number]

export const gameTypeLabels: { [k in GameTypes]: string } = {
  liarsdice: 'Liars Dice',
  onenightwerewolf: 'One Night Ultimate Werewolf',
}
