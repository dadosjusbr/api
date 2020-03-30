<template>
  <div class="agencyContainer">
    <div class="resume">
     <md-card>
      <md-card-content>
      </md-card-content>
    </md-card>
    </div>
    
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
      <md-button  v-if="checkPreviousYear"
        v-on:click="previousYear()">
        <img src="../../assets/previous.png">
      </md-button>
      <md-button class="deactivatedButton" v-else><img src="../../assets/previousd.png"></md-button>
      <a class="year"> {{ this.currentYear }} </a>
      <md-button
        v-if="checkNextYear"
        v-on:click="nextYear()"
      >
      <img src="../../assets/next.png">
      </md-button>
      <md-button class="deactivatedButton" v-else><img src="../../assets/nextd.png"></md-button>
    </div>
    <bar-graph class="graph" :options="chartOptions" :series="series" />
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
      currentYear: new Date().getFullYear(),
      data: {},
      series: [],
      chartOptions: {
        colors: ["#c9e4ca", "#87bba2", "#364958"],
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
            breakpoint: 601,
            options: {
              legend: {
                position: "bottom",
                offsetX: -10,
                offsetY: 0
              },
              yaxis: {
                decimalsInFloat: 2,
                labels: {
                  show: true,
                  minWidth: 0,
                  maxWidth: 50,
                  style: {
                    colors: [],
                    fontSize: "12px",
                    fontFamily: "Helvetica, Arial, sans-serif",
                    fontWeight: 600,
                    cssClass: "apexcharts-yaxis-label"
                  },
                  formatter: function(value) {
                    return "R$ " + (value / 1000000).toFixed(1) + "M";
                  }
                }
              },
              xaxis: {
                labels: {
                  rotate: -45,
                  rotateAlways: true
                }
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
              return "R$ " + (value / 1000000).toFixed(1) + "M";
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
      if (this.currentYear >= new Date().getFullYear()) {
        return false;
      } else {
        return true;
      }
    },
    checkPreviousYear() {
      if (this.currentYear <= 2018) {
        return false;
      } else {
        return true;
      }
    }
  },
  methods: {
    async fetchData() {
      var response = await this.$http.get(
        "/orgao/totais/PB/" + this.agencyName + "/" + this.currentYear
      );
      while (response.data.MonthTotals == null) {
        this.currentYear -= 1;
        response = await this.$http.get(
          "/orgao/totais/PB/" + this.agencyName + "/" + this.currentYear
        );
      }
      this.data = response.data;
      this.generateSeries();
    },
    generateSeries() {
      if (this.data.MonthTotals.length != 12) {
        this.addMonthsWithNoValue();
      }
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
    addMonthsWithNoValue() {
      var existingMonths = new Array();
      this.data.MonthTotals.forEach(monthTotal => {
        existingMonths.push(monthTotal.Month);
      });
      for (let i = 1; i <= 12; i++) {
        if (!existingMonths.includes(i)) {
          this.data.MonthTotals.push({
            Month: i,
            Others: 0,
            Perks: 0,
            Wage: 0
          });
        }
      }
      this.data.MonthTotals.sort((a, b) => {
        return a.Month - b.Month;
      });
    },
    async nextYear() {
      this.currentYear = this.currentYear + 1;
      let resp = await this.$http.get(
        "/orgao/totais/PB/" + this.agencyName + "/" + this.currentYear
      );
      if (resp.data.MonthTotals == null) {
        alert("Não existem dados disponíveis para o ano: " + this.currentYear);
        this.currentYear = this.currentYear - 1;
      } else {
        this.data = resp.data;
        this.generateSeries();
      }
    },
    async previousYear() {
      this.currentYear = this.currentYear - 1;
      let resp = await this.$http.get(
        "/orgao/totais/PB/" + this.agencyName + "/" + this.currentYear
      );
      if (resp.data.MonthTotals == null) {
        alert("Não existem dados disponíveis para o ano: " + this.currentYear);
        this.currentYear = this.currentYear + 1;
      } else {
        this.data = resp.data;
        this.generateSeries();
      }
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
  font-size: 1em;
  font-weight: bold;
}

.agencyName {
  font-size: 1.5em;
  margin-left: 15%;
  font-weight: bold;
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
  float: left;
  width: 17%;
  height: 30em;
  padding: 1px;
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
  background-color: white;
  border: none;
  color: white;
}

.md-card {
  width: 12%;
  margin-left: 2%;
  background-color: #2ab38b;
  height: 32em;
  border-style: solid;
  float: left;
}

.agencyYear {
  float:left;
}

.md-button {
  min-width: 0%;
  height: 36px;
  margin-top: -3%;
}

.graph {
  float: right;
}

@media only screen and (max-width: 379px) {

  .buttonContainer {
    float: left;
    width: 60%;
    height: 0em;
    padding: 0px; 
    position: relative;
    margin-left: 40%;
    margin-top: -9%;
  }

  .md-card {
    display: none;
  }

  .agencyName {
    font-size: 1.1em;
    margin-top: 2%;
    margin-left: 2%;
  }
}

@media only screen and (min-width: 380px) and (max-width: 600px) {

  .buttonContainer {
    float: left;
    width: 45%;
    height: 0em;
    padding: 0px; 
    position: relative;
    margin-left: 40%;
    margin-top: -9%;
  }

  .md-card {
    display: none;
  }

  .agencyName {
    margin-top: 2%;
  }
}

@media only screen and (min-width: 601px) and  (max-width: 770px) {

  .buttonContainer {
    float: left;
    width: 30%; 
    height: 2em;
    /* padding: 1px; */
    position: relative;
    margin-left: 34%;
    margin-top: -6%;  
  }

  .md-card {
  display: none;
  width: 12%;
  margin-left: 2%;
  background-color: #2ab38b;
  height: 32em;
  border-style: solid;
  float: left;
}

  .agencyName {
    margin-left: 2%;
  } 
}

@media only screen and  (min-width: 771px) and (max-width: 1025px) {

    .buttonContainer {
      float: left;
      width: 30%; 
      height: 2em;
      /* padding: 1px; */
      position: relative;
      margin-left: 34%;
      margin-top: -4%;  
    }

}

</style>