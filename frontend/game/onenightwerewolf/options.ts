import { z } from 'zod'
import { roles } from './roles'

export const optionsSchema = z.object({
  roles: z.array(z.enum(roles)),
})

export type Options = z.infer<typeof optionsSchema>
