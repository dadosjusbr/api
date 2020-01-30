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
      <a> Sobre </a>
      <a> Contato </a>
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
  background-color: #182825;
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
  width: 15%;
}

.navMenus {
  background-color: #182825;
  text-align: center;
}

a {
  color: white;
  display: inline-block;
  margin: 10px;
  font-family: "Montserrat", sans-serif;
}

.navMenus a:hover {
  background-color: black;
}
</style>
