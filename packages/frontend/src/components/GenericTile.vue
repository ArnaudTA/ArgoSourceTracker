<script setup lang="ts">

defineProps<{
    title?: string
}>()
</script>

<!-- `Kind: ${kind}\nName: ${name}` -->
<template>

    <div class="node" :title="title">

        <div class="node-content">
            <slot></slot>
        </div>
    </div>
</template>

<style scoped>
.node {
    background-color: #fff;
    color: #363c4a;
    padding-left: 1em;
    padding-right: 1em;
    margin: 10px;
    box-shadow: 1px 1px 1px #ccd6dd;
    border-radius: 4px;
    border: 1px solid transparent;
    min-width: 282px;
    min-height: 52px;
}

.node-content {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
}

.node-kind-icon {
    cursor: pointer;
    text-align: center;
    left: 0;
    top: 8px;
    width: 60px;
    line-height: 2;
    color: #495763;
    text-decoration: none;
    justify-self: center;
}
</style>