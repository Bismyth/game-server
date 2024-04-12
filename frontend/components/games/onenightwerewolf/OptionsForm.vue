<script setup lang="ts">
import FormWrap from '@/components/FormWrap.vue'
import type { ComponentExposed } from 'vue-component-type-helpers'
import { ref } from 'vue'
import { useLobbyStore } from '@/stores/lobby'
import { optionsSchema, type Options } from '@/game/onenightwerewolf/options'

const emit = defineEmits<{
  submit: [T: Options]
}>()

const submit = (data: Options) => {
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
    <template #default>
      <h1>Options Form</h1>
    </template>
  </FormWrap>
</template>
