<template>
  <div class="center">
    <div class="circle">{{ info.name + ": " + info.value }}</div>
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
.circle {
  width: 250px;
  height: 250px;
  border-radius: 50%;
  border: purple 3px solid;
  font-size: 25px;
  line-height: 250px;
  text-align: center;
  align-self: center;
  background: #d8bfd8;
}

.center {
  padding: 10px;
}
</style>
