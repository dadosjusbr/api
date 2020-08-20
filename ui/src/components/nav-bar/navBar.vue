<template>
  <div class="navBar">
    <div class="auxBar">
      <img src="../../assets/icon_dadosjusbr_colored.svg" />
      <router-link to="/">
        <b style="color: white; font-size: 18px;"> Ver Gráficos </b>
      </router-link>
    </div>
  </div>
</template>

<script>
export default {
  name: "navBar",
};
</script>

<style scoped>
.navBar {
  background-color: #7f3d8b;
  height: 185px;
  padding: 50px 100px 0px 100px;
}
.auxBar {
  display: flex;
  justify-content: space-between;
  border-bottom: 1px solid white;
  padding-right: 100px;
  padding-bottom: 35px;
  align-items: center;
}
</style>
