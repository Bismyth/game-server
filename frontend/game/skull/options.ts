import { z } from 'zod'

export const optionsSchema = z.object({
  discardRandom: z.boolean(),
})

export type Options = z.infer<typeof optionsSchema>
