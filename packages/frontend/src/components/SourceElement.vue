<script setup lang="ts">
import type { TypesChartSummary } from '../api/Api'
import { Card } from 'primevue'
import { ref } from 'vue'
import { statusClass } from '../utils/client';

defineProps<TypesChartSummary>()

const expanded = ref<boolean>(false)
</script>

<template>
    <Card :title="chart" :class="statusClass[status ?? 'None']" class="card">
        <template #title>
            {{ chart }}
        </template>
        <template #subtitle>
            revision: {{ revision }}
        </template>
        <template #content>
            <div>
                repo: {{ repoURL }}
            </div>
            <div v-if="error">
                error: {{ error }}
            </div>
            <div v-if="newTags" class="mt-2">
                <span>Newer Tags</span>
                <div v-for="tag in expanded ? newTags : newTags.slice(0)" :key="tag">
                    {{ tag }}
                </div>
            </div>
        </template>
    </Card>
</template>

<style>
.card {
    position: relative;
}
</style>