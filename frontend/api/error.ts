import { z } from 'zod'
import { parseData } from './main'
import { useErrorStore } from '@/stores/error'

const errorMessageSchema = z.string()

const handle = (data: unknown) => {
  const parsedData = parseData(data, errorMessageSchema)

  const errorStore = useErrorStore()
  errorStore.add({
    type: 'danger',
    message: parsedData,
  })
}

export default { handle }
