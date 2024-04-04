<script setup lang="ts">
import api from '@/api'
import { useLobbyStore } from '@/stores/lobby'
import { createZodPlugin } from '@formkit/zod'
import { z } from 'zod'

const lobby = useLobbyStore()

const optionsSchema = z.object({
  startingDice: z.number().int().min(1).max(99),
})

const [zodPlugin, submitHandler] = createZodPlugin(optionsSchema, async (formData) => {
  api.lobby.options(lobby.id, formData)
})
</script>

<template>
  <h1>Lobby Options</h1>
  <FormKit type="form" :plugins="[zodPlugin]" @submit="submitHandler">
    <FormKit
      type="number"
      placeholder="Select Starting Dice"
      name="startingDice"
      number="integer"
    />
  </FormKit>
</template>
