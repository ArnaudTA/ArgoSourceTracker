<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { ApplicationTrackRecord, ParserApplicationSummary } from "../api/Api";
import { client } from '../utils/client';
import GenealogyElement from '../components/GenealogyElement.vue';
import ApplicationSummary from '../components/ApplicationSummary.vue';

const props = defineProps<{
    name: string
}>()

const history = ref<ApplicationTrackRecord[]>([])
const application = ref<ParserApplicationSummary>()
onMounted(async () => {
    history.value = (await (client.api.v1AppsOriginList(props.name))).data
    application.value = (await (client.api.v1AppsDetail(props.name))).data
})
</script>
<template>
    <div id="container">
        <div id="charts">
            <ApplicationSummary v-if="application" v-bind:summary="application" :name="name" />

        </div>
        <div id="history">
            <h2>
                Managed By
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
}

</style>