<script setup lang="ts">
import { useSideMenuStore } from '../../stores/sideMenu';

const sideMenuStore = useSideMenuStore()
const slots = defineSlots<{
  condensed?(): any
  default(): any
}>()
</script>
<template>
  <div class="drawer">
    <template v-if=slots.condensed>
      <slot name="condensed" class="condensed"></slot>
      <slot class="content"></slot>
    </template>
    <template v-else>
      <slot v-if="!sideMenuStore.condensed"></slot>
    </template>
  </div>
</template>

<style scoped>
.drawer {
  padding: 0.5rem;
  min-height: 2rem;
}
.condensed {
  width: 4rem;
}
</style>