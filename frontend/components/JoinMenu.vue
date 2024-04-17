<script setup lang="ts">
import api from '@/api'
import DarkModeToggle from '@/components/DarkModeToggle.vue'
import FullLogo from '@/components/FullLogo.vue'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps<{
  id?: string
}>()

const join = () => {
  if (name.value !== '') {
    localStorage.setItem('nickname', name.value)
  }

  const roomName = name.value || previousName.value

  if (props.id) {
    api.http.joinRoom(props.id, roomName)
  } else {
    api.http.createRoom(roomName)
  }
}

const name = ref('')

const previousName = ref('')
const router = useRouter()

onMounted(() => {
  if (props.id) {
    if (localStorage.getItem(`room:${props.id}`)) {
      router.replace({ name: 'room', params: { id: props.id } })
    }
  }

  api.http.validateTokens()

  const nick = localStorage.getItem('nickname')
  if (nick) {
    previousName.value = nick
  }
})
</script>

<template>
  <div class="box container b-primary wrapper">
    <div class="logo-header mb-6">
      <div class="outer"></div>
      <FullLogo :width="500" />
      <div class="outer">
        <DarkModeToggle />
      </div>
    </div>
    <div class="name">
      <div class="box mb-3 b-primary">
        <h2 class="title is-4 has-text-primary">Choose your Nickname...</h2>
        <form>
          <div class="field">
            <p class="control">
              <input
                class="input"
                type="text"
                :placeholder="previousName || `Enter Nickname...`"
                v-model="name"
              />
            </p>
          </div>
          <hr />
          <div class="field is-grouped is-grouped-centered">
            <div class="control">
              <button class="button is-primary" @click.prevent="join">
                {{ id ? 'Join Lobby' : 'Create Lobby' }}
              </button>
            </div>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
