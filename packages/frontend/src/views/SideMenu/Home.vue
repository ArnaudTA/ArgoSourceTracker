<script setup lang="ts">
import { ref } from "vue";
import { Listbox } from 'primevue'
import { useHomeStore } from '../../stores/home.ts'
import { ApplicationApplicationStatus } from "../../api/Api.ts";
import GenericDrawer from "../../components/SideMenu/GenericDrawer.vue";

const homeStore = useHomeStore()

const cities = ref<Array<{ code: ApplicationApplicationStatus }>>([
    { code: ApplicationApplicationStatus.Outdated },
]);
</script>

<template>
    <GenericDrawer>
        <Listbox v-model="homeStore.filters" :options="cities" multiple optionLabel="name" class="w-full md:w-56" />
    </GenericDrawer>
</template>
