import { z } from 'zod'
import { parseData } from './main'
import { useErrorStore } from '@/stores/error'

const errorMessageSchema = z.string()

export const errorHandle = (data: unknown) => {
  const parsedData = parseData(data, errorMessageSchema)

  const errorStore = useErrorStore()
  errorStore.add({
    type: 'danger',
    message: `Server Error: ${parsedData}`,
  })
}

export const kickMessage = () => {
  const errorStore = useErrorStore()
  errorStore.add({
    type: 'danger',
    message: `You were kicked from the room`,
  })
}
