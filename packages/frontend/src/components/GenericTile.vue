<script setup lang="ts">

export type TileStatus = 'Up-to-date' | 'Ignored' | 'Outdated' | 'Checking' | 'Unknown' | 'Error' | 'None'
withDefaults(defineProps<{
    title?: string
    status: TileStatus
}>(), {
    status: 'None'
})

const statusClass: Record<TileStatus, string> = {
    "Up-to-date": "uptodate",
    Checking: "checking",
    Error: "error",
    Ignored: "ignored",
    Unknown: "unknown",
    Outdated: "outdated",
    None: "none"
}
</script>

<template>

    <div :class="`node ${statusClass[status]}`" :title="title">
        <div class="node-content">
            <slot></slot>
        </div>
    </div>
</template>

<style scoped>
.node {
    background-color: #fff;
    color: #363c4a;
    padding: 0.5rem 1rem 0.5rem 1rem;
    box-shadow: 1px 1px 1px #ccd6dd;
    border-radius: 4px;
    border: 1px solid transparent;
    min-width: 282px;
    min-height: 52px;
    position: relative;
    padding-bottom: 0.5rem;
}

.node-content {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    gap: 2rem;
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

.node::after {
    height: 0.2rem;
    content: "";
    position: absolute;
    bottom: 0;
    left: 0;
    width: 100%;
}
.node.uptodate::after {
    background-color: rgb(2, 231, 2);
}
.node.outdated::after {
    background-color: rgb(255, 178, 35);
}
.node.error::after {
    background-color: red;
}
.node.ignored::after {
    background-color: rgba(211, 211, 211, 0.562);
}
</style>