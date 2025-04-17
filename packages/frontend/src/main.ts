import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router/index.js'
import PrimeVue from 'primevue/config';

createApp(App)
    .use(PrimeVue)
    .use(createPinia())
    .use(router)
    .mount('#app')