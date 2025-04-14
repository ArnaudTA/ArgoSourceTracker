<script setup lang="ts">
import type { ParserApplicationSummary } from '../api/Api'
import SourceElement from './SourceElement.vue';

// Props
defineProps<{
    summary: ParserApplicationSummary
    name: string
}>()

// Badge class utilitaire
const statusBadge = (status?: string) => {
    switch (status?.toLowerCase()) {
        case 'Up-to-date':
            return 'text-green-600 bg-green-100'
        case 'outdated':
            return 'text-yellow-600 bg-yellow-100'
        default:
            return 'text-gray-600 bg-gray-100'
    }
}
</script>

<template>
    <div class="p-4 space-y-4">
        <div class="bg-white shadow-md rounded-2xl p-6 border border-gray-200">
            <div class="flex justify-between items-center mb-4">
                <div>
                    <h2 class="text-xl font-semibold text-gray-800">{{ summary.instance }}</h2>
                    <p class="text-sm text-gray-500">Status:
                        <span :class="statusBadge(summary.status)">
                            {{ summary.status }}
                        </span>
                    </p>
                </div>
            </div>
        </div>
    </div>
    <h3>Charts</h3>
    <SourceElement v-for="(chart, index) in summary.charts" :key="index" class="border-b hover:bg-gray-50"
        v-bind="chart" />
</template>
