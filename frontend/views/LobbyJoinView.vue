<script setup lang="ts">
import api from '@/api'
import { onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const goHome = () => {
  router.push('/')
}

const route = useRoute()

onMounted(() => {
  api.lobby.join(route.params.id.toString())
})
</script>

<template>
  <span class="mb-2">Failed to join lobby</span>

  <button class="button" @click="goHome">Go Home</button>
</template>
