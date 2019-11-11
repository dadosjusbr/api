<template>
  <div>
    <nav-bar />
    <!-- <div id="nav">
      <router-link to="/">Home</router-link>|
      <router-link to="/about">About</router-link>
    </div>-->
    <orgao-summary :orgaoName="orgaoName" :infos="infos" />
  </div>
</template>

<script>
import orgaoSummary from "@/components/orgao/orgaoSummary.vue";
import navBar from "@/components/nav-bar/navBar.vue";

export default {
  components: {
    orgaoSummary,
    navBar
  },
  data() {
    return {
      orgaoName: "TJPB",
      infos: ["info 1", "info2", "info 3", "info4"]
    };
  }
};
</script>

<style>
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
