<template>
  <div class="agencyContainer">
    <h2 class="agencyName">
      <router-link
        :to="{
          name: 'agency',
          params: { agencyName: this.agencyName.toLowerCase() }
        }"
      >
        {{ this.agencyName.toUpperCase() }}
      </router-link>
    </h2>
    <div class="buttonContainer">
      <button
        class="button btn btn-dark"
        v-if="checkPreviousYear"
        v-on:click="previousYear()"
      >
        &#60;
      </button>
      <button class="deactivatedButton" v-else>&#60;</button>
      <a class="year"> {{ this.currentYear }} </a>
      <button
        class="button btn btn-dark"
        v-if="checkNextYear"
        v-on:click="nextYear()"
      >
        &#62;
      </button>
      <button class="deactivatedButton" v-else>&#62;</button>
    </div>
    <bar-graph :options="chartOptions" :series="series" />
  </div>
</template>

<script>
import barGraph from "@/components/state/barGraph.vue";
/* ignore */

export default {
  name: "agency",
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
      currentYear: 2019,
      data: {},
      series: [],
      chartOptions: {
        colors: ["#991040", "#F9CD30", "#00AEEF"],
        chart: {
          stacked: true,
          toolbar: {
            show: false
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
        yaxis: {
          decimalsInFloat: 2,
          labels: {
            show: true,
            minWidth: 0,
            maxWidth: 160,
            style: {
              colors: [],
              fontSize: "16px",
              fontFamily: "Helvetica, Arial, sans-serif",
              fontWeight: 600,
              cssClass: "apexcharts-yaxis-label"
            },
            formatter: function(value) {
              return value / 1000000 + "M R$";
            }
          }
        },
        xaxis: {
          categories: [
            "JAN",
            "FEV",
            "MAR",
            "ABR",
            "MAI",
            "JUN",
            "JUL",
            "AGO",
            "SET",
            "OUT",
            "NOV",
            "DEZ"
          ]
        },
        legend: {
          position: "right",
          offsetY: 120
        },
        fill: {
          opacity: 1
        },
        dataLabels: {
          enabled: false
        }
      }
    };
  },
  computed: {
    checkNextYear() {
      if (this.currentYear >= 2020) {
        return false;
      } else {
        return true;
      }
    },
    checkPreviousYear() {
      if (this.currentYear <= 2015) {
        return false;
      } else {
        return true;
      }
    }
  },
  methods: {
    async fetchData() {
      const response = await this.$http.get(
        "/orgao/totais/PB/" + this.agencyName + "/" + this.currentYear
      );
      this.data = response.data;
      this.generateSeries();
    },
    generateSeries() {
      let others = this.data.MonthTotals.map(month => month["Others"]);
      // eslint-disable-next-line
      console.log(others);
      let wages = this.data.MonthTotals.map(month => month["Wage"]);
      let perks = this.data.MonthTotals.map(month => month["Perks"]);
      this.series = [
        {
          name: "Outros",
          data: others
        },
        {
          name: "Indenizações",
          data: perks
        },
        {
          name: "Remunerações",
          data: wages
        }
      ];
    },
    async nextYear() {
      this.currentYear = this.currentYear + 1;
      await this.$http
        .get("/orgao/totais/" + this.agencyName + "/" + this.currentYear)
        .then(response => (this.data = response.data));
    },
    async previousYear() {
      this.currentYear = this.currentYear - 1;
      await this.$http
        .get("/orgao/totais/" + this.agencyName + "/" + this.currentYear)
        .then(response => (this.data = response.data));
    }
  },
  async mounted() {
    this.fetchData();
  }
};
</script>

<style scoped>
a {
  color: black;
}

.year {
  color: black;
  font-size: 20px;
}

.agencyName {
  font-size: 30px;
  line-height: 40px;
  padding-left: 25px;
  text-decoration: underline;
}
.button {
  background-color: #182825;
  border: none;
  color: white;
  text-decoration: none;
  font-size: 20px;
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
  border-block-end: solid 1px black;
  overflow: auto;
  margin-top: 5px;
  margin-bottom: 5px;
  margin-right: 5px;
  margin-left: 5px;
}

.deactivatedButton {
  background-color: grey; /* Green */
  border: none;
  color: white;
  text-decoration: none;
  font-size: 30px;
  position: relative;
  top: 10px;
  width: 50px;
}
</style>
