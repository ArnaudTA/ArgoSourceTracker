<script setup lang="ts">
import { goToApp } from '../utils/client'
import GenericTile from '../components/GenericTile.vue'
import { useHomeStore } from '../stores/home'
import { Button, MeterGroup, Paginator, ProgressBar, Toolbar } from 'primevue'
import { TypesApplicationStatus } from '../api/Api'

const homeStore = useHomeStore()

</script>

<template>
    <div class="layout">
        <Toolbar style="padding: 1rem">
            <template #start>
                <Button icon="pi pi-refresh" @click="homeStore.refreshApps()" label="Refresh"></Button>
            </template>
            <template #end>
                Total: {{ homeStore.total }}
            </template>
        </Toolbar>

        <Toolbar style="padding: 1rem">
            <template #center>
                <Paginator :rows="homeStore.limit" :totalRecords="homeStore.total"
                    :rowsPerPageOptions="[5, 10, 25, 100]" @update:last="homeStore.setPage" @page="homeStore.setPage">
                </Paginator>
            </template>
        </Toolbar>
        <div class="content">
            <MeterGroup :value="homeStore.meters" :max="homeStore.applications.length" class="mb-4" />
            <ProgressBar :class="{ invisible: !homeStore.isLoading }" mode="indeterminate" style="height: 6px"
                class="mb-4">
            </ProgressBar>
            <p v-if="homeStore.errorMessage">errorMessage</p>
            <div v-else class="appList">
                <GenericTile v-for="application in homeStore.applications" class="tile" @click="goToApp(application)"
                    :status="(application.status as TypesApplicationStatus)">
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
#pagination {
    width: 100%;
    height: 3rem;
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    margin-bottom: 2rem;
}

#pagination-nav {
    height: 3rem;
    display: flex;
    gap: .5rem;
    flex-direction: row;
    justify-content: space-between;
}

.appList {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(350px, auto));
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