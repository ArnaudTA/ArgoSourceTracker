import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router/index.js'
import PrimeVue from 'primevue/config';
import Theme from '@primeuix/themes/nora';
import 'virtual:uno.css'

createApp(App)
    .use(PrimeVue, {
        theme: {
            preset: Theme,
        },
        options: {
            darkModeSelector: '.my-app-dark',
            cssLayer: false
        }
    })
    .use(createPinia())
    .use(router)
    .mount('#app')