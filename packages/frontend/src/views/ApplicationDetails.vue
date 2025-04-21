<script setup lang="ts">
import { ref, watch } from 'vue';
import { TypesSummary, TypesParent } from "../api/Api";
import { client } from '../utils/client';
import GenealogyElement from '../components/GenealogyElement.vue';
import ApplicationSummary from '../components/ApplicationSummary.vue';
import { Button, Timeline, Toolbar } from 'primevue';

const props = defineProps<{
    name: string
    namespace: string
}>()

const history = ref<TypesParent[]>([])
const application = ref<TypesSummary>()

watch(props, getInfos, { immediate: true })

async function getInfos() {
    history.value = (await (client.api.v1AppsOriginList(props.name, props.namespace))).data
    application.value = (await (client.api.v1AppsDetail(props.name, props.namespace))).data
}

</script>
<template>
    <div class="page">
        <Toolbar>
            <template #start>
                <div class="p-4 space-y-4">
                    <div>
                        <h2 style="display: inline;">{{ name }}</h2>
                        <UrlIcon style="display: inline; margin-left: 1rem;"
                            :applicationUrl="application?.applicationUrl">
                        </UrlIcon>
                    </div>
                    <p class="text-sm text-gray-500">Status:
                        <span>
                            {{ application?.status }}
                        </span>
                    </p>
                </div>
            </template>
            <template #end>
                <Button v-if="application?.applicationUrl" icon="pi pi-external-link" as="a" variant="link"
                    label="ArgoCD" :href="application.applicationUrl" target="_blank"></Button>
            </template>
        </Toolbar>
        <div class="flex flex-row p-8 gap-16">
            <div class="grow">
                <h2 class="text-3xl mb-4">Charts</h2>
                <ApplicationSummary v-if="application" v-bind:summary="application" :name="name"
                    :namespace="namespace" />
            </div>
            <div class="grow">
                <h2 class="text-3xl mb-4">Genealogy</h2>
                <Timeline :value="history" align="left">
                    <template #content="slotProps">
                        <GenealogyElement v-bind="slotProps.item" class="my-2" />
                    </template>
                </Timeline>
            </div>
        </div>
    </div>
</template>

<style>
.p-timeline-event-opposite {
    flex: 0 !important;
}
</style>