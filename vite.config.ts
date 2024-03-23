import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
  ],
  build: {
    outDir: '.output'
  },
  server: {
    port: 8080,
    proxy: {
      '/ws': 'ws://localhost:8081/'
    }
  },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./frontend', import.meta.url))
    }
  }
})
