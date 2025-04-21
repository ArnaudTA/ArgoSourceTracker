<script setup lang="ts">
import { Button, Card } from 'primevue';
import type { TypesParent } from '../api/Api';
import { goToApp } from '../utils/client';

const props = defineProps<TypesParent>()
</script>

<template>
    <Card :title="`Kind: ${kind}\nName: ${name}`" @click="kind === 'Application' && goToApp(props)" :status="'None'">
        <template #title>{{ namespace+"/"+name }}</template>
        <template #subtitle>{{ kind }}</template>
        <template #footer>
            <Button v-if="applicationUrl" icon="pi pi-external-link" as="a" variant="link" label="ArgoCD"></Button>
        </template>
    </Card>
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

.names {}
</style>