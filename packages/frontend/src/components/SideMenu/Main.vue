<script setup lang="ts">
import { onMounted } from 'vue'
import logoUrl from '../../assets/logo.png'
import { Button } from "primevue";
import GenericDrawer from "./GenericDrawer.vue";
import { useSideMenuStore } from '../../stores/sideMenu';
import router from '../../router';

const sideMenuStore = useSideMenuStore()
onMounted(async () => {
    sideMenuStore.checkHealth()
})

function toggleDarkMode() {
    document.documentElement.classList.toggle('my-app-dark');
}

</script>

<template>
    <div id="side-menu" :class="{ menu: true, condensed: sideMenuStore.condensed }">
        <div>
            <div class="header" :class="{ menu: true, condensed: sideMenuStore.condensed }">
                <img :src="logoUrl" alt="" srcset="">
                <span v-if="!sideMenuStore.condensed">ArgoCD Source Tracker</span>
                <Button link class="expand-button" @click.prevent="sideMenuStore.condensed = !sideMenuStore.condensed">{{
                    !sideMenuStore.condensed ? "<" : ">" }}</Button>
            </div>
        </div>
        <GenericDrawer>
            <Button icon="pi pi-table" link label="Applications" @click="router.push({ name: 'Home' })"></Button>
        </GenericDrawer>
        <router-view name="menu"></router-view>
        <div class="separator"></div>
        <GenericDrawer>
            Api status: {{ sideMenuStore.health }}
        </GenericDrawer>
        <GenericDrawer>
            <Button variant="link" as="a" href="/ui/docs">Documentation</Button>
        </GenericDrawer>
        <GenericDrawer>
            <Button variant="outlined" label="Toggle Dark Mode" @click="toggleDarkMode()" />
        </GenericDrawer>
    </div>
</template>

<style scoped>
#side-menu {
    padding: 1rem 0;
    display: flex;
    flex-direction: column;
    width: 14rem;
}

#side-menu.condensed {
    width: 4rem;
}

img {
    width: 4rem;
    height: 4rem;
}

.separator {
    flex: 1;
}

.header {
    display: flex;
    align-content: center;
    justify-items: center;
    align-items: center;
}

.header.condensed {
    flex-direction: column;
}

.expand-button {
    width: 3rem;
}
</style>