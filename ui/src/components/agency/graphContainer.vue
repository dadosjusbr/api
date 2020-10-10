<template>
  <b-container fluid class="graphContainer">
    <b-row class="graphHeader">
      <b-col class="d-flex align-items-center justify-content-center">
        Remuneração por quantidade de empregados em {{ this.date.month }} de
        {{ this.date.year }}
      </b-col>
    </b-row>
    <b-row >
      <b-col cols="12" class="graph d-flex align-items-center justify-content-center">
        <apexcharts
          width="100%"
          height="500"
          type="bar"
          :options="chartOptions"
          :series="series"
        />
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
import VueApexCharts from "vue-apexcharts";

export default {
  name: "graphContainer",
  components: {
    apexcharts: VueApexCharts,
  },
  props: {
    series: {
      type: Array,
      default: [],
    },
    date: {
      type: Object,
      default: null,
    },
  },
  data: function() {
    return {
      chartOptions: {
        legend: {
          show: false,
        },
        colors: ["#c9a0d0", "#513658"],
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
            breakpoint: 500,
            options: {
              chart: {
                width: "95%",
              },
              yaxis: {
                labels: {
                  maxWidth: 120,
                  style: {
                    colors: [],
                    fontSize: "12px",
                    fontFamily: "Helvetica, Arial, sans-serif",
                    fontWeight: 600,
                    cssClass: "apexcharts-yaxis-label",
                  },
                },
              },
            },
          },
        ],
        plotOptions: {
          bar: {
            horizontal: true,
            barHeight: "70%",
          },
        },
        yaxis: {
          decimalsInFloat: 2,
          labels: {
            show: true,
            minWidth: 0,
            maxWidth: 160,
            style: {
              colors: [],
              fontSize: "14px",
              fontFamily: "Helvetica, Arial, sans-serif",
              fontWeight: 600,
              cssClass: "apexcharts-yaxis-label",
            },
          },
        },
        xaxis: {
          categories: [
            "> R$ 50 mil",
            "R$ 40~50 mil",
            "R$ 30~40 mil",
            "R$ 20~30 mil",
            "R$ 10~20 mil",
            "< R$ 10 mil",
          ],
          title: {
            text: "Quantidade de funcionários",
            margin: 10,
            style: {
              fontSize: "16px",
              fontWeight: "bold",
              fontFamily: undefined,
              color: "#263238",
            },
          },
        },
        fill: {
          opacity: 1,
        },
        dataLabels: {
          enabled: false,
        },
      },
    };
  },
};
</script>

<style scoped>
.graphTooltip {
  height: 65px;
  width: 144px;
  background-color: white;
}

.graphContainer {
  text-align: center;
  background-color: white;
  text-align: center;
}

.graphHeader {
  background-color: #3e5363;
  font-size: 1.5em;
  color: #ffffff;
  height: 5em !important;
}

.graph {
  margin-bottom: 5%;
  height: 100%;
}

@media only screen and (max-width: 700px) {
  .graphHeader {
    height: 7em !important;
  }
}
</style>
