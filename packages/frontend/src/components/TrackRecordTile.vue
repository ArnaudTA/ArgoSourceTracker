<script setup lang="ts">
defineProps<{
    name: string
    kind: string
}>()
</script>

<template>
    <div class="tile">
        {{ name }}
        {{ kind }}
    </div>
</template>

<style scoped>
.tile {
    min-width: 10rem;
    background-color: lightgray;
    border-radius: .2rem;
    box-shadow: 0 0 .3rem 0rem lightgray;
    border-color: gray;
    border-width: 1px;
    border-style: solid;
}
</style>