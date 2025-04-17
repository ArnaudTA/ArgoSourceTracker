import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import UnoCSS from 'unocss/vite'

export default defineConfig({
  plugins: [vue(), UnoCSS()],
  server: {
    port: 3000,
    host: true,
    proxy: {
      '/api': {
        target: 'http://api:8080/',
      },
      '/metrics': {
        target: 'http://api:8081/',
      },
      '/swagger': {
        target: 'http://api:8080/',
      }
    }
  },
  
})