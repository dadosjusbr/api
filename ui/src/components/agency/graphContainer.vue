<template>
  <div>
    <div>
      <button v-on:click="previousMonth()" class="button">Anterior</button>
      <button v-on:click="nextMonth()" class="button">Proximo</button>
      <div>{{ this.salaryData }}</div>
      <div>{{ this.others }}</div>
      <div>{{ this.names }}</div>
      <div>{{ this.wages }}</div>
      <div>{{ "other 0 :" + this.others[0] }}</div>
    </div>
    <graph-point
      width="500"
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
          }
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
.container {
  border: coral 2px solid;
  height: 500px;
  padding-top: 50px;
  align-self: center;
}
.button {
  background-color: #4caf50; /* Green */
  border: none;
  color: white;
  padding: 15px 32px;
  text-align: center;
  text-decoration: none;
  display: inline-block;
  font-size: 16px;
  background-color: #008cba;
  padding: 10px 24px;
}
</style>
