<template>
  <div class="agencyContainer">
    <h2 class="agencyName">{{ this.agencyName }}</h2>
    <div class="buttonContainer">
      <button class="button" v-on:click="previousYear()">&#8249;</button>
      <button class="button" v-on:click="nextYear()">&#8250;</button>
    </div>
    <bar-graph :options="chartOptions" :series="series" />
  </div>
</template>

<script>
import barGraph from "@/components/state/barGraph.vue";
/* ignore */

export default {
  name: "agency",
  currentYear: 2019,
  components: {
    barGraph
  },
  props: {
    agencyName: {
      type: String,
      default: ""
    }
  },
  data() {
    return {
      data: {},
      series: [],
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
    async fetchData() {
      const response = await this.$http.get(
        "/orgao/totais/TJPB/" + this.currentYear
      );
      this.data = response.data;
      this.generateSeries();
    },
    generateSeries() {
      let others = this.data.MonthTotals.map(month => month["Others"]);
      let wages = this.data.MonthTotals.map(month => month["Wage"]);
      let perks = this.data.MonthTotals.map(month => month["Perks"]);
      this.series = [
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
    },
    async nextYear() {
      await this.$http
        .get("/orgao/totais/" + this.agencyName + "/" + this.currentYear + 1)
        .then(response => (this.data = response.data));
    },
    async previousYear() {
      await this.$http
        .get("/orgao/totais/" + this.agencyName + "/" + this.currentYear - 1)
        .then(response => (this.data = response.data));
    }
  },
  async mounted() {
    this.fetchData();
  }
};
</script>

<style scoped>
.agencyName {
  font-family: "Montserrat", sans-serif;
  font-size: 25 px;
  line-height: 40px;
  padding-left: 25px;
}
.button {
  background-color: #4caf50; /* Green */
  border: none;
  color: white;
  text-decoration: none;
  font-size: 30px;
  position: relative;
  top: 10px;
  width: 50px;
}
.buttonContainer {
  width: 200px;
  height: auto;
  margin: 0 auto;
  padding: 10px;
  position: relative;
}
.agencyContainer {
  border: 1px solid firebrick;
}
</style>
