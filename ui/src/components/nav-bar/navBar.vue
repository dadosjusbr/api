<template>
  <div>
    <div class="navMenus">
      <md-toolbar class="md-primary">
        <div class="logoContainer">
          <router-link to="/">
            <img src="../../assets/white_logo.png" />
          </router-link>
        </div>
        <div class="links">
          <router-link to="/sobre">
            <a> Sobre </a>
          </router-link>
        </div>
      </md-toolbar>
    </div>
  </div>
</template>

<script>
export default {
  name: "navBar",
};
</script>

<style scoped>
.links {
  margin: 10px 10px 10px 70%;
}
a {
  font-size: 15px;
  color: white;
}
a:hover {
  font-weight: 1000;
  color: lightblue;
}
.md-toolbar {
  background-color: #2ab38b;
  height: 3em;
  min-height: 0px;
  padding: 0 5%;
}
.logoContainer {
  width: 13%;
  margin-left: 8%;
}
@media only screen and (max-width: 800px) {
  .logoContainer {
    width: 20%;
    margin-left: 8%;
  }
  .links {
    margin: 10px 10px 10px 60%;
  }
}
@media only screen and (max-width: 700px) {
  .logoContainer {
    width: 25%;
    margin-left: 8%;
  }
  .links {
    margin: 10px 10px 10px 50%;
  }
}
@media only screen and (max-width: 450px) {
  .logoContainer {
    width: 30%;
    margin-left: 8%;
  }
  .links {
    margin: 10px 10px 10px 40%;
  }
}
@media only screen and (max-width: 320px) {
  .logoContainer {
    width: 30%;
    margin-left: 8%;
  }
  .links {
    margin: 10px 10px 10px 40%;
  }
}
</style>
