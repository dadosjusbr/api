<template>
  <div>
    <state-page-container />
  </div>
</template>

<script>
import statePageContainer from "@/components/state/statePageContainer.vue";

export default {
  components: {
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
  background-color: #6a757a;
}
.cor-2 {
  background-color: #6a757a;
}
.cor-3 {
  background-color: #6a757a;
}
.cor-4 {
  background-color: #6a757a;
}
.cor-5 {
  background-color: #6a757a;
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
