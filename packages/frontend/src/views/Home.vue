<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { type ParserApplicationSummary } from '../api/Api'
import ApplicationTile from '../components/ApplicationTile.vue'
import { client } from '../utils/client'

const health = ref<string>('')

const applications = ref<Record<string, ParserApplicationSummary>>({})

onMounted(async () => {
    try {
        health.value = (await client.api.v1HealthList()).data.status ?? 'Inconnu'
        if (health.value == "ok") {
            applications.value = (await client.api.v1AppsList()).data
        }
    } catch (error) {
        console.error('Error fetching health status:', error)
    }
})
</script>

<template>
    <h1>Vue + Gin Monorepo</h1>
    <div class="card">
        <p>Backend Health Status: {{ health }}</p>
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
    gap: 2rem;
}
</style>