<script setup lang="ts">
import { lobbyDataSchema } from '@/api/lobby'
import { gameTypeLabels } from '@/game'
import { useLobbyStore } from '@/stores/lobby'
import FormWrap from './FormWrap.vue'
import type { z } from 'zod'
import api from '@/api'

const lobby = useLobbyStore()

type FormResult = z.infer<typeof lobbyDataSchema>

const submit = (data: FormResult) => {
  api.lobby.change(data)
}

const init = () => {
  return { id: lobby.id }
}
</script>

<template>
  <h1>Lobby Options</h1>

  <FormWrap :schema="lobbyDataSchema" :init="init" @submit="submit">
    <template #default="context">
      <div class="field is-horizontal">
        <div class="field-label is-normal">
          <label class="label">Game Type</label>
        </div>
        <div class="field-body">
          <div class="field is-narrow">
            <div class="control">
              <div class="select is-fullwidth">
                <select v-model="context.data.gameType">
                  <option selected disabled>Select Game Type</option>
                  <option v-for="(value, key) in gameTypeLabels" :key="key" :value="key">
                    {{ value }}
                  </option>
                </select>
              </div>
            </div>
          </div>
        </div>
      </div>

      <button type="submit">Submit</button>
    </template>
  </FormWrap>
</template>
