<script setup lang="ts">
defineProps<{ applicationUrl: string }>()
</script>

<template>
    <div v-if="applicationUrl"><a :href="applicationUrl" target="_blank" title="Open application">&#128279;</a></div>
</template>