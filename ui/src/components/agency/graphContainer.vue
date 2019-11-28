<template>
  <div class="graphContainer">
    <div class="buttonContainer">
      <button v-on:click="previousMonth()" class="button">&#8249;</button>
      <button v-on:click="nextMonth()" class="button">&#8250;</button>
    </div>
    <graph-point
      width="90%"
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
  border: 1px solid #6a757a;
}
</style>
