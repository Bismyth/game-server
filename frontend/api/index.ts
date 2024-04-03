import { setSendMessage, validateUUID, isNilUUID } from './main'
import { handleIncomingMessage } from './routes'
import user from './user'
import lobby from './lobby'

export default {
  handleIncomingMessage,
  setSendMessage,
  user,
  lobby,
  validateUUID,
  isNilUUID,
}
