<template>
  <div class="circle">{{ info }}</div>
</template>

<script>
export default {
  name: "infoCard",
  props: {
    info: {
      type: String,
      default: "Não foi possível carregar a informação"
    }
  }
};
</script>

<style scoped>
.circle {
  width: 300px;
  height: 300px;
  border-radius: 50%;
  border: purple 3px solid;
  font-size: 25px;
  line-height: 300px;
  text-align: center;
  background: #d8bfd8;
}
</style>
