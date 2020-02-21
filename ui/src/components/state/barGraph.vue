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
  width: 1000px;
  height: auto;
  margin: 0 auto;
  padding: 10px;
  position: relative;
}
</style>
