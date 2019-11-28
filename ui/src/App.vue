<template>
  <div>
    <nav-bar />
    <div class="divisoria-colorida">
      <div class="cor-1"></div>
      <div class="cor-2"></div>
      <div class="cor-3"></div>
      <div class="cor-4"></div>
      <div class="cor-5"></div>
    </div>
    <router-view />
  </div>
</template>

<script>
import navBar from "@/components/nav-bar/navBar.vue";

export default {
  components: {
    navBar
  }
};
</script>

<style></style>
