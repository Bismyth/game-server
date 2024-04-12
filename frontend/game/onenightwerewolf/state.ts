import { z } from 'zod'
import { roles } from './roles'

const publicStateSchema = z.object({})
const privateStateScehma = z.object({
  role: z.enum(roles),
})

export type PublicState = z.infer<typeof publicStateSchema>
export type PrivateState = z.infer<typeof privateStateScehma>

export const stateSchema = z.object({
  public: publicStateSchema.nullable(),
  private: privateStateScehma.nullable(),
})
