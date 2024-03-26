<script setup lang="ts">
import { RouterView } from 'vue-router'
import NavBar from '@/components/NavBar.vue'
import ErrorMessage from '@/components/ErrorMessage.vue'

import { useSocketStore } from '@/stores/socket'

useSocketStore()
</script>

<template>
  <NavBar />
  <ErrorMessage />
  <RouterView />
</template>
