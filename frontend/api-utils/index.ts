import { setSendMessage, validateUUID, isNilUUID } from './main'
import { handleIncomingMessage } from './routes'
import room from './room'
import game from './game'
import http from './http'

export default {
  handleIncomingMessage,
  setSendMessage,
  room,
  game,
  http,
  validateUUID,
  isNilUUID,
}
