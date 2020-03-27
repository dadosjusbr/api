<template>
  <div class="graphContainer">
    <div class="buttonContainer">
      <button v-on:click="previousMonth()" class="button btn btn-dark">
        &#8249;
      </button>
      <a> {{ this.months[this.currentMonthAndYear.month] }} </a>
      <button v-on:click="nextMonth()" class="button btn btn-dark">
        &#8250;
      </button>
    </div>
    <graph-bar
      width="100%"
      type="scatter"
      :options="chartOptions"
      :series="series"
    ></graph-bar>
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
      salaryData: [],
      currentMonthAndYear: { year: 2019, month: 1 },
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
      let year, month;
      if (this.currentMonthAndYear.month === 12) {
        year = this.currentMonthAndYear.year + 1;
        month = 1;
      } else {
        year = this.currentMonthAndYear;
        month = this.currentMonthAndYear.month + 1;
      }
      this.currentMonthAndYear = { year, month };
      this.$http
        .get("/orgao/salario/" + this.agencyName + "/" + year + "/" + month)
        .then(response => this.generateSeries(response.data));
    },
    previousMonth() {
      var year, month;
      if (this.currentMonthAndYear.month === 1) {
        year = this.currentMonthAndYear.year - 1;
        month = 12;
      } else {
        year = this.currentMonthAndYear.year;
        month = this.currentMonthAndYear.month - 1;
      }
      this.currentMonthAndYear = { year, month };
      this.$http
        .get("/orgao/salario/" + this.agencyName + "/" + year + "/" + month)
        .then(response => this.generateSeries(response.data));
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
.button {
  background-color: #182825; /* Green */
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
.graphContainer {
  margin-top: 5px;
  text-align: center;
  overflow: hidden;
}
a {
  font-family: "Montserrat", sans-serif;
  font-size: 14px;
  color: black;
}
</style>
