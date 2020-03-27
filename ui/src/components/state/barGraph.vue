<template>
  <div class="graphContainer">
    <apexcharts
      width="100%"
      height="400"
      type="bar"
      :options="options"
      :series="series"
    />
  </div>
</template>

<script>
import VueApexCharts from "vue-apexcharts";

export default {
  name: "barGraph",
  components: {
    apexcharts: VueApexCharts
  },
  props: {
    options: {
      type: Object,
      default: null
    },
    series: {
      type: Array,
      default: null
    }
  }
};
</script>

<style scoped>
.graphContainer {
  width: 69%;
  height: auto;
  padding: 10px;
  position: relative;
}

@media only screen and  (max-width: 379px) {
  .graphContainer {
    width: 100%;
    height: auto;
    padding: 0%;
    position: relative;
    margin-top: -7%;
  }
}

@media only screen and (min-width: 380px) and  (max-width: 600px) {
  .graphContainer {
    width: 100%;
    height: auto;
    padding: 8px;
    position: relative;
    margin-top: -7%;
  }
}

@media only screen and (min-width: 601px) and  (max-width: 770px) {
  .graphContainer {
    width: 102%;
    height: auto;
    padding: 10px;
    position: relative;
  }
}

@media only screen and  (min-width: 771px) and (max-width: 1025px)  {
  .graphContainer {
    width: 86%;
    height: auto;
    padding: 10px;
    position: relative;
  }
}



</style>
