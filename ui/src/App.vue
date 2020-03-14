<template>
  <div>
    <nav-bar />
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

<style>
@import url("https://fonts.googleapis.com/css?family=Montserrat&display=swap");
</style>


