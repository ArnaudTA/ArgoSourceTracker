<script setup lang="ts">
import { ref } from 'vue';
import type { ParserSource } from '../api/Api';
import GenericTile from './GenericTile.vue';

defineProps<ParserSource>()

const expanded = ref<boolean>(false)
</script>

<template>
    <GenericTile :title="chart">
        <div>
            <div>
                repo: {{ repoURL || '-' }}</div>
            <div>
                chart: {{ chart || '-' }}</div>
            <div>
                revision: {{ revision || '-' }}</div>

        </div>
        <template v-if="newTags">
            <div>
                <span>Newer Tags</span>
                <div >
                    <div v-for="tag in expanded ? newTags : newTags.slice(0, 3)" :key="tag">{{ tag }}</div>
                </div>
            </div>
        </template>
    </GenericTile>
</template>
