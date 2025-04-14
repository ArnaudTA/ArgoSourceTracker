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
        <h4>{{ name }}</h4>
    </div>
</template>

<style scoped>
.tile {
    cursor: pointer;
    text-align: left;
    min-width: 8rem;
    min-height: 4rem;
    background-color: lightgray;
    border-radius: .2rem;
    box-shadow: 0 0 .3rem 0rem lightgray;
    border-color: gray;
    border-width: 1px;
    border-style: solid;
    padding: 1rem;
}
</style>