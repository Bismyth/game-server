<script setup lang="ts">
import { ref } from 'vue'

const emit = defineEmits<{
  changeTab: [number]
}>()

const tabIndex = ref(0)
const tabs = ['Dice', 'Value']

const handleTabChange = (e: Event, i: number) => {
  e.preventDefault()
  tabIndex.value = i
  emit('changeTab', i)
}
</script>

<template>
  <div class="tabs mb-1">
    <ul>
      <li v-for="(t, i) in tabs" :key="i" :class="{ 'is-active': i === tabIndex }">
        <a @click="(e) => handleTabChange(e, i)">{{ t }}</a>
      </li>
    </ul>
  </div>
</template>
