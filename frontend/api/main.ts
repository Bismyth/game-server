import { useErrorStore } from '@/stores/error'
import { z, type ZodTypeAny } from 'zod'

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export let sendMessage = (message: any) => {
  const error = useErrorStore()
  error.add({
    message: 'Send function not initilised.',
    type: 'danger',
  })
}

export const setSendMessage = (sendFn: typeof sendMessage) => {
  sendMessage = sendFn
}

export const parseData = <T extends ZodTypeAny>(data: unknown, schema: T): z.infer<T> => {
  const result = schema.safeParse(data)
  if (!result.success) {
    // some error
    throw Error('invalid data packet')
  }

  return result.data
}
