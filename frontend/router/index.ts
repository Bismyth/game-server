import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/views/HomeView.vue'
import LobbyView from '@/views/LobbyView.vue'
import api from '@/api'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/lobby/:id',
      name: 'lobby',
      component: LobbyView,
      beforeEnter: (to, from) => {
        if (!api.validateUUID(to.params.id.toString())) {
          return { name: 'home' }
        }
      },
    },
  ],
})

export default router
