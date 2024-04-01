import { useErrorStore } from '@/stores/error'
import { z, type ZodTypeAny } from 'zod'
import type { OPacketType } from './packetTypes'

export const sendMessage = (payload: { type: OPacketType; data: any }) => {
  internalSend(payload)
}

// eslint-disable-next-line @typescript-eslint/no-unused-vars
let internalSend = (payload: any) => {
  const error = useErrorStore()
  error.add({
    message: 'Send function not initilised.',
    type: 'danger',
  })
}

const uuidTest = new RegExp(
  /^[0-9a-f]{8}-[0-9a-f]{4}-[0-7][0-9a-f]{3}-[089ab][0-9a-f]{3}-[0-9a-f]{12}$/,
)

export const validateUUID = (id: string) => uuidTest.test(id)

export const setSendMessage = (sendFn: typeof sendMessage) => {
  internalSend = sendFn
}

export const parseData = <T extends ZodTypeAny>(data: unknown, schema: T): z.infer<T> => {
  const result = schema.safeParse(data)
  if (!result.success) {
    // some error
    throw Error('invalid data packet')
  }

  return result.data
}
