import Theme from '@primeuix/themes/nora'
import { createPinia } from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'
import PrimeVue from 'primevue/config'
import { createApp } from 'vue'
import App from './App.vue'
import router from './router/index.js'
import 'virtual:uno.css'

const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)

createApp(App)
    .use(PrimeVue, {
        theme: {
            preset: Theme,
        },
        options: {
            darkModeSelector: '.my-app-dark',
            cssLayer: false,
        },
    })
    .use(pinia)
    .use(router)
    .mount('#app')
