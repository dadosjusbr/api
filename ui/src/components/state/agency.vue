<template>
  <div>
    <md-empty-state
      v-show="noDataAvailable"
      md-rounded
      md-icon="highlight_off"
      md-label="Não existem dados para esse ano :("
      md-description="Talvez o órgão não tenha disponibilizado os dados em seu site."
    >
    </md-empty-state>

    <div class="agencyContainer" v-show="!noDataAvailable">
      <h2 class="agencyName">
        <router-link
          :to="{
            name: 'agency',
            params: {
              agencyName: this.agency.Name.toLowerCase(),
              year: this.yearWithData,
              month: this.monthWithData,
            },
          }"
        >
          {{ this.agency.Name.toUpperCase() + " - " + this.agency.FullName }}
        </router-link>
      </h2>

      <div class="buttonContainer" v-show="!simplifyComponent">
        <md-button v-if="checkPreviousYear" v-on:click="previousYear()">
          <img src="../../assets/previous.png" />
        </md-button>
        <md-button class="deactivatedButton" v-else
          ><img src="../../assets/previousd.png"
        /></md-button>
        <a class="year"> {{ this.currentYear }} </a>
        <md-button v-if="checkNextYear" v-on:click="nextYear()">
          <img src="../../assets/next.png" />
        </md-button>
        <md-button class="deactivatedButton" v-else
          ><img src="../../assets/nextd.png"
        /></md-button>
      </div>
      <div class="agencyContent">
        <div style="width: 90%; align-self: center;">
          <div class="remunerationMenu">
            <div class="menuHeader">
              <div style="width: 90%">
                <p>
                  Total Remunerações {{ this.currentYear }}: R$
                  {{ this.totals.totalRemuneration }}M
                </p>
              </div>
              <div style="width: 5%">
                <md-icon id="tooltip-target-1">info</md-icon>
              </div>
            </div>
            <div class="employeesClassification" style="padding-top: 15px">
              <div class="employeeClass">
                <div
                  style="background-color: #364958;"
                  :class="[!this.dataFilter.wage ? 'squareOpac' : '', 'square']"
                  v-on:click="filterWage()"
                ></div>
                <p>Salario: {{ this.totals.wageTotal }}M</p>
              </div>
              <div class="employeeClass">
                <div
                  style="background-color: #87bba2;"
                  :class="[
                    !this.dataFilter.perks ? 'squareOpac' : '',
                    'square',
                  ]"
                  v-on:click="filterPerks()"
                ></div>
                <p>Indenizações: {{ this.totals.perksTotal }}M</p>
              </div>
              <div class="employeeClass">
                <div
                  style="background-color: #c9e4ca;"
                  :class="[
                    !this.dataFilter.others ? 'squareOpac' : '',
                    'square',
                  ]"
                  v-on:click="filterOthers()"
                ></div>
                <p>Outros: {{ this.totals.othersTotal }}M</p>
              </div>
              <div class="employeeClass">
                <div
                  style="background-color: #000000;"
                  :class="[
                    !this.dataFilter.noData ? 'squareOpac' : '',
                    'square',
                  ]"
                  v-on:click="filterNoData()"
                ></div>
                <p>Sem Dados</p>
              </div>
            </div>
          </div>
        </div>
        <div class="graphContainer">
          <div
            style="height: 59px;background-color: rgba(155, 155, 155, 0.4); line-height: 57px;"
          >
            <p>Soma do valor de remunerações por mês</p>
          </div>
          <div style="margin-left: 20%;">
            <bar-graph :options="chartOptions" :series="chartDataToPlot" />
          </div>
        </div>
        <div style="width: 90%; align-self: center; text-align: right;">
          <button v-on:click="routerToOMA()" class="moreInfoButton">
            Mais informação
            <router-link
              :to="{
                name: 'agency',
                params: {
                  agencyName: this.agency.Name.toLowerCase(),
                  year: this.yearWithData,
                  month: this.monthWithData,
                },
              }"
            >
            </router-link>
          </button>
        </div>
      </div>
    </div>
    <b-tooltip target="tooltip-target-1" triggers="hover">
      - Salário: valor recebido de acordo com a prestação de serviços, em
      decorrência do contrato de trabalho.
      <br />
      - Outros: Qualquer remuneração recebida por um funcionário que não seja
      proveniente de salário ou indenizações. Exemplos de outros são: diárias,
      gratificações, remuneração por função de confiança, benefícios pessoais ou
      eventuais...
      <br />
      - Indenizações: valores eventuais, por exemplo, auxílios alimentação,
      saúde, escolar...
    </b-tooltip>
  </div>
