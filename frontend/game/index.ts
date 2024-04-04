export const gameTypes = ['liarsdice', 'debug'] as const

export type GameTypes = (typeof gameTypes)[number]

export const gameTypeLabels: { [k in GameTypes]: string } = {
  liarsdice: 'Liars Dice',
  debug: 'Debug',
}
