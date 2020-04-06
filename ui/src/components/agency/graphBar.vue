<template>
  <div class="graph">
    <apexcharts
      width="100%"
      height="500"
      type="bar"
      :options="options"
      :series="series"
    />
  </div>
</template>

<script>
import VueApexCharts from "vue-apexcharts";

export default {
  name: "graphBar",
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
  align-self: center;
}
</style>