</template>

<script>
import barGraph from "@/components/state/barGraph.vue";

export default {
  name: "agency",
  components: {
    barGraph,
  },
  props: {
    agency: {
      type: Object,
      default: null,
    },
    simplifyComponent: {
      type: Boolean,
      default: false,
    },
    year: {
      type: Number,
      default: new Date().getFullYear(),
    },
  },
  data() {
    return {
      totals: {
        totalRemuneration: 0,
        totalWage: 0,
        totalPerks: 0,
        totalOthers: 0,
      },
      dataFilter: {
        wage: true,
        perks: true,
        others: true,
        noData: true,
      },
      monthWithData: 0,
      yearWithData: 0,
      currentYear: this.year,
      noDataAvailable: false,
      data: {},
      series: [],
      chartDataToPlot: [],
      chartOptions: {
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
              if (value == 29000321) return "Não existem dados para esse mês";
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
          show: false,
          position: "right",
          offsetY: 120,
        },
        fill: {
          opacity: 1,
          image: {
            src: [
              "https://catalogue.accasoftware.com/img/Prodotti/2920/PREVIEW/hachura-30.1.750x527-1_1563779607.PNG",
            ],
          },
        },
        dataLabels: {
          enabled: false,
        },
      },
    };
  },
  computed: {
    checkNextYear() {
      if (this.currentYear >= new Date().getFullYear()) {
        return false;
      } else {
        return true;
      }
    },
    checkPreviousYear() {
      if (this.currentYear <= 2018) {
        return false;
      } else {
        return true;
      }
    },
  },
  methods: {
    filterWage() {
      if (this.dataFilter.wage) {
        this.chartDataToPlot.splice(2, 1);
        this.chartDataToPlot.splice(2, 0, { data: [], name: "" });
      } else {
        this.chartDataToPlot.splice(2, 1);
        this.chartDataToPlot.splice(2, 0, this.series[2]);
      }
      this.dataFilter.wage = !this.dataFilter.wage;
    },
    filterPerks() {
      if (this.dataFilter.perks) {
        this.chartDataToPlot.splice(1, 1);
        this.chartDataToPlot.splice(1, 0, { data: [], name: "" });
      } else {
        this.chartDataToPlot.splice(1, 1);
        this.chartDataToPlot.splice(1, 0, this.series[1]);
      }
      this.dataFilter.perks = !this.dataFilter.perks;
    },
    filterOthers() {
      if (this.dataFilter.others) {
        this.chartDataToPlot.splice(0, 1);
        this.chartDataToPlot.splice(0, 0, { data: [], name: "" });
      } else {
        this.chartDataToPlot.splice(0, 1);
        this.chartDataToPlot.splice(0, 0, this.series[0]);
      }
      this.dataFilter.others = !this.dataFilter.others;
    },
    filterNoData() {
      if (this.dataFilter.noData) {
        this.chartDataToPlot.splice(3, 1);
        this.chartDataToPlot.splice(3, 0, { data: [], name: "" });
      } else {
        this.chartDataToPlot.splice(3, 1);
        this.chartDataToPlot.splice(3, 0, this.series[3]);
      }
      this.dataFilter.noData = !this.dataFilter.noData;
    },
    routerToOMA() {
      this.$router.push({
        name: "agency",
        params: {
          agencyName: this.agency.Name.toLowerCase(),
          year: this.yearWithData,
          month: this.monthWithData,
        },
      });
    },
    async fetchData() {
      var response = await this.$http.get(
        "/orgao/totais/PB/" + this.agency.Name + "/" + this.currentYear
      );
      if (this.simplifyComponent == true && response.data.MonthTotals == null) {
        this.noDataAvailable = true;
      } else {
        while (response.data.MonthTotals == null) {
          this.currentYear -= 1;
          response = await this.$http.get(
            "/orgao/totais/PB/" + this.agency.Name + "/" + this.currentYear
          );
        }
      }
      this.data = response.data;
      this.yearWithData = this.currentYear;
      this.monthWithData = response.data.MonthTotals.length;
      this.sumTotals();
      this.generateSeries();
    },
    sumTotals() {
      let remunerationTotal = 0;
      let wageTotal = 0;
      let perksTotal = 0;
      let othersTotal = 0;
      this.data.MonthTotals.forEach((month) => {
        remunerationTotal =
          remunerationTotal + month.Wage + month.Others + month.Perks;
        wageTotal = wageTotal + month.Wage;
        perksTotal = perksTotal + month.Perks;
        othersTotal = othersTotal + month.Others;
      });
      this.totals.totalRemuneration = (remunerationTotal / 1000000).toFixed(1);
      this.totals.wageTotal = (wageTotal / 1000000).toFixed(0);
      this.totals.perksTotal = (perksTotal / 1000000).toFixed(0);
      this.totals.othersTotal = (othersTotal / 1000000).toFixed(0);
    },
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
          noDataMarker.push(29000321);
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
      this.chartDataToPlot = [
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
    async nextYear() {
      this.currentYear = this.currentYear + 1;
      let resp = await this.$http.get(
        "/orgao/totais/PB/" + this.agency.Name + "/" + this.currentYear
      );
      if (resp.data.MonthTotals == null) {
        alert("Não existem dados disponíveis para o ano: " + this.currentYear);
        this.currentYear = this.currentYear - 1;
      } else {
        this.data = resp.data;
        this.generateSeries();
        this.sumRemuneration();
      }
    },
    async previousYear() {
      this.currentYear = this.currentYear - 1;
      let resp = await this.$http.get(
        "/orgao/totais/PB/" + this.agency.Name + "/" + this.currentYear
      );
      if (resp.data.MonthTotals == null) {
        alert("Não existem dados disponíveis para o ano: " + this.currentYear);
        this.currentYear = this.currentYear + 1;
      } else {
        this.data = resp.data;
        this.generateSeries();
        this.sumRemuneration();
      }
    },
  },
  async mounted() {
    this.fetchData();
  },
};
</script>

