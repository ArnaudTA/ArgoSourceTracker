<script setup lang="ts">
import type { TypesSummary } from '../api/Api';
import router from '../router';
defineProps<{
    application: TypesSummary
}>()
</script>

<template>
    <div class="tile" @click="router.push({
        name: 'Application',
        params: {
            name: application.name,
            namespace: application.namespace
        }
    })">
        <h4>{{ application.namespace +"/"+ application.name }}</h4>
    </div>
</template>

<style scoped>
.tile {
    cursor: pointer;
    text-align: left;
    min-width: 8rem;
    min-height: 4rem;
    border-radius: .2rem;
    box-shadow: 0 0 .3rem 0rem lightgray;
    border-width: 1px;
    border-style: solid;
    padding: 1rem;
}
</style>