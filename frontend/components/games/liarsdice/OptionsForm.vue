<script setup lang="ts">
import { z } from 'zod'
import FormWrap from '@/components/FormWrap.vue'
import type { ComponentExposed } from 'vue-component-type-helpers'
import { ref } from 'vue'
import { useLobbyStore } from '@/stores/lobby'

const optionsSchema = z.object({
  startingDice: z.number().int().min(1).max(99),
})

const emit = defineEmits<{
  submit: [T: FormResult]
}>()

type FormResult = z.infer<typeof optionsSchema>
const submit = (data: FormResult) => {
  emit('submit', data)
}

const formWrap = ref<ComponentExposed<typeof FormWrap<typeof optionsSchema>> | null>(null)

const submitForm = () => {
  formWrap.value?.submit()
}

const lobby = useLobbyStore()

const init = () => {
  formWrap.value?.init(lobby.options)
}

defineExpose({
  submit: submitForm,
  init,
})
</script>

<template>
  <FormWrap :schema="optionsSchema" @submit="submit" ref="formWrap">
    <template #default="context">
      <div class="field">
        <label class="label">Starting Dice</label>
        <div class="control">
          <input class="input" v-model="context.data.startingDice" type="number" />
        </div>
      </div>
    </template>
  </FormWrap>
</template>
