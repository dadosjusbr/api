<template>
  <div class="container">
    <h1 class="agency Name text-center">{{ this.agencyName.toUpperCase() }}</h1>
    <h2 class="entityName text-center">{{ this.year }}</h2>
    <bar-graph class="graph" :options="chartOptions" :series="series" />
  </div>
</template>

<script>
import barGraph from "@/components/state/barGraph.vue";

export default {
  data() {
    return {
      agencyName: this.$route.params.agencyName,
      year: this.$route.params.year,
      chartOptions: {
        events: {
          markerClick: function(
            event,
            chartContext,
            { seriesIndex, dataPointIndex, config }
          ) {
            alert("oiiii  ");
          },
        },

        colors: ["#c9e4ca", "#87bba2", "#364958", "#000000"],
        chart: {
          stacked: true,
          toolbar: {
            show: false,
          },
          zoom: {
            enabled: true,
          },
        },
        responsive: [
          {
            breakpoint: 601,
            options: {
              legend: {
                position: "bottom",
                offsetX: -10,
                offsetY: 0,
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
                    cssClass: "apexcharts-yaxis-label",
                  },
                  formatter: function(value) {
                    return "R$ " + (value / 1000000).toFixed(1) + "M";
                  },
                },
              },
              xaxis: {
                labels: {
                  rotate: -45,
                  rotateAlways: true,
                },
              },
            },
          },
        ],
        plotOptions: {
          bar: {
            horizontal: false,
          },
        },
        yaxis: {
          decimalsInFloat: 2,
          title: {
            text: "Total remunerções",
            offsetY: 10,
            style: {
              fontSize: "14px",
              fontWeight: "bold",
              fontFamily: undefined,
              color: "#263238",
            },
          },
          labels: {
            show: true,
            minWidth: 0,
            maxWidth: 160,
            style: {
              colors: [],
              fontSize: "16px",
              fontFamily: "Helvetica, Arial, sans-serif",
              fontWeight: 600,
              cssClass: "apexcharts-yaxis-label",
            },
            formatter: function(value) {
              if (value == 5000321) return "Não existem dados para esse mês";
              return "R$ " + (value / 1000000).toFixed(1) + "M";
            },
          },
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
            "DEZ",
          ],
          title: {
            text: "Meses",
            offsetX: 6,
            style: {
              fontSize: "15px",
              fontWeight: "bold",
              fontFamily: undefined,
              color: "#263238",
            },
          },
        },
        legend: {
          position: "right",
          offsetY: 120,
        },
        fill: {
          opacity: 1,
        },
        dataLabels: {
          enabled: false,
        },
      },
      data: {},
      series: []
    };
  },
  async mounted() {
    const { data } = await this.$http.get(
      "/orgao/totais/PB/" + this.agencyName + "/" + this.year
    );
    this.data = data;
    this.generateSeries()
  },
  methods: {
    generateSeries() {
      if (this.data.MonthTotals.length != 12) {
        this.addMonthsWithNoValue();
      }
      let others = this.data.MonthTotals.map((month) => month["Others"]);
      let wages = this.data.MonthTotals.map((month) => month["Wage"]);
      let perks = this.data.MonthTotals.map((month) => month["Perks"]);
      let noDataMarker = [];
      wages.forEach((wage) => {
        if (wage === 0) {
          noDataMarker.push(5000321);
        } else {
          noDataMarker.push(0);
        }
      });
      this.series = [
        {
          name: "Outros",
          data: others,
        },
        {
          name: "Indenizações",
          data: perks,
        },
        {
          name: "Salário",
          data: wages,
        },
        {
          name: "Sem dados",
          data: noDataMarker,
        },
      ];
    },
    addMonthsWithNoValue() {
      var existingMonths = new Array();
      this.data.MonthTotals.forEach((monthTotal) => {
        existingMonths.push(monthTotal.Month);
      });
      for (let i = 1; i <= 12; i++) {
        if (!existingMonths.includes(i)) {
          this.data.MonthTotals.push({
            Month: i,
            Others: 0,
            Perks: 0,
            Wage: 0,
          });
        }
      }
      this.data.MonthTotals.sort((a, b) => {
        return a.Month - b.Month;
      });
    },
  },
  components: {
    barGraph,
  },
};
</script>

<style scoped>
.agencyName {
  font-size: 3rem;
  margin-top: 2%;
  margin-bottom: 0%;
}
.container {
  text-align: center;
  align-items: center;
  margin-top: 1%;
}
.entityName {
  font-size: 1.5em;
  margin-top: 2%;
  font-weight: bold;
}
.graph {
  align-self: center;
  width: 100%;
}
</style>
