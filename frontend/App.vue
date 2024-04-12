<script setup lang="ts">
import { RouterView } from 'vue-router'

import { useSocketStore } from '@/stores/socket'
import LoadingView from '@/views/LoadingView.vue'

const socket = useSocketStore()
</script>

<template>
  <div v-if="socket.userReady">
    <RouterView />
  </div>

  <LoadingView v-else />
</template>
