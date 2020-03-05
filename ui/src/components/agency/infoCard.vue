<template>
  <div class="center">
    <div class="circle">
      {{ info.name.replace("_", " ") + ":" + "\n" + info.value }}
    </div>
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
  width: 275px;
  height: 275px;
  border-radius: 50%;
  border: #182825 3px solid;
  font-size: 17px;
  line-height: 250px;
  text-align: center;
  font-family: "Montserrat", sans-serif;
  background: #e1e8ed;
}

.center {
  padding: 10px;
  align-items: flex-start;
}
</style>
