<script setup lang="ts">
import ErrorStore from '@/components/ErrorStore.vue'
import JoinMenu from '@/components/JoinMenu.vue'
import SlideOut from '@/components/SlideOut.vue'
</script>

<template>
  <main class="centerize">
    <ErrorStore />
    <JoinMenu :id="($route.query.id ?? '').toString()" />
  </main>
  <SlideOut />
</template>
