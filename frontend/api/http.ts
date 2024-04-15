import router from '@/router'
import { useErrorStore } from '@/stores/error'
import { z } from 'zod'
import { fromZodError } from 'zod-validation-error'

const roomJoinResponseSchema = z.object({
  status: z.enum(['success', 'error']),
  error: z.string(),
  token: z.string(),
  id: z.string().uuid(),
})

const handleJoinResponse = (response: unknown) => {
  const es = useErrorStore()
  const result = roomJoinResponseSchema.safeParse(response)
  if (!result.success) {
    const zodMessage = fromZodError(result.error)
    es.add({
      type: 'danger',
      message: zodMessage.message,
    })
    return
  }
  if (result.data.status === 'error') {
    return
  }

  localStorage.setItem(`room:${result.data.id}`, result.data.token)

  router.replace({ name: 'room', params: { id: result.data.id } })
}

const createRoom = async (name: string) => {
  const es = useErrorStore()

  const res = await fetch('/api/room/create', {
    method: 'POST',
    body: JSON.stringify({ name }),
  })

  if (res.status !== 200) {
    const responseText = await res.text()
    es.add({
      type: 'danger',
      message: responseText,
    })
    return
  }

  const data = await res.json()
  handleJoinResponse(data)
}

const joinRoom = async (id: string, name: string) => {
  const es = useErrorStore()

  const res = await fetch('/api/room/join', {
    method: 'POST',
    body: JSON.stringify({ name, id }),
  })

  if (res.status !== 200) {
    const responseText = await res.text()
    es.add({
      type: 'danger',
      message: responseText,
    })
    return
  }

  const data = await res.json()
  handleJoinResponse(data)
}

export default { createRoom, joinRoom }
