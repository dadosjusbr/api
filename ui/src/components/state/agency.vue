<template>
  <div class="agencyContainer">
    <h2 class="agencyName">
      <router-link
        :to="{
          name: 'agency',
          params: { agencyName: this.agencyName.toLowerCase() }
        }"
      >
        {{ this.agencyName }}
      </router-link>
    </h2>
    <div class="buttonContainer">
      <button
        class="button"
        v-if="checkPreviousYear"
        v-on:click="previousYear()"
      >
        &#60;
      </button>
      <button class="deactivatedButton" v-else>&#60;</button>
      <a> {{ this.currentYear }} </a>
      <button class="button" v-if="checkNextYear" v-on:click="nextYear()">
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
  font-family: "Montserrat", sans-serif;
  color: black;
}

.agencyName {
  font-family: "Montserrat", sans-serif;
  font-size: 25 px;
  line-height: 40px;
  padding-left: 25px;
}
.button {
  background-color: #182825;
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
  border: 1px solid #6a757a;
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
