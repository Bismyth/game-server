import { useUserStore } from '@/stores/user'
import { parseData, sendMessage } from './main'
import { z } from 'zod'
import { OPacketType } from './packetTypes'

const idLSKey = 'id'

export const userMessageSchema = z.object({
  id: z.string().uuid(),
  name: z.string(),
})

type userMessage = z.infer<typeof userMessageSchema>

export const handleUserInit = (data: unknown) => {
  const parsedData = parseData(data, userMessageSchema)

  localStorage.setItem(idLSKey, parsedData.id ?? 'wrong uuid')

  const user = useUserStore()
  user.data = parsedData
}

export const handleUserChange = (data: unknown) => {
  const parsedData = parseData(data, userMessageSchema)

  const user = useUserStore()
  user.data = parsedData
}

const getLocalId = () => {
  return localStorage.getItem(idLSKey)
}

const sendChange = (data: userMessage) => {
  sendMessage({
    type: OPacketType.UserChange,
    data,
  })
}

export default { getLocalId, sendChange }
