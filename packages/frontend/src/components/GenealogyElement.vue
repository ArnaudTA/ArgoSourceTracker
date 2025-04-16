<script setup lang="ts">
import type { ApplicationGenealogyItem } from '../api/Api';
import { goToApp } from '../utils/client';
import GenericTile from './GenericTile.vue';
import UrlIcon from './UrlIcon.vue';

const props = defineProps<ApplicationGenealogyItem>()
</script>

<template>
    <GenericTile :title="`Kind: ${kind}\nName: ${name}`" @click="kind === 'Application' && goToApp(props)" :status="'None'">
        <div class="kind">
            {{ kind }}
        </div>
        <div class="names">
            <div>{{ name }}</div>
            <div>{{ namespace }}</div>

        </div>
        <UrlIcon v-if="applicationUrl" :applicationUrl="applicationUrl" class="url"></UrlIcon>
    </GenericTile>
</template>

<style>
.kind {
    align-content: center;
    width: 6rem;
    font-size: small;
    text-align: center;
}
.url {
    align-content: center;
    font-size: large;
    text-align: center;
}

.names {
    flex: 1;
}
</style>