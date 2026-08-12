<script setup lang="ts">
import FormWrap from '@/components/FormWrap.vue'
import type { ComponentExposed } from 'vue-component-type-helpers'
import { ref } from 'vue'
import { useRoomStore } from '@/stores/room'
import { optionsSchema, type Options } from '@/game/nothanks/options'

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

const room = useRoomStore()

const init = () => {
  formWrap.value?.init(room.data.options)
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
        <label class="label">Total Rounds</label>
        <div class="control">
          <input class="input" v-model="context.data.rounds" type="number" />
        </div>
      </div>
    </template>
  </FormWrap>
</template>
