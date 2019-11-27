<template>
  <apexcharts
    width="100%"
    height="350"
    type="scatter"
    :options="options"
    :series="series"
  />
</template>

<script>
import VueApexCharts from "vue-apexcharts";

export default {
  name: "graphPoint",
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

<style>
.graph {
  width: 95%;
  height: auto;
  margin: 0 auto;
  padding: 10px;
  position: relative;
}
</style>
