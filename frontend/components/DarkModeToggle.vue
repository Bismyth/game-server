<script setup lang="ts">
import { Icon } from '@iconify/vue'
import { ref } from 'vue'

const darkMode = ref<boolean | null>(null)

const changeMode = (e: Event) => {
  const target = e.target as HTMLAnchorElement

  const attr = target.attributes.getNamedItem('mode')
  if (attr === null) {
    return
  }

  switch (attr.nodeValue) {
    case 'dark':
      document.documentElement.setAttribute('data-theme', 'dark')
      darkMode.value = true
      break
    case 'light':
      document.documentElement.setAttribute('data-theme', 'light')
      darkMode.value = false
      break
    default:
      document.documentElement.setAttribute('data-theme', '')
      darkMode.value = null
      break
  }
}
</script>

<template>
  <div class="navbar-item has-dropdown is-hoverable">
    <a class="navbar-link">
      <Icon :icon="`fa6-solid:${darkMode === null ? 'desktop' : darkMode ? 'moon' : 'sun'}`" />
    </a>

    <div class="navbar-dropdown">
      <a
        class="navbar-item"
        :class="{ 'is-selected': darkMode === true }"
        @click="changeMode"
        mode="dark"
      >
        <span class="icon-text is-small no-mouse">
          <span class="icon">
            <Icon icon="fa6-solid:moon" />
          </span>
          <span>Dark</span>
        </span>
      </a>
      <a
        class="navbar-item"
        :class="{ 'is-selected': darkMode === false }"
        @click="changeMode"
        mode="light"
      >
        <span class="icon-text is-small no-mouse">
          <span class="icon">
            <Icon icon="fa6-solid:sun" />
          </span>
          <span>Light</span>
        </span>
      </a>
      <a
        class="navbar-item"
        :class="{ 'is-selected': darkMode === null }"
        @click="changeMode"
        mode="system"
      >
        <span class="icon-text is-small no-mouse">
          <span class="icon">
            <Icon icon="fa6-solid:desktop" />
          </span>
          <span>System</span>
        </span>
      </a>
    </div>
  </div>
</template>
