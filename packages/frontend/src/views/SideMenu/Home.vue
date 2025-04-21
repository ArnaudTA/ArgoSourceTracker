<script setup lang="ts">
import { useHomeStore } from "../../stores/home.ts";
import GenericDrawer from "../../components/SideMenu/GenericDrawer.vue";
import { MultiSelect } from "primevue";
import { onMounted } from "vue";

const homeStore = useHomeStore()
onMounted(homeStore.refreshApps)
</script>

<template>
    <GenericDrawer>
        <MultiSelect v-model="homeStore.selectedStatus" :options="homeStore.StatusChoice" showClear optionLabel="name"
            filter placeholder="Select Status" size="small" class="w-full" />
    </GenericDrawer>
</template>
