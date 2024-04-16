<script setup lang="ts" generic="T extends ZodTypeAny">
import { computed, ref } from 'vue'
import { reactive } from 'vue'
import type { z, ZodTypeAny } from 'zod'
import { fromZodError, ValidationError } from 'zod-validation-error'
import ErrorMessages from './ErrorMessages.vue'
import type { ErrorMessage } from '@/stores/error'

type FormData = Partial<FormDataResult>
type FormDataResult = z.infer<typeof props.schema>

interface Props {
  schema: T
}

const props = defineProps<Props>()
// const optionalSchema = props.schema.partial()

const data = reactive<FormData>({})
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

  const result = props.schema.safeParse(data)
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
  for (const k in initData) {
    let value = initData[k]
    if (value) {
      data[k] = value
    }
  }
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
