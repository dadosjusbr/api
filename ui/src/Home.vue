<template>
  <div>
    <state-page-container />
  </div>
</template>

<script>
import statePageContainer from "@/components/state/statePageContainer.vue";

export default {
  components: {
    statePageContainer,
  },
};
</script>

<style>
@import url("https://fonts.googleapis.com/css?family=Montserrat&display=swap");

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

body {
  margin: 0 !important;
}

#nav a {
  font-weight: bold;
  color: #2c3e50;
}

#nav a.router-link-exact-active {
  color: #42b983;
}
</style>
