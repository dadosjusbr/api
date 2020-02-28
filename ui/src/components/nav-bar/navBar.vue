<template>
  <div>
    <div class="logo">
      <router-link to="/">
        <img src="../../assets/logo.png" class="active" />
      </router-link>
    </div>

    <div class="navMenus">
      <router-link to="/">
        <a> Início </a>
      </router-link>
      <router-link to="/sobre">
        <a> Sobre </a>
      </router-link>
      <router-link to="/contato">
        <a> Contato </a>
      </router-link>
    </div>
  </div>
</template>

<script>
export default {
  name: "navBar"
};
</script>

<style>
img {
  width: 100%;
  height: 100%;
  border: 0;
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

.navMenus {
  background-color: black;
  text-align: center;
}

a {
  color: white;
  display: inline-block;
  margin: 10px;
  font-family: "Montserrat", sans-serif;
}

.navMenus a:hover {
  font-weight: 1000;
  color: lightblue;
}
</style>
