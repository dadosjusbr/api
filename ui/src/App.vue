<template>
  <div class="app">
    <link
      rel="stylesheet"
      href="https://fonts.googleapis.com/css?family=Roboto:300,400,500,700,400italic|Material+Icons"
    />
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
    pageFooter,
  },
};
</script>

<style>
@import url("https://fonts.googleapis.com/css?family=Roboto+Condensed");
.app {
  font-family: "Roboto Condensed", sans-serif;
  background-color: #F4F4F4;
}
</style>
