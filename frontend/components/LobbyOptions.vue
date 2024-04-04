<script setup lang="ts">
import api from '@/api'
import { lobbyDataSchema } from '@/api/lobby'
import { gameTypeLabels } from '@/game'
import { useLobbyStore } from '@/stores/lobby'
import { createZodPlugin } from '@formkit/zod'

const lobby = useLobbyStore()

const [zodPlugin, submitHandler] = createZodPlugin(lobbyDataSchema, async (formData) => {
  api.lobby.change(formData)
})
</script>

<template>
  <h1>Lobby Options</h1>
  <FormKit type="form" :plugins="[zodPlugin]" @submit="submitHandler">
    <FormKit type="hidden" name="id" v-model="lobby.id" />
    <FormKit
      type="select"
      :options="gameTypeLabels"
      placeholder="Select Game Type"
      name="gameType"
    />
  </FormKit>
</template>
