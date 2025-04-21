<script setup lang="ts">
import type { TypesSummary } from '../api/Api'
import SourceElement from './SourceElement.vue';

// Props
defineProps<{
    summary: TypesSummary
    name: string
    namespace: string
}>()
</script>

<template>
    <h2>Charts</h2>
    <p v-if="summary.charts.length === 0">No Charts found</p>
    <div class="chart-list">
        <SourceElement v-for="(chart, index) in summary.charts" :key="index" class="border-b hover:bg-gray-50"
            v-bind="chart" />
    </div>
</template>

<style scoped>
.chart-list {
    display: flex;
    min-width: 25rem;
    justify-items: stretch;
    align-items: stretch;
    justify-content: stretch;
    justify-self: stretch;
}
</style>
