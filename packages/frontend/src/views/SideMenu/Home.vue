<script setup lang="ts">
import { useHomeStore } from "../../stores/home.ts";
import GenericDrawer from "../../components/SideMenu/GenericDrawer.vue";
import { Checkbox } from "primevue";
import { onMounted } from "vue";

const homeStore = useHomeStore()
onMounted(homeStore.refreshApps)
</script>

<template>
    <GenericDrawer>
        <div v-for="status of homeStore.StatusChoice" :key="status.key" class="flex items-center gap-3 mb-3">
            <Checkbox v-model="homeStore.selectedStatus" :inputId="status.key" name="status" :value="status.name" />
            <label :for="status.key">{{ status.name }} ({{ homeStore.stats[status.key] ?? 0 }})</label>
        </div>
    </GenericDrawer>
</template>
