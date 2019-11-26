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
    <!-- <agency-page-container /> -->
    <state-page-container />
  </div>
</template>

<script>
import navBar from "@/components/nav-bar/navBar.vue";
// import agencyPageContainer from "@/components/agency/agencyPageContainer.vue";
import statePageContainer from "@/components/state/statePageContainer.vue";

export default {
  components: {
    // agencyPageContainer,
    navBar,
    statePageContainer
  }
};
</script>

<style>
@import url("https://fonts.googleapis.com/css?family=Montserrat&display=swap");
.divisoria-colorida {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr 1fr 1fr;
  height: 14px;
  position: relative;
  width: 100%;
}

.cor-1 {
  background-color: #d00000;
}
.cor-2 {
  background-color: #d00000;
}
.cor-3 {
  background-color: #d00000;
}
.cor-4 {
  background-color: #d00000;
}
.cor-5 {
  background-color: #d00000;
}

#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

#nav {
  padding: 30px;
}

#nav a {
  font-weight: bold;
  color: #2c3e50;
}

#nav a.router-link-exact-active {
  color: #42b983;
}
</style>
