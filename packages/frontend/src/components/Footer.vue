<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { type ParserApplicationSummary } from '../api/Api'
import { client } from '../utils/client'

const health = ref<string>('Checking...')

const applications = ref<Record<string, ParserApplicationSummary>>({})

onMounted(async () => {
    try {
        health.value = await client.api.v1HealthList()
            .then(res => res.data.status)
            .catch(_ => "Can't reach")
        if (health.value == "ok") {
            applications.value = (await client.api.v1AppsList()).data
        }
    } catch (error) {
        console.error('Error fetching health status:', error)
    }
})
</script>

<template>
    <div id="footer" class="menu">
        <p>API status: {{ health }}</p>
    </div>
</template>

<style scoped>
#footer {
    padding: 1rem;
    height: auto;
}

#footer>* {
    line-height: 0;
    cursor: pointer;
}
</style>