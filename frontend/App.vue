<script setup lang="ts">
import { RouterView } from 'vue-router'
import NavBar from '@/components/NavBar.vue'

import { useSocketStore } from '@/stores/socket'

useSocketStore()
</script>

<template>
  <NavBar />

  <RouterView />
</template>
