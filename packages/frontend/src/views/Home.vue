<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { type ApplicationApplicationSummary } from '../api/Api'
import { client, goToApp } from '../utils/client'
import GenericTile, { TileStatus } from '../components/GenericTile.vue'


const applications = ref<ApplicationApplicationSummary[]>([])
const isLoading = ref(true)
const errorMessage = ref('')

const FilterOptions = ['all', 'outdated', 'standard'] as const
let filterChoice = ref<number>(2)

async function fetchApps() {
    isLoading.value = true
    client.api.v1AppsList({
        filter: FilterOptions[filterChoice.value]
    })
        .then(res => applications.value = res.data)
        .catch(reason => errorMessage.value = reason)
        .finally(() => isLoading.value = false)
}

onMounted(fetchApps)
</script>

<template>
    <div class="button-zone">
        <div>
            Items: {{ Object.keys(applications).length }}
        </div>
        <button type="button"
            @click="filterChoice = filterChoice + 1 > FilterOptions.length - 1 ? 0 : filterChoice + 1; fetchApps()">{{
                FilterOptions[filterChoice] ?? filterChoice }}</button>
        <button type="button" @click="fetchApps()">Refresh</button>
    </div>
    <p v-if="isLoading">Loading...</p>
    <p v-else-if="errorMessage">errorMessage</p>
    <div v-else class="appList">
        <GenericTile v-for="application in applications" class="tile" @click="goToApp(application)"
            :status="(application.status as TileStatus)">
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

.button-zone {
    right: 0;
    width: max-content;
    display: flex;
    flex-direction: row;
    gap: 1rem;
}

.button-zone button {
    background-color: white;
    border: solid 1px 1px;
    width: 6rem;
    padding: 0.5rem;
}

table {
    border-spacing: 0;
}

td {
    padding-right: 0.5rem;
}
</style>