<template>
  <div>
    <div class="navMenus">
      <md-toolbar class="md-primary">
        <div class="logoContainer">
          <img src="../../assets/white_logo.png" />
       </div>
      </md-toolbar>
    </div>
  </div>
</template>

<script>
export default {
  name: "pageFooter"
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
    background-color: #2ab38b;
    padding: 0 5%;
    z-index: 5;
    height: 3em;
    min-height: 0px;
}

.logoContainer {
  flex: right;
  width: 13%;
  margin: 0 auto;
  margin-left: 8%;
}

.logoContainer img {
  width: 90%;
  margin-left: 550%;
}

.links {
  margin-top: 0%;
  margin-right: 3%;
}

</style>