<template>
  <div class="container">
    <h2 class="score">
      {{ info.name.replace("_", " ") }}
      <p class="text">{{ info.value }}</p>
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
  text-align: center;
  background: #e1e8ed;
  border: 3px solid #182825;
  border-radius: 30px;
  padding-top: 30px;
  font-size: 33px;
  height: 180px;
  width: 200px;
  color: black;
  font-weight: 500;
  margin-bottom: 0;
}

.text {
  padding-top: 40px;
  position: relative;
  font-weight: 200;
  margin-top: -20px !important;
  font-size: 20px;
}
</style>
