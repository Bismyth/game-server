import '@/assets/main.scss'
import '@formkit/themes/genesis'
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { plugin, defaultConfig } from '@formkit/vue'
import App from './App.vue'
import router from './router'

const app = createApp(App)
const pinia = createPinia()
app.use(pinia)
app.use(router)
app.use(plugin, defaultConfig)

app.mount('#app')
