<script setup lang="ts">
import { TypesApplicationStatus } from '../api/Api';
import { statusClass } from '../utils/client';


withDefaults(defineProps<{
    title?: string
    status: TypesApplicationStatus | 'None'
}>(), {
    status: 'None'
})

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
</style>