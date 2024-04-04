<script setup lang="ts">
import { RouterView } from 'vue-router'
import NavBar from '@/components/NavBar.vue'
import ErrorMessage from '@/components/ErrorMessage.vue'

import { useSocketStore } from '@/stores/socket'
import LoadingView from '@/views/LoadingView.vue'

const socket = useSocketStore()
</script>

<template>
  <div v-if="socket.userReady">
    <NavBar />
    <ErrorMessage />
    <RouterView />
  </div>

  <LoadingView v-else />
</template>
