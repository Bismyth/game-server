import { z } from 'zod'

export const optionsSchema = z.object({
  startingDice: z.number().int().min(1).max(99),
})

export type Options = z.infer<typeof optionsSchema>