<style scoped>
.menuHeader {
  width: 100%;
  height: 59px;
  background-color: rgba(155, 155, 155, 0.4);
  line-height: 57px;
  text-align: center;
  display: flex;
  flex-direction: row;
}

.year {
  margin-top: 10px;
}

.agencyContent {
  display: flex;
  flex-direction: column;
}

.remunerationMenu {
  height: 150px;
  width: 400px;
  background-color: white;
  margin-bottom: 50px;
  margin-top: 50px;
}

.graphContainer {
  background-color: white;
  min-width: 90%;
  align-self: center;
}

.agencyName {
  font-size: 20px;
  font-weight: 900;
  align-self: center;
  font-size: 20px;
}

a {
  color: #4a4a4a;
}

.buttonContainer {
  margin-top: 15px;
  width: 100%;
  height: 27px;
  text-align: center;
}
.md-button {
  margin-top: -3px;
}

.moreInfoButton {
  margin: 15px 0 15px 0px;
  width: 150px;
  height: 48px;
  background-color: #7f3d8b;
  border: solid #7f3d8b;
  color: white;
  font-size: 17px;
}

.agencyContainer {
  min-height: 900px;
  padding-top: 50px;
  margin-bottom: 30px;
  margin-right: 5px;
  margin-left: 5px;
  background-color: rgba(155, 155, 155, 0.2);
}
</style>
