<script setup lang="ts">
import type { ApplicationTrackRecord } from '../api/Api';
import GenericTile from './GenericTile.vue';

defineProps<ApplicationTrackRecord>()
</script>

<template>
    <GenericTile :title="`Kind: ${kind}\nName: ${name}`">
        <div>
            <div class="node-title direction-left">
                {{ name }}</div>
            <div class="node-title direction-left">
                {{ kind }}</div>
        </div>
        <div v-if="applicationUrl" class="node-status-icon"><a :href="applicationUrl" target="_blank"
                title="Open application">&#128279;</a></div>
    </GenericTile>
</template>
