import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/views/HomeView.vue'
import LobbyView from '@/views/LobbyView.vue'
import GameView from '@/views/GameView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/room/:id',
      name: 'room',
      component: LobbyView,
    },
    {
      path: '/game/:id',
      name: 'game',
      component: GameView,
    },
  ],
})

export default router
