<template>
  <div>
    <div>
      <button v-on:click="previousMonth()">Anterior</button>
      <button v-on:click="nextMonth()">Proximo</button>
    </div>
    <graph-point
      width="500"
      type="line"
      :options="chartOptions"
      :series="series"
    ></graph-point>
  </div>
</template>

<script>
import graphPoint from "@/components/agency/graphPoint.vue";

export default {
  name: "graphContainer",
  components: {
    graphPoint
  },
  data: function() {
    return {
      chartOptions: {
        xaxis: {
          categories: [1991, 1992, 1993, 1994, 1995, 1996, 1997, 1998]
        }
      },
      series: [
        {
          data: [30, 40, 45, 50, 49, 60, 70, 81]
        }
      ]
    };
  },
  methods: {
    nextMonth() {
      this.chartOptions = {
        chart: {
          id: "vuechart-example"
        },
        xaxis: {
          categories: [1991]
        }
      };
      this.series = [
        {
          name: "Vue Chart",
          data: [30]
        }
      ];
    },
    previousMonth() {
      this.chartOptions = {
        chart: {
          id: "vuechart-example"
        },
        xaxis: {
          categories: [1991, 1995]
        }
      };
      this.series = [
        {
          name: "Vue Chart",
          data: [55, 44]
        }
      ];
    }
  }
};
</script>

<style scoped>
.container {
  border: coral 2px solid;
  height: 500px;
  padding-top: 50px;
}
</style>
