<template>
  <div class="container">
    <h1 class="stateName text-center">{{ this.agencyName.toUpperCase() }}</h1>
    <h2 class="entityName text-center">{{ this.year + " - " + this.month }}</h2>
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
      month: this.$route.params.month,
      series: [
        {
          name: "Outros",
          data: [
            5378919.910000007,
            5416046.700000006,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0
          ]
        },
        {
          name: "Indenizações",
          data: [1517733.75, 1566334.7500000014, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
        },
        {
          name: "Salário",
          data: [
            12516638.410000043,
            12444599.890000043,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0
          ]
        },
        {
          name: "Sem dados",
          data: [
            0,
            0,
            5000321,
            5000321,
            5000321,
            5000321,
            5000321,
            5000321,
            5000321,
            5000321,
            5000321,
            5000321
          ]
        }
      ],
      chartOptions: {
        events: {
          markerClick: function(
            event,
            chartContext,
            { seriesIndex, dataPointIndex, config }
          ) {
            alert("oiiii  ");
          }
        },

        colors: ["#c9e4ca", "#87bba2", "#364958", "#000000"],
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
          title: {
            text: "Total remunerções",
            offsetY: 10,
            style: {
              fontSize: "14px",
              fontWeight: "bold",
              fontFamily: undefined,
              color: "#263238"
            }
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
              cssClass: "apexcharts-yaxis-label"
            },
            formatter: function(value) {
              if (value == 5000321) return "Não existem dados para esse mês";
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
          ],
          title: {
            text: "Meses",
            offsetX: 6,
            style: {
              fontSize: "15px",
              fontWeight: "bold",
              fontFamily: undefined,
              color: "#263238"
            }
          }
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
  components: {
    barGraph
  }
};
</script>

<style scoped>
.stateName {
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
