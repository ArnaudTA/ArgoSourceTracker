<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { type ParserApplicationSummary } from '../api/Api'
import ApplicationTile from '../components/ApplicationTile.vue'
import { client } from '../utils/client'


const applications = ref<Record<string, ParserApplicationSummary>>({})

onMounted(async () => {
    try {
        applications.value = (await client.api.v1AppsList()).data
    } catch (error) {
        console.error('Error fetching health status:', error)
    }
})
</script>

<template>
    <div class="card">
    </div>

    <div class="appList">
        <ApplicationTile v-for="[name, application] in Object.entries(applications)" :name="name"
            :application="application" />
    </div>
</template>

<style>
.appList {
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;
    gap: 2rem;
}
</style>