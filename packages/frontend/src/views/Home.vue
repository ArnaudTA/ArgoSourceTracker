<script setup lang="ts">
import { onMounted } from 'vue'
import { goToApp } from '../utils/client'
import GenericTile from '../components/GenericTile.vue'
import { useHomeStore } from '../stores/home'
import { Button } from 'primevue'
import Dashboard from '../components/Dashboard.vue'
import { ApplicationApplicationStatus } from '../api/Api'

const homeStore = useHomeStore()


onMounted(homeStore.fetchApps)
</script>

<template>
    <div class="layout">
        <Dashboard class="dashboard">
            <Button @click="homeStore.fetchApps()">Refresh</Button>
            <div>
                Items: {{ Object.keys(homeStore.applications).length }}
            </div>
        </Dashboard>
        <div class="content">

            <p v-if="homeStore.isLoading">Loading...</p>
            <p v-else-if="homeStore.errorMessage">errorMessage</p>
            <div v-else class="appList">
                <GenericTile v-for="application in homeStore.applications" class="tile" @click="goToApp(application)"
                    :status="(application.status as ApplicationApplicationStatus)">
                    <table>
                        <tbody>
                            <tr>
                                <td>Name:</td>
                                <td>{{ application.name }}</td>
                            </tr>

                            <tr>
                                <td>Namespace:</td>
                                <td>{{ application.namespace }}</td>
                            </tr>

                            <tr>
                                <td>Status:</td>
                                <td>{{ application.status }}</td>
                            </tr>
                        </tbody>
                    </table>
                </GenericTile>
            </div>
        </div>
    </div>
</template>

<style scoped>
.appList {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
    /* Optionnel : espace entre les éléments */
    gap: 2rem;
}

.tile {
    padding: 0.5rem;
    display: flex;
    flex-direction: row;
}

table {
    border-spacing: 0;
}

td {
    padding-right: 0.5rem;
}

.content {
    padding: 1.5rem;
}

.dashboard {
    display: flex;
    gap: 3rem;
    justify-content: space-between;
}
</style>