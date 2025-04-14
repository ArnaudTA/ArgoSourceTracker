<script setup lang="ts">
import router from '../router';

</script>

<template>
    <div id="header" class="menu">
        <div id="HeaderHome" @click="router.push({ name: 'Home' })">
            <h2>
                Argo Source Tracker
            </h2>
        </div>
    </div>
</template>

<style>
#header {
    padding: 1rem;
    height: auto;
}

#HeaderHome>* {
    line-height: 0;
    cursor: pointer;
}
</style>