<template>
  <div>
    <div class="topnav">
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
.topnav {
  background-color: white;
  height: 170px;
}

.topnav a {
  display: block;
  text-align: center;
}

.topnav img:hover {
  background-color: whitesmoke;
}

.topnav img.active {
  align-items: center;
  width: 282px;
  height: 158px;
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
