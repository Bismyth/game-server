<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const lobbyId = computed(() => route.params.id.toString())
</script>

<template>
  <main class="centerize">
    <JoinMenu :id="lobbyId" />
  </main>
</template>
