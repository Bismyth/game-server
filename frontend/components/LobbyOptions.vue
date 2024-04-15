<script setup lang="ts">
import { roomDataSchema } from '@/api/room'
import { gameTypeLabels } from '@/game'
import { useRoomStore } from '@/stores/room'
import FormWrap from './FormWrap.vue'
import ModalWrap from './ModalWrap.vue'
import type { z } from 'zod'
import api from '@/api'
import { ref, watch } from 'vue'
import type { ComponentExposed } from 'vue-component-type-helpers'
import IconButton from './IconButton.vue'

const room = useRoomStore()

type FormResult = z.infer<typeof roomDataSchema>

const submit = (data: FormResult) => {
  api.room.change(data)
  showOptions.value = false
}

const showOptions = ref(false)
const onClose = () => {
  showOptions.value = false
}

const openModal = () => {
  showOptions.value = true
}

const formWrap = ref<ComponentExposed<typeof FormWrap<typeof roomDataSchema>> | null>(null)

watch(showOptions, (newValue) => {
  if (newValue) {
    formWrap.value?.init({
      gameType: room.gameType,
    })
  }
})
</script>

<template>
  <IconButton @click="openModal" icon="fa6-solid:pencil" label="Edit" />
  <ModalWrap :shown="showOptions" title="Lobby Options" @close="onClose">
    <template #body>
      <FormWrap :schema="roomDataSchema" @submit="submit" ref="formWrap">
        <template #default="context">
          <div class="field">
            <label class="label">Game Type</label>
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
        </template>
      </FormWrap>
    </template>
    <template #footer>
      <div class="buttons">
        <button class="button is-link" @click="formWrap?.submit">Submit</button>
        <button class="button" @click="onClose">Cancel</button>
      </div>
    </template>
  </ModalWrap>
</template>
