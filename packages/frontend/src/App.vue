<script setup lang="ts">
import Footer from './components/Footer.vue';
import Header from './components/Header.vue';

</script>

<template>
  <div class="app">
    <Header />
    <router-view />
    <div id="separator"></div>
    <Footer />

  </div>
</template>

<style>
.app {
  margin: 0 0;
  padding: 0rem;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}
#separator {
  flex: 1;
}
</style>