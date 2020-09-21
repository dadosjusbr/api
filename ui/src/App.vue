<template>
  <div class="app">
    <link
      rel="stylesheet"
      href="https://fonts.googleapis.com/css?family=Roboto:300,400,500,700,400italic|Material+Icons"
    />
    <b-container fluid>
      <nav-bar />
      <router-view />
      <page-footer />
    </b-container>
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
  font-family: "Roboto Condensed";
  background-color: #ffffff;
}

@media only screen and (max-width: 600px) {
  .contentContainer {
    padding-bottom: 368px;
  }
}
</style>
