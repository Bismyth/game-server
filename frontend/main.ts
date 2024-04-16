import '@/assets/main.scss'
import 'floating-vue/dist/style.css'
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import FloatingVue from 'floating-vue'

const app = createApp(App)
const pinia = createPinia()
app.use(FloatingVue)
app.use(pinia)
app.use(router)
app.mount('#app')
