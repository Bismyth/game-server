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
  <div class="dropdown is-hoverable is-right">
    <div class="dropdown-trigger">
      <button class="button" aria-haspopup="true" aria-controls="dropdown-menu">
        <Icon :icon="`fa6-solid:${darkMode === null ? 'desktop' : darkMode ? 'moon' : 'sun'}`" />
        <Icon class="icon is-small" icon="fa6-solid:angle-down" />
      </button>
    </div>
    <div class="dropdown-menu" id="dropdown-menu" role="menu">
      <div class="dropdown-content">
        <a
          class="dropdown-item"
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
          class="dropdown-item"
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
          class="dropdown-item"
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
  </div>
</template>
