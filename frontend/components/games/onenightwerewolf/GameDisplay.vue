<script setup lang="ts">
import game from '@/game/onenightwerewolf'
import { useLobbyStore } from '@/stores/lobby'
import { watch, onMounted, defineAsyncComponent } from 'vue'
import { roleActions, type Role } from '@/game/onenightwerewolf/roles'
const lobby = useLobbyStore()

onMounted(() => {
  game.create(lobby.id)
})

const getRoleActionC = (role: Role | undefined) => {
  return defineAsyncComponent({
    // the loader function
    loader: () => {
      if (!role) {
        return import('./RoleLoading.vue')
      }

      if (!roleActions.includes(role)) {
        return import('./NoRoleAction.vue')
      }

      return import(`./roles/${role}/GameDisplay.vue`)
    },
    timeout: 3000,
  })
}
let RoleAction = getRoleActionC(undefined)

watch(
  () => game.data.private,
  (nv) => {
    if (nv) {
      RoleAction = getRoleActionC(nv?.role)
    }
  },
)
</script>

<template>
  <h1>Game</h1>
  <pre>{{ game.data.public }}</pre>
  <pre>{{ game.data.private }}</pre>
  <RoleAction />
</template>
