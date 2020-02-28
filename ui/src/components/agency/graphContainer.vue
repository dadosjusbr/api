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
    <graph-point
      width="100%"
      type="scatter"
      :options="chartOptions"
      :series="series"
    ></graph-point>
  </div>
</template>

<script>
import graphPoint from "@/components/agency/graphPoint.vue";

export default {
  name: "graphContainer",
  components: {
    graphPoint
  },
  data: function() {
    return {
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
        tooltip: {
          custom: function({ series, seriesIndex, dataPointIndex }) {
            return (
              '<div class="arrow_box">' +
              "<span>" +
              series[seriesIndex][dataPointIndex] +
              "</span>" +
              "</div>"
            );
          },
          colors: ["#00AEEF"]
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
        .get("/orgao/salario/TJPB/" + year + "/" + month)
        .then(response => (this.salaryData = response.data));
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
        .get("/orgao/salario/TJPB/" + year + "/" + month)
        .then(response => (this.salaryData = response.data));
    }
  },
  computed: {
    series: function() {
      let dataToPlot = this.salaryData.map((employee, index) => [
        employee["Total"],
        index + 1
      ]);
      return [{ name: "total", data: dataToPlot }];
    },
    names: function() {
      return this.salaryData.map(employee => employee["Name"]);
    },
    wages: function() {
      return this.salaryData.map(employee => employee["Wage"]);
    },
    others: function() {
      return this.salaryData.map(employee => employee["Others"]);
    },
    perks: function() {
      return this.salaryData.map(employee => employee["Perks"]);
    }
  },
  mounted() {
    this.$http
      .get(
        "/orgao/salario/TJPB/" +
          this.currentMonthAndYear.year +
          "/" +
          this.currentMonthAndYear.month
      )
      .then(response => (this.salaryData = response.data));
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
