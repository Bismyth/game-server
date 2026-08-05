/* eslint-disable prefer-const */
import api from '@/api'
import { useRoomStore } from '@/stores/room'
import { computed, reactive, ref } from 'vue'
import { z } from 'zod'

const publicStateSchema = z.object({
  tilesPlaced: z.record(z.uuid(), z.number().int()),
  tilesRevealed: z.record(z.uuid(), z.array(z.boolean())),
  bid: z.number().int(),
  passed: z.array(z.uuid()),
  points: z.record(z.uuid(), z.number().int()),
  flipper: z.uuid(),
  gameOver: z.boolean(),
  playerLeft: z.boolean(),
  turnOrder: z.array(z.string()),
  turn: z.uuid(),
  round: z.number().int(),
})

type publicStateT = z.infer<typeof publicStateSchema>

const privateStateScehma = z.object({
  tiles: z.array(z.boolean()),
  tilesPlaced: z.array(z.boolean()),
})

const stateSchema = z.object({
  public: publicStateSchema.nullable(),
  private: privateStateScehma.nullable(),
})

type privateStateT = z.infer<typeof privateStateScehma>

type fullStateT = z.infer<typeof stateSchema>

const create = () => {
  api.game.handleAction.fn = handleAction
  api.game.handleEvent.fn = handleEvent
  api.game.handleState.fn = handleState

  resetValues()

  ready()
}

const resetValues = () => {
  showGameOver.value = false
  gameData.publicState = undefined
  gameData.privateState = undefined
  gameData.currentOptions = []
}

interface GameData {
  publicState: publicStateT | undefined
  privateState: privateStateT | undefined
  currentOptions: string[]
}

const gameData = reactive<GameData>({
  publicState: undefined,
  privateState: undefined,
  currentOptions: [],
})

const newRoundData = ref<fullStateT | undefined>()

const showCall = ref(false)
const showGameOver = ref(false)

const ready = () => {
  api.game.ready()
}

const takeAction = (option: string, data?: any) => {
  api.game.action(option, data)
}

const place = (tile: boolean) => {
  takeAction('place', { tile })
}

const bid = (amount: number) => {
  takeAction('bid', { bid: amount })
}

const pass = () => {
  takeAction('pass')
}

const flip = (player: string) => {
  takeAction('flip', { player })
}

const handleState = (data: unknown) => {
  const room = useRoomStore()

  const result = stateSchema.safeParse(data)
  if (!result.success) {
    //todo: better error
    console.log(result.error.format())
    console.error('bad data')
    return
  }

  if (result.data.public?.gameOver && !showGameOver.value) {
    showGameOver.value = true
  }
  if (result.data.public) {
    if (gameData.publicState !== undefined && result.data.public.round > 1 && result.data.public.round !== gameData.publicState?.round) {
      newRoundData.value = result.data
    } else {
      gameData.publicState = result.data.public
    }
  }
  if (result.data.private) {
    gameData.privateState = result.data.private
  }

  // Current User Action

}

const userActions = computed(() => {
  const room = useRoomStore()
  const data = {
    showMessage: false,
    message: '',
    showBid: false,
    showPass: false,
  }
  if (gameData.publicState?.gameOver) {
    data.showMessage = true
    data.message = "The Game is Over!"
  }
  else if (hasNextRound.value) {
    if (gameData.publicState?.playerLeft) {
      data.showMessage = true
      data.message = "A Player has left"
    } else {
      let flipped = 0
    for (const p in gameData.publicState?.tilesRevealed) {
      flipped += gameData.publicState?.tilesRevealed[p].length
    }


    let message = ``

    if (gameData.publicState?.flipper == room.data.userId) {
      message += 'You'
    } else {
      message += `${room.users.names[gameData.publicState?.flipper ?? '']}`
    }

    let flippedSkull = false
    for (const p in gameData.publicState?.tilesRevealed) {
      for (const t of gameData.publicState.tilesRevealed[p]) {
        if (t) {
          flippedSkull = true
          break
        }
      }
    }

    if (flippedSkull) {
      message += ' flipped a skull!'
    } else {
      message += ` flipped all ${gameData.publicState?.bid} roses`
    }
    data.showMessage = true
    data.message = message
    }


    

  }
    else if (gameData.publicState?.tilesPlaced[room.data.userId] == 0) {
    data.showMessage = true
    data.message = "Place your inital tile"
  } else if (gameData.publicState?.flipper != '00000000-0000-0000-0000-000000000000') {
    const total = gameData.publicState?.bid
  
    let flipped = 0
    for (const p in gameData.publicState?.tilesRevealed) {
      flipped += gameData.publicState.tilesRevealed[p].length
    }
    data.showMessage = true

    if (gameData.publicState?.flipper == room.data.userId) {
      data.message = `Flip ${total} roses. Don't Flip a Skull! Roses: ${flipped}/${total}`
    } else {
      data.message = `${room.users.names[gameData.publicState?.flipper ?? '']} needs to flip ${total} roses. ${flipped}/${total}`
    }

  } else if (gameData.publicState?.turn == room.data.userId) {
    data.showMessage = true
    data.showBid = true
    if (gameData.publicState.bid == 0) {
      data.message = "Place a tile or start the bidding"
    } else {
      data.message = "Raise the bid or pass"
      data.showPass = true
    }
  }
  return data
})



const nextRound = () => {
  if (newRoundData.value === undefined) {
    return
  }

  if (newRoundData.value.public) {
    gameData.publicState = newRoundData.value.public
  }
  if (newRoundData.value.private) {
    gameData.privateState = newRoundData.value.private
  }

  newRoundData.value = undefined
}

const hasNextRound = computed(() => newRoundData.value !== undefined)

const handleAction = (data: unknown) => {
  const result = z.array(z.string()).safeParse(data)
  if (!result.success) {
    //Todo: better error response
    console.error('bad data')
    return
  }
  gameData.currentOptions = result.data
}

const handleEvent = (data: unknown) => {
  console.log(data)
}

const currentHand = computed(() => {
  const totalTiles = gameData.privateState?.tiles ?? []
  const placedTiles = gameData.privateState?.tilesPlaced ?? []

  const hasSkull = totalTiles.includes(true)
  const placedSkull = placedTiles.includes(true)

  const cHand = []

  if (hasSkull && !placedSkull) {
    cHand.push(true)
  }
  
  const restSize = totalTiles.length - placedTiles.length - cHand.length
  for (let x = 0; x < restSize; x++) {
    cHand.push(false)
  }

  return cHand
})



const canFlip = (id: string): boolean => {
  return (gameData.publicState?.tilesPlaced[id] ?? 0) > (gameData.publicState?.tilesRevealed[id].length ?? 0)
}



export default {
  create,
  ready,
  place,
  bid,
  pass,
  flip,
  gameData,
  showCall,
  showGameOver,
  currentHand,
  nextRound,
  hasNextRound,
  canFlip,
  userActions
}
