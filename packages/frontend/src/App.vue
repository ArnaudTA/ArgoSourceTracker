<script setup lang="ts">
import Footer from './components/Footer.vue';
import Header from './components/Header.vue';

</script>

<template>
  <div class="app">
    <Header />
    <div id="content">
      <router-view />

    </div>
    <div id="separator"></div>
    <Footer />

  </div>
</template>

<style scoped>
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
#content{
  padding: 2rem 2rem;
}
</style>