<script setup lang="ts">
import { ref } from 'vue';
import type { ApplicationApplicationStatus, ApplicationChartSummary } from '../api/Api';
import GenericTile from './GenericTile.vue';

defineProps<ApplicationChartSummary>()

const expanded = ref<boolean>(false)
</script>

<template>
    <GenericTile :title="chart" :status="(status as ApplicationApplicationStatus)">
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
                <div>
                    <div v-for="tag in expanded ? newTags : newTags.slice(0, 3)" :key="tag">{{ tag }}</div>
                </div>
            </div>
        </template>
    </GenericTile>
</template>
