<template>
  <div>
    <button v-on:click="previousYear()">Anterior</button>
    <button v-on:click="nextYear()">Proximo</button>
    <div>{{ this.data.MonthTotals }}</div>
    <div>{{ this.data }}</div>
    <bar-graph :options="chartOptions" :series="series" />
  </div>
</template>

<script>
import barGraph from "@/components/state/barGraph.vue";

export default {
  name: "agency",
  currentYear: 2019,
  components: {
    barGraph
  },
  data() {
    return {
      data: {},
      chartOptions: {
        chart: {
          stacked: true,
          toolbar: {
            show: true
          },
          zoom: {
            enabled: true
          }
        },
        responsive: [
          {
            breakpoint: 480,
            options: {
              legend: {
                position: "bottom",
                offsetX: -10,
                offsetY: 0
              }
            }
          }
        ],
        plotOptions: {
          bar: {
            horizontal: false
          }
        },
        xaxis: {
          categories: ["JAN", "FEV", "MAR"]
        },
        legend: {
          position: "right",
          offsetY: 40
        },
        fill: {
          opacity: 1
        }
      }
    };
  },
  methods: {
    nextYear() {
      this.$http
        .get("/orgao/totais/TJPB/" + this.currentYear + 1)
        .then(response => (this.data = response.data));
    },
    previousYear() {
      this.$http
        .get("/orgao/totais/TJPB/" + this.currentYear - 1)
        .then(response => (this.data = response.data));
    }
  },
  computed: {
    series() {
      let others = this.data.MonthTotals.map(month => month["Others"]);
      let wages = this.data.MonthTotals.map(month => month["Wage"]);
      let perks = this.data.MonthTotals.map(month => month["Perks"]);
      return [
        {
          name: "others",
          data: others
        },
        {
          name: "perks",
          data: perks
        },
        {
          name: "wages",
          data: wages
        }
      ];
    }
  },
  mounted() {
    this.$http
      .get("/orgao/totais/TJPB/4")
      .then(response => (this.data = response.data));
  }
};
</script>

<style></style>
