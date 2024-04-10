import { setSendMessage, validateUUID, isNilUUID } from './main'
import { handleIncomingMessage } from './routes'
import user from './user'
import lobby from './lobby'
import game from './game'

export default {
  handleIncomingMessage,
  setSendMessage,
  user,
  lobby,
  game,
  validateUUID,
  isNilUUID,
}
