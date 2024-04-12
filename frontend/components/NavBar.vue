<script setup lang="ts">
import DarkModeToggle from '@/components/DarkModeToggle.vue'
import { useUserStore } from '@/stores/user'
import { RouterLink } from 'vue-router'
import api from '@/api'

const user = useUserStore()

const changeName = () => {
  const username = prompt('Update username:')?.trim()

  if (username) {
    api.user.sendNameChange(username)
  }
}
</script>

<template>
  <div class="container">
    <nav class="navbar" role="navigation" aria-label="main navigation">
      <div class="navbar-brand">
        <RouterLink class="navbar-item" to="/">
          <h1 class="title is-size-4">Game Server</h1>
        </RouterLink>

        <a
          role="button"
          class="navbar-burger"
          aria-label="menu"
          aria-expanded="false"
          data-target="navbarBasicExample"
        >
          <span aria-hidden="true"></span>
          <span aria-hidden="true"></span>
          <span aria-hidden="true"></span>
          <span aria-hidden="true"></span>
        </a>
      </div>

      <div id="navbarBasicExample" class="navbar-menu">
        <div class="navbar-start">
          <a class="navbar-item"> Home </a>
        </div>

        <div class="navbar-end">
          <div class="navbar-item">
            <div class="pr-3">Name: {{ user.data?.name }}</div>
            <a class="button is-link" @click="changeName"> Change </a>
          </div>
          <div class="navbar-item">
            <DarkModeToggle />
          </div>
        </div>
      </div>
    </nav>
  </div>
</template>
