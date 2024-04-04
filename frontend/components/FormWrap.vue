<script setup lang="ts" generic="T extends ZodTypeAny">
import { onMounted } from 'vue'
import { reactive } from 'vue'
import type { z, ZodTypeAny } from 'zod'

type FormData = Partial<FormDataResult>
type FormDataResult = z.infer<typeof props.schema>

interface Props {
  schema: T
  init?: () => FormData
}

const props = defineProps<Props>()
// const optionalSchema = props.schema.partial()

const data = reactive<FormData>({})
const errors = reactive({})

const emit = defineEmits<{
  submit: [data: FormDataResult]
}>()

const handleSubmit = (e: Event) => {
  e.preventDefault()

  const result = props.schema.safeParse(data)
  if (!result.success) {
    // handle error
    console.log(result.error)
    return
  }

  emit('submit', result.data)
}

onMounted(() => {
  const initData = props.init?.()
  if (initData === undefined) {
    return
  }
  for (const k in initData) {
    let value = initData[k]
    if (value) {
      data[k] = value
    }
  }
})
</script>

<template>
  <form @submit="handleSubmit">
    <slot :data="data" :errors="errors"></slot>
  </form>
</template>
