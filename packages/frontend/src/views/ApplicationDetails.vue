<script setup lang="ts">
import { ref, watch } from 'vue';
import { ApplicationGenealogyItem, ApplicationApplicationSummary } from "../api/Api";
import { client } from '../utils/client';
import GenealogyElement from '../components/GenealogyElement.vue';
import ApplicationSummary from '../components/ApplicationSummary.vue';

const props = defineProps<{
    name: string
    namespace: string
}>()

const history = ref<ApplicationGenealogyItem[]>([])
const application = ref<ApplicationApplicationSummary>()
watch(props, getInfos, { immediate: true })
async function getInfos() {
    history.value = (await (client.api.v1AppsOriginList(props.name, props.namespace))).data
    application.value = (await (client.api.v1AppsDetail(props.name, props.namespace))).data
}
</script>
<template>
    <div id="container">
        <div>
            <ApplicationSummary v-if="application" v-bind:summary="application" :name="name" :namespace="namespace" />
        </div>
        <div id="history">
            <h2>
                Managed By ↓
            </h2>
            <GenealogyElement v-for="record in history" v-bind="record" />
        </div>
    </div>
</template>

<style scoped>
#container {
    padding: 1rem;
    display: flex;
    gap: 2rem;
}

#container>* {
    flex: 1;
}

#history {
    display: flex;
    flex-direction: column;
    width: max-content;
    gap: 2rem;
}

#history>h2 {
    display: inline;
    margin: 0;
}
</style>