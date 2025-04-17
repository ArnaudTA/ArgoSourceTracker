<script setup lang="ts">
import { ref, watch } from 'vue';
import { ApplicationGenealogyItem, ApplicationApplicationSummary } from "../api/Api";
import { client } from '../utils/client';
import GenealogyElement from '../components/GenealogyElement.vue';
import ApplicationSummary from '../components/ApplicationSummary.vue';
import Dashboard from '../components/Dashboard.vue';

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
    <div class="page">

        <Dashboard>
            <div class="p-4 space-y-4">
                <div class="bg-white shadow-md rounded-2xl p-6 border border-gray-200">
                    <div class="flex justify-between items-center mb-4">
                        <div>
                            <div>
                                <h2 style="display: inline;">{{ name }}</h2>
                                <UrlIcon style="display: inline; margin-left: 1rem;"
                                    :applicationUrl="application?.applicationUrl"></UrlIcon>
                            </div>
                            <p class="text-sm text-gray-500">Status:
                                <span>
                                    {{ application?.status }}
                                </span>
                            </p>
                        </div>
                    </div>
                </div>
            </div>
        </Dashboard>
        <div id="container">
            <div>
                <ApplicationSummary v-if="application" v-bind:summary="application" :name="name"
                    :namespace="namespace" />
            </div>
            <div>
                <h2>
                    Managed By ↓
                </h2>
                <div id="history">

                    <GenealogyElement v-for="record in history" v-bind="record" />
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
#container {
    padding: 0 3rem;
    display: flex;
    gap: 2rem;
    justify-content: space-between;
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
</style>