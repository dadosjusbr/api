<template>
  <div class="container">
    <h2 class="score">
      {{ info.value }}
      <p class="text">{{ info.name.replace("_", " ") }}</p>
    </h2>
  </div>
</template>

<script>
export default {
  name: "infoCard",
  props: {
    info: {
      type: Object,
      default: null
    }
  }
};
</script>

<style scoped>
.container {
  height: 200px;
  width: 200px;
}

.score {
  padding-top: 60px;
  text-align: center;
  background: white;
  border: 3px solid rgb(4, 4, 173);
  margin-top: 10px;
  font-size: 20px;
  height: 150px;
  width: 180px;
  color:rgb(4, 4, 173);
  font-weight: 500;
}

.text {
  padding-top: 60px;
  padding-left: 5px;
  text-align: left;
  position: relative;
  font-weight: 200;
  margin-top: -20px !important;
  font-size: 15px;
  color: rgb(4, 4, 173);
}
</style>
