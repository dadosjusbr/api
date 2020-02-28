<template>
  <div class="footer">
    <div class="logoContainer">
      <img src="../../assets/white_logo.png" />
    </div>
  </div>
</template>

<script>
export default {
  name: "pageFooter"
};
</script>

<style scoped>
.footer {
  padding-top: 8px;
  padding-bottom: 5px;
  background-color: black;
  height: 70px;
  margin-top: 5px;
}

.logoContainer {
  height: 60px;
  width: 250px;
  margin: 0 auto;
}

img {
  width: 100%;
  height: 100%;
  border: 0;
}
</style>
