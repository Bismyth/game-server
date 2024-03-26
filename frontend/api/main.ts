type OutgoingMessage = any // ! TODO: not any

let sendMessage = (message: OutgoingMessage) => {
  // Todo: Throw error for unintilised sendMessage func
}

const setSendMessage = (sendFn: typeof sendMessage) => {
  sendMessage = sendFn
}

const handleMessage = (message: string) => {}

export default { setSendMessage, handleMessage }
