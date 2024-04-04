<script setup lang="ts">
import { useLobbyStore } from '@/stores/lobby'
import { z } from 'zod'
import FormWrap from '@/components/FormWrap.vue'
import api from '@/api'

const lobby = useLobbyStore()

const optionsSchema = z.object({
  startingDice: z.number().int().min(1).max(99),
})

type FormResult = z.infer<typeof optionsSchema>
const submit = (data: FormResult) => {
  api.lobby.options(lobby.id, data)
}
</script>

<template>
  <h1>Liars Dice Options</h1>
  <FormWrap :schema="optionsSchema" @submit="submit">
    <template #default="context">
      <div class="field is-horizontal">
        <div class="field-label is-normal">
          <label class="label">Starting Dice</label>
        </div>
        <div class="field-body">
          <div class="field">
            <div class="control">
              <input class="input" v-model="context.data.startingDice" type="number" />
            </div>
          </div>
        </div>
      </div>
      <div class="field is-horizontal">
        <div class="field-label">
          <!-- Left empty for spacing -->
        </div>
        <div class="field-body">
          <div class="field">
            <div class="control">
              <button class="button is-link" type="submit">Send</button>
            </div>
          </div>
        </div>
      </div>
    </template>
  </FormWrap>
</template>
