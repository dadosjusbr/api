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
    <page-footer />
  </div>
</template>

<script>
import navBar from "@/components/nav-bar/navBar.vue";
import pageFooter from "@/components/footer/pageFooter.vue";

export default {
  components: {
    navBar,
    pageFooter
  }
};
</script>

<style></style>
