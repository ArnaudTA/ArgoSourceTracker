<script setup lang="ts">
import type { ParserApplicationSummary } from '../api/Api';
import router from '../router';
defineProps<{
    name: string,
    application: ParserApplicationSummary
}>()
</script>

<template>
    <div class="tile" @click="router.push({
        name: 'Application',
        params: {
            name: name
        }
    })">
        {{ name }}
    </div>
</template>

<style scoped>
.tile {
    min-width: 10rem;
    background-color: lightgray;
    border-radius: .2rem;
    box-shadow: 0 0 .3rem 0rem lightgray;
    border-color: gray;
    border-width: 1px;
    border-style: solid;
}
</style>