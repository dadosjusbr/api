<template>
  <div class="graphContainer">
    <div class="buttonContainer">
      <md-button v-if="checkPreviousYear" v-on:click="previousMonth()">
        <img src="../../assets/previous.png" />
      </md-button>
      <md-button class="deactivatedButton" v-else
        ><img src="../../assets/previousd.png"
      /></md-button>
      <a>
        {{
          this.months[this.currentMonthAndYear.month] +
            ", " +
            this.currentMonthAndYear.year
        }}
      </a>
      <md-button v-if="checkNextYear" v-on:click="nextMonth()">
        <img src="../../assets/next.png" />
      </md-button>
      <md-button class="deactivatedButton" v-else
        ><img src="../../assets/nextd.png"
      /></md-button>
    </div>
    <graph-bar :options="chartOptions" :series="series"></graph-bar>
  </div>
</template>

<script>
import graphBar from "@/components/agency/graphBar.vue";

export default {
  name: "graphContainer",
  components: {
    graphBar
  },
  data: function() {
    return {
      agencyName: this.$route.params.agencyName,
      series: [],
      months: {
        1: "Janeiro",
        2: "Fevereiro",
        3: "Março",
        4: "Abril",
        5: "Maio",
        6: "Junho",
        7: "Julho",
        8: "Agosto",
        9: "Setembro",
        10: "Outubro",
        11: "Novembro",
        12: "Dezembro"
      },
      currentMonthAndYear: { year: 2020, month: 1 },
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
            horizontal: true
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
            }
          }
        },
        xaxis: {
          categories: [
            "> R$ 50 mil",
            "R$ 40~50 mil",
            "R$ 30~40 mil",
            "R$ 20~30 mil",
            "R$ 10~20 mil",
            "< R$ 10 mil"
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
  methods: {
    nextMonth() {
      let month = this.currentMonthAndYear.month;
      let year = this.currentMonthAndYear.year;
      if (this.currentMonthAndYear.month === 12) {
        month = 1;
        year = year + 1;
      } else {
        month = month + 1;
      }
      this.currentMonthAndYear = { month, year };
      this.$http
        .get("/orgao/salario/" + this.agencyName + "/" + year + "/" + month)
        .then(response => this.generateSeries(response.data));
    },
    previousMonth() {
      let month = this.currentMonthAndYear.month;
      let year = this.currentMonthAndYear.year;
      if (this.currentMonthAndYear.month === 1) {
        month = 12;
        year = year - 1;
      } else {
        month = this.currentMonthAndYear.month - 1;
      }
      this.currentMonthAndYear = { month, year };
      this.$http
        .get("/orgao/salario/" + this.agencyName + "/" + year + "/" + month)
        .then(response => this.generateSeries(response.data));
    },
    checkNextYear() {
      if (this.currentMonthAndYear.year >= new Date().getFullYear()) {
        return false;
      } else {
        return true;
      }
    },
    checkPreviousYear() {
      if (this.currentMonthAndYear.year <= 2018) {
        return false;
      } else {
        return true;
      }
    },
    generateSeries(data) {
      this.series = [
        {
          name: "Membros",
          data: [
            data.Members["-1"],
            data.Members["50000"],
            data.Members["40000"],
            data.Members["30000"],
            data.Members["20000"],
            data.Members["10000"]
          ]
        },
        {
          name: "Servidores",
          data: [
            data.Servers["-1"],
            data.Servers["50000"],
            data.Servers["40000"],
            data.Servers["30000"],
            data.Servers["20000"],
            data.Servers["10000"]
          ]
        },
        {
          name: "Inativos",
          data: [
            data.Inactives["-1"],
            data.Inactives["50000"],
            data.Inactives["40000"],
            data.Inactives["30000"],
            data.Inactives["20000"],
            data.Inactives["10000"]
          ]
        }
      ];
    }
  },
  async mounted() {
    const { data } = await this.$http.get(
      "/orgao/salario/" +
        this.agencyName +
        "/" +
        this.currentMonthAndYear.year +
        "/" +
        this.currentMonthAndYear.month
    );
    this.generateSeries(data);
  }
};
</script>

<style scoped>
.buttonContainer {
  margin: 0 auto;
  width: 70%;
  height: 10%;
  position: relative;
}
.graphContainer {
  text-align: center;
  overflow: hidden;
  /* background-color:  rgb(4, 4, 173); */
  margin-bottom: 10px;
}
a {
  color: black;
  font-size: 1.4em;
  font-weight: bold;
}
</style>
