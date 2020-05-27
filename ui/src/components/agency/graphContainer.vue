<template>
  <div class="graphContainer">
    <graph-bar :options="chartOptions" :series="series"></graph-bar>
  </div>
</template>

<script>
import graphBar from "@/components/agency/graphBar.vue";

export default {
  name: "graphContainer",
  components: {
    graphBar,
  },
  props: {
    series: {
      type: Array,
      default: [],
    },
  },
  data: function() {
    return {
      chartOptions: {
        colors: ["#c9e4ca", "#87bba2", "#364958"],
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
            breakpoint: 480,
            options: {
              legend: {
                position: "bottom",
                offsetX: -10,
                offsetY: 0,
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
          title: {
            text: "Remuneração",
            offsetX: 6,
            style: {
              fontSize: "16px",
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
    };
  },
};
</script>

<style scoped>
.graphContainer {
  text-align: center;
  overflow: hidde;
  margin-bottom: 10px;
}

</style>
