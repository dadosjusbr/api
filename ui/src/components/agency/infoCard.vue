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
  height: 120px;
  width: 180px;
}

.score {
  padding-top: 29%;
  text-align: center;
  background: #87bba2;
  border: 1px solid black;
  margin-top: 10px;
  font-size: 20px;
  height: 120px;
  width: 180px;
  color: black;
  font-weight: 500;
}

.text {
  padding-top: 50px;
  padding-left: 5px;
  text-align: left;
  position: relative;
  font-weight: 410;
  margin-top: -20px !important;
  font-size: 16px;
  color: black;
}
</style>
