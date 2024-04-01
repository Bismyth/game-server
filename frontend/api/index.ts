import { setSendMessage, validateUUID } from './main'
import { handleIncomingMessage } from './routes'
import user from './user'
import lobby from './lobby'

export default {
  handleIncomingMessage,
  setSendMessage,
  user,
  lobby,
  validateUUID,
}
