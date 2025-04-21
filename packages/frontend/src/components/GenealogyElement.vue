<script setup lang="ts">
import type { TypesParent } from '../api/Api'
import { Button, Card } from 'primevue'
import { goToApp } from '../utils/client'

const props = defineProps<TypesParent>()
</script>

<template>
    <Card class="hidden" :title="`Kind: ${kind}\nName: ${name}`" status="None"
        @click="kind === 'Application' && goToApp(props)">
        <template #title>
            {{ `${namespace}/${name}` }}
        </template>
        <template #subtitle>
            {{ kind }}
        </template>
        <template #footer>
            <Button v-if="applicationUrl" icon="pi pi-external-link" as="a" variant="link" label="ArgoCD"
                :href="applicationUrl" target="_blank" rel="noopener" @click.stop="() => { }" />
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

.hover\:bg-gray-50:hover {
  --un-bg-opacity: 1;
  background-color: var(--p-card-background);
}
</style>
