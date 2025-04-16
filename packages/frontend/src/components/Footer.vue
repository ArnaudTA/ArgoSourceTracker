<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { client } from '../utils/client'

const CHECKING_TEXT = 'Checking...'
const health = ref(CHECKING_TEXT)

async function checkHealth() {
    health.value = CHECKING_TEXT
    health.value = await client.api.v1HealthList()
        .then(res => res.data.status)
        .catch(_ => "Can't reach")
        .finally(() => setTimeout(() => {
            checkHealth()
        }, 20000))
}
onMounted(async () => {
    checkHealth()
})
</script>

<template>
    <div id="footer" class="menu">
        <p>API status: {{ health }}</p>
    </div>
</template>

<style scoped>
#footer {
    height: auto;
}

#footer>* {
    line-height: 0;
    cursor: pointer;
}
</style>