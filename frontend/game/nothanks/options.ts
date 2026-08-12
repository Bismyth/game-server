import { z } from 'zod'

export const optionsSchema = z.object({
  rounds: z.number().int().min(1).max(5),
})

export type Options = z.infer<typeof optionsSchema>
