<template>
  <div class="graphContainer">
    <div class="graphHeader">Remuneração por quantidade de empregados</div>
    <div style="graph">
      <graph-bar :options="chartOptions" :series="series"></graph-bar>
    </div>
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
        legend:{
          show: false
        },
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
.graphContainer {
  text-align: center;
  overflow: hidde;
  margin-bottom: 10px;
  background-color: white;
  height: 526px;
  text-align: center;
}

.graphHeader {
  width: 100%;
  height: 59px;
  box-shadow: 20px 20px 300px 0 rgba(0, 0, 0, 0.05);
  background-color: rgba(155, 155, 155, 0.3);
  line-height: 59px;
}

.graph{
  padding-left: 14%;
  height: 100%;
  padding-top: 23px;
}
</style>
