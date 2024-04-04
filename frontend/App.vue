<script setup lang="ts">
import { RouterView } from 'vue-router'
import NavBar from '@/components/NavBar.vue'

import { useSocketStore } from '@/stores/socket'
import LoadingView from '@/views/LoadingView.vue'
import ErrorStore from './components/ErrorStore.vue'

const socket = useSocketStore()
</script>

<template>
  <div v-if="socket.userReady">
    <NavBar />
    <ErrorStore />
    <RouterView />
  </div>

  <LoadingView v-else />
</template>
