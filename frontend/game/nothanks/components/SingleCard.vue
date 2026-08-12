<script setup lang="ts">
const props = defineProps<{
  value: number
  style?: { boxShadow: string }
  clickable?: boolean
}>()
</script>

<template>
  <div class="card b-primary" :class="{ click: clickable }" :style="style">
    <span>{{ value }}</span>
  </div>
</template>

<style scoped>
.card {
  width: 2.5rem;
  justify-content: center;
  height: 3.5rem;
  align-items: center;
  display: flex;
  border-radius: 5px !important;
  box-shadow: none;
  font-size: 1.5em;
  -webkit-user-select: none;
  user-select: none;
}

.click {
  cursor: pointer;
}
</style>
