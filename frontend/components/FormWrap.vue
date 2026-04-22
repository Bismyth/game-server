<script setup lang="ts" generic="T extends ZodObject">
import { computed, ref } from 'vue'
import type { z, ZodObject } from 'zod'
import { fromZodError, ValidationError } from 'zod-validation-error'
import ErrorMessages from './ErrorMessages.vue'
import type { ErrorMessage } from '@/stores/error'

type FormData = Partial<FormDataResult>
type FormDataResult = z.infer<typeof props.schema>

const props = defineProps<{
  schema: T
}>()
// const optionalSchema = props.schema.partial()

const data = ref<FormData>({})
const errors = ref<ValidationError | undefined>()

const errorMessages = computed(() => {
  if (!errors.value) {
    return []
  }

  const message: ErrorMessage = {
    type: 'danger',
    message: errors.value.message,
  }

  return [message]
})

const removeErrors = () => {
  errors.value = undefined
}

const emit = defineEmits<{
  submit: [data: FormDataResult]
}>()

const handleSubmit = (e: Event) => {
  e.preventDefault()
  const result = props.schema.safeParse(data.value)
  if (!result.success) {
    errors.value = fromZodError(result.error)
    return
  }

  emit('submit', result.data)
}

const init = (initData: FormData) => {
  if (initData === undefined) {
    return
  }
  data.value = { ...initData }
}

const formRef = ref<HTMLFormElement | null>(null)

const submit = () => {
  formRef.value?.requestSubmit()
}

defineExpose({
  init,
  submit,
})
</script>

<template>
  <form @submit="handleSubmit" ref="formRef">
    <slot name="errors" :errors="errors">
      <ErrorMessages :messages="errorMessages" @delete="removeErrors" />
    </slot>

    <slot :data="data"></slot>
  </form>
</template>
