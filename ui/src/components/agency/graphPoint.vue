<template>
  <div>
    <div>{{ this.options }}</div>
    <div>{{ this.series }}</div>
    <apexcharts
      width="50%"
      height="350"
      type="scatter"
      :options="options"
      :series="series"
    />
  </div>
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

<style></style>
