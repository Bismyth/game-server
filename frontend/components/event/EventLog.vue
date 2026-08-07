<script setup lang="ts">
import { nextTick, onMounted, onUnmounted, ref } from 'vue'
import BBCodeText from './BBCodeText.vue'
import apiUtils from '@/api-utils/index.ts'
import IconButton from '../IconButton.vue'
import { Icon } from '@iconify/vue'

const MAX_HISTORY_LENGTH = 100

const eventHistory = ref<string[]>([])

const showBody = ref(true)

const chatContainer = ref<HTMLDivElement | null>(null)

onMounted(() => {
  apiUtils.game.handleEvent.fn = async (v: string) => {
    eventHistory.value.push(v)
    if (eventHistory.value.length > MAX_HISTORY_LENGTH) {
      const size = eventHistory.value.length
      eventHistory.value.splice(0, size - MAX_HISTORY_LENGTH)
    }
    await nextTick()
    if (chatContainer.value) {
      chatContainer.value.scrollTop = chatContainer.value.scrollHeight
    }
  }
})

const toggleEvents = async () => {
  showBody.value = !showBody.value
  await nextTick()
  if (showBody.value && chatContainer.value) {
    console.log('trying to set scroll')
    chatContainer.value.scrollTop = chatContainer.value.scrollHeight
  }
}

onUnmounted(() => {
  apiUtils.game.handleEvent.fn = undefined
})
</script>

<template>
  <div class="log-window panel is-primary">
    <div class="panel-heading log-header" @click="toggleEvents">
      <span class="mr-4">Events</span>
      <button class="log-button">
        <Icon :icon="`mdi:arrow-${showBody ? 'up' : 'down'}`" />
      </button>
    </div>
    <div class="panel-block log-messages" v-show="showBody" ref="chatContainer">
      <span v-if="eventHistory.length == 0">No messages yet...</span>
      <BBCodeText v-else v-for="(e, i) in eventHistory" :key="i" :text="e" />
    </div>
  </div>
</template>

<style>
.log-window {
  position: fixed;
  bottom: 16px;
  right: 16px;
  z-index: 999;
  overflow: hidden;
  max-width: 400px;
}
.log-messages {
  max-height: 200px;
  min-width: 300px;
  overflow-y: scroll;
}
.log-button {
  margin-left: auto;
  background-color: rgba(0, 0, 0, 0.5);
  align-self: center;
  border-radius: 9999px;
  flex-grow: 0;
  flex-shrink: 0;
  font-size: 1.25rem;
  width: 1.25rem;
  height: 1.25rem;
}

.log-header {
  padding: 0.5em 1em;
  display: flex;
  align-items: center;
}
</style>
