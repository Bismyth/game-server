import api from '@/api'
import { reactive } from 'vue'
import { stateSchema } from './state'
import type { PrivateState, PublicState } from './state'

const create = (id: string) => {
  api.game.handleAction.fn = handleAction
  api.game.handleEvent.fn = handleEvent
  api.game.handleState.fn = handleState

  ready(id)
}

const data = reactive<{
  public: PublicState | undefined
  private: PrivateState | undefined
}>({
  public: undefined,
  private: undefined,
})

const handleAction = (data: unknown) => {
  console.log(data)
}

const handleState = (i: unknown) => {
  const result = stateSchema.safeParse(i)
  if (!result.success) {
    //todo: better error
    console.log(result.error.format())
    console.error('bad data')
    return
  }

  if (result.data.public) {
    data.public = result.data.public
  }

  if (result.data.private) {
    data.private = result.data.private
  }
}

const ready = (id: string) => {
  api.game.ready(id)
}

const handleEvent = (data: unknown) => {
  console.log(data)
}

export default { create, data, ready }
