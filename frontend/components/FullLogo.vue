<script setup lang="ts">
import FullLogo from '@/assets/Logo_full.svg?component'

defineProps<{
  width?: number
}>()
</script>

<template>
  <FullLogo class="f-primary" :style="{ width: `${width ?? 400}px` }" />
</template>
