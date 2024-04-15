import { CallBackFunc, sendMessage } from './main'
import { OPacketType } from './packetTypes'

const handleAction = new CallBackFunc<unknown>()
export const handleGameAction = (data: unknown) => {
  handleAction.run(data)
}

const handleEvent = new CallBackFunc<unknown>()
export const handleGameEvent = (data: unknown) => {
  handleEvent.run(data)
}

const handleState = new CallBackFunc<unknown>()
export const handleGameState = (data: unknown) => {
  handleState.run(data)
}

const startGame = () => {
  sendMessage({
    type: OPacketType.GameStart,
  })
}

const action = (option: string, data: any) => {
  sendMessage({
    type: OPacketType.GameAction,
    data: {
      option,
      data,
    },
  })
}

const ready = () => {
  sendMessage({
    type: OPacketType.GameReady,
  })
}

export default { startGame, action, ready, handleAction, handleEvent, handleState }
