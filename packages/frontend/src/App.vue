<script setup lang="ts">
import SideMenu from './components/SideMenu/Main.vue';
</script>

<template>
  <SideMenu id="side-menu" />
  <router-view class="page" />


</template>

<style scoped>
#side-menu {
  position: sticky;
  top: 0;
  height: 100vh;
}

.page {
  flex: 1;
}
</style>