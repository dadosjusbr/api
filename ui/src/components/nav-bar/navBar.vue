<template>
  <div>
    <div class="logo">
      <router-link to="/">
        <img src="../../assets/logo.png" class="active" />
      </router-link>
    </div>

    <div class="navMenus">
      <md-toolbar class="md-primary">
        <div class="logoContainer">
          <img src="../../assets/white_logo.png" />
       </div>
       <div class="links">
        <router-link to="/">
          <a class="i"> Início </a>
        </router-link>
        <router-link to="/sobre">
          <a> Sobre </a>
        </router-link>
        <router-link to="/contato">
          <a> Contato </a>
        </router-link>
       </div>
      </md-toolbar>
    </div>
  </div>
</template>

<script>
export default {
  name: "navBar"
};
</script>

<style lang="css" scoped>
.logo img {
  width: 80%;
  height: 80%;
  border: 0;
  margin-top: 8%;
}

.logo {
  background-color: white;
  height: 185px;
  width: 330px;
  margin: 0 auto;
  
}

.logo img:hover {
  background-color: whitesmoke;
}

a {
  font-size: 15px;
  color: white;
  margin: 10px;
  font-family: "Montserrat", sans-serif;
}

.navMenus a:hover {
  font-weight: 1000;
  color: lightblue;
}

.md-toolbar {
    text-align: center;
    width: 100%;
    background-color: #352eba;
    padding: 0 5%;
    z-index: 5;
    height: 3em;
    min-height: 0px;
}

.logoContainer {
  width: 13%;
  margin: 0 auto;
  margin-left: 8%;
}

.logoContainer img {
  width: 80%;
}

.links {
  margin-top: 0%;
  margin-right: 3%;
}

@media only screen and (min-width: 380px) and (max-width: 600px) {
  .logo {
    background-color: white;
    width: 90%;
    margin: 0 auto;
    margin-bottom: 0%;
    height: 0%;
  }

  .logo img {
    width: 60%;
    height: 60%;
    border: 0;
    margin-top: 1%;
    margin-left: 15%;
  }

  .logoContainer {
    display: none;
    width: 0%;
    margin: 0 0;
    margin-left: 0%;
  }

  .links {
    margin-top: 0%;
    margin-right: 0%;
    margin-left: 13%;
  }
}

@media only screen and (min-width: 601px) and  (max-width: 770px) {
  .logoContainer {
    width: 20%;
  }
}

@media only screen and (min-width: 771px) and (max-width: 1025px) {
  .logoContainer {
    width: 15%;
    margin: 0 auto;
    margin-left: 7%;
  }
}


</style>

