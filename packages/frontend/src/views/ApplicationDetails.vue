<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { ApplicationTrackRecord } from "../api/Api";
import { client } from '../utils/client';
import TrackRecordTile from '../components/TrackRecordTile.vue';
const props = defineProps<{
    name: string
}>()

const history = ref<ApplicationTrackRecord[]>([])
onMounted(async () => {
    history.value = (await (client.api.v1AppsOriginList(props.name))).data
})
</script>
<template>
    <TrackRecordTile v-for="record in history" v-bind="record" />
</template>