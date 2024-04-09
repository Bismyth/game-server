<script setup lang="ts">
import { ref } from 'vue'

const faces = [2, 3, 4, 5, 6]

const props = defineProps<{ currentBid: string }>()

const emit = defineEmits<{
  bid: [string]
}>()

const formSubmitHandler = (e: Event) => {
  e.preventDefault()

  const target = e.target as HTMLFormElement

  const formData = new FormData(target)

  const formValues = [...formData.values()]

  if (formValues.length !== 2) {
    errorMessage.value = 'Please enter all fields'
    return
  }

  const a = parseInt(formValues[0].toString())
  const f = parseInt(formValues[1].toString())

  if (Number.isNaN(a) || Number.isNaN(f)) {
    errorMessage.value = 'Please enter all fields'
    return
  }

  if (props.currentBid !== '') {
    const [ca, cf] = props.currentBid.split(',').map((v) => parseInt(v))

    if (a < ca) {
      errorMessage.value = 'The amount you have bid is too low'
      return
    }

    if (a == ca && f <= cf) {
      errorMessage.value = 'Raise the amount or the face (or both)'
      return
    }
  }

  emit('bid', `${a},${f}`)

  target.reset()
}

const handleReset = (e: Event) => {
  e.preventDefault()
  errorMessage.value = ''
}

const errorMessage = ref('')
</script>

<template>
  <form @submit="formSubmitHandler" @reset="handleReset">
    <input type="submit" hidden />
    <div class="field">
      <div class="field has-addons">
        <p class="control">
          <input class="input" type="text" placeholder="Amount..." name="amount" />
        </p>
        <p class="control">
          <span class="select">
            <select name="face">
              <option disabled selected>Face...</option>
              <option v-for="o in faces" :key="o" :value="o">{{ o }}</option>
            </select>
          </span>
        </p>
        <p class="control">
          <button class="button is-primary" type="submit">Bid</button>
        </p>
      </div>
      <p class="help is-danger">{{ errorMessage }}</p>
    </div>
  </form>
</template>
