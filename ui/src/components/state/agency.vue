<template>
  <b-container fluid>
    <md-empty-state
      v-show="noDataAvailable"
      md-rounded
      md-icon="highlight_off"
      md-label="Não existem dados para esse ano :("
      md-description="Talvez o órgão não tenha disponibilizado os dados em seu site."
    >
    </md-empty-state>

    <b-container fluid class="p-xl-5 p-0" v-show="!noDataAvailable">
      <b-row class="mt-3 mb-2 align-items-center">
        <b-col>
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
          <b class="agencyName">
            {{
              this.agency.FullName + " (" + this.agency.Name.toUpperCase() + ")"
            }}
          </b>
        </router-link>
         </b-col>
      </b-row >
      <b-row class="buttonContainer align-items-center justify-content-center mt-3 mb-2" v-show="!simplifyComponent"> 
          <md-button v-if="checkPreviousYear" v-on:click="previousYear()">
            <img
              style="height: 30px; width:30px"
              src="../../assets/previous.svg"
            />
          </md-button>
          <md-button class="deactivatedButton" v-else
            ><img src="../../assets/previousd.png"
          /></md-button>
          <a class="year"> {{ this.currentYear }} </a>
          <md-button v-if="checkNextYear" v-on:click="nextYear()">
            <img style="height: 30px; width:30px" src="../../assets/next.svg" />
          </md-button>
          <md-button class="deactivatedButton" v-else
            ><img style="height: 30px; width:30px;" src="../../assets/nextd.svg"
          /></md-button>
      </b-row>
      
      <b-row class="menuHeader " >
        <b-col class="pl-1 mt-xl-3">
        <p>
          Total de Remunerações em {{ this.currentYear }}: R$
          {{ this.totals.totalRemuneration }}M
        </p>
        </b-col>
        <b-col cols="1" class="mt-xl-3" :id="this.agency.Name">
          <img style="width: 30%;" src="../../assets/icon-info.svg" />
        </b-col>
      <b-col cols="6" class="d-none d-xl-block" style="background-color: white;"> </b-col>
      </b-row>
      <b-row class="remunerationMenu" style="color: #ffffff;">
          <b-col class="employeeClass mt-2">
            <div
              style="background-color: #364958;"
              :class="[!this.dataFilter.wage ? 'squareOpac' : '', 'square']"
              v-on:click="filterWage()"
            >
              <img
                style="width: 100%; height: 100%;"
                src="../../assets/icon-salario-oma.svg"
              />
            </div>
            <p>Salario: {{ this.totals.totalWage }}M</p>
          </b-col>
          
          <b-col class="employeeClass mt-2">
            <div
              style="background-color: #c9e4ca;"
              :class="[
                !this.dataFilter.benefits ? 'squareOpac' : '',
                'square',
              ]"
              v-on:click="filterBenefits()"
            >
              <img
                style="width: 100%; height: 100%;"
                src="../../assets/icon-beneficio-oma.svg"
              />
            </div>
            <p>Benefícios: {{ this.totals.totalBenefits }}M</p>
          </b-col>
          <b-col class="employeeClass mt-2">
            <div
              style="background-color: #000000;"
              :class="[
                !this.dataFilter.noData ? 'squareOpac' : '',
                'square',
              ]"
              v-on:click="filterNoData()"
            ></div>
            <p>Sem Dados</p>
          </b-col>
        <b-col cols="6" class="d-none d-xl-block" style="background-color: white;"> </b-col>
      </b-row>
      <b-row class="graphHeader mt-5">
        <b-col class="mt-xl-3 align-items-center justify-content-center">
          <p>Total de Remunerações por Mês em {{ this.year }}</p>
        </b-col>
      </b-row >
      
      <b-row class="auxDivGraph">
        <b-col cols="2" class="d-none d-xl-block"></b-col>
        <b-col><bar-graph :options="chartOptions" :series="chartDataToPlot" /></b-col>
        <b-col cols="2" class="d-none d-xl-block"></b-col>
      </b-row>
        <b-row class="ml-1 justify-content-xl-end" style="width: 100%; align-self: center; text-align: right;">
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
            <img
              style="width: 295px; height: 83px;"
              src="../../assets/button_explorar_meses.svg"
            />
          </router-link>
        </b-row>
    </b-container>
    <b-tooltip :target="this.agency.Name" triggers="hover">
      - Salário: valor recebido de acordo com a prestação de serviços, em
      decorrência do contrato de trabalho.
      <br />
      - Benefícios: Qualquer remuneração recebida por um funcionário que não
      seja proveniente de salário. Exemplos de benefícios são: diárias,
      gratificações, remuneração por função de confiança, benefícios pessoais ou
      eventuais, auxílios alimentação, saúde, escolar...
    </b-tooltip>
    <b-row class="mt-0" >
        <b-col style="border-top: 1px solid black;"></b-col>
  </b-row>
  </b-container>
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
      maxMonth: 0,
      totals: {
        totalRemuneration: 0,
        totalWage: 0,
        totalBenefits: 0,
      },
      dataFilter: {
        wage: true,
        benefits: true,
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
        colors: ["#c9e4ca", "#364958", "#000000"],
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
              legend: {
                position: "bottom",
                offsetX: -10,
                offsetY: 0,
              },
              chart: {
                width: "100%",
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
            text: "Total de Remunerações",
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
        this.chartDataToPlot.splice(1, 1);
        this.chartDataToPlot.splice(1, 0, { data: [], name: "" });
      } else {
        this.chartDataToPlot.splice(1, 1);
        this.chartDataToPlot.splice(1, 0, this.series[1]);
      }
      this.dataFilter.wage = !this.dataFilter.wage;
    },
    filterBenefits() {
      if (this.dataFilter.benefits) {
        this.chartDataToPlot.splice(0, 1);
        this.chartDataToPlot.splice(0, 0, { data: [], name: "" });
      } else {
        this.chartDataToPlot.splice(0, 1);
        this.chartDataToPlot.splice(0, 0, this.series[0]);
      }
      this.dataFilter.benefits = !this.dataFilter.benefits;
    },
    filterNoData() {
      if (this.dataFilter.noData) {
        this.chartDataToPlot.splice(2, 1);
        this.chartDataToPlot.splice(2, 0, { data: [], name: "" });
      } else {
        this.chartDataToPlot.splice(2, 1);
        this.chartDataToPlot.splice(2, 0, this.series[2]);
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
        "/orgao/totais/" + this.agency.Name + "/" + this.currentYear
      );
      if (this.simplifyComponent == true && response.data.MonthTotals == null) {
        this.noDataAvailable = true;
      } else {
        while (response.data.MonthTotals == null && this.currentYear > 2018) {
          this.currentYear -= 1;
          response = await this.$http.get(
            "/orgao/totais/" + this.agency.Name + "/" + this.currentYear
          );
        }
      }
      this.data = response.data;
      this.yearWithData = this.currentYear;
      this.monthWithData =
        response.data.MonthTotals[response.data.MonthTotals.length - 1].Month;
      this.sumTotals();
      this.generateSeries();
    },
    sumTotals() {
      let maxMonth = 0;
      let remunerationTotal = 0;
      let totalWage = 0;
      let perksTotal = 0;
      let othersTotal = 0;
      this.data.MonthTotals.forEach((month) => {
        let monthSum = month.Wage + month.Others + month.Perks;
        if (maxMonth < monthSum)
          maxMonth = month.Wage + month.Others + month.Perks;
        remunerationTotal = remunerationTotal + monthSum;
        totalWage = totalWage + month.Wage;
        perksTotal = perksTotal + month.Perks;
        othersTotal = othersTotal + month.Others;
      });
      this.maxMonth = maxMonth;
      this.totals.totalRemuneration = (remunerationTotal / 1000000).toFixed(1);
      this.totals.totalWage = (totalWage / 1000000).toFixed(0);
      this.totals.totalBenefits = (
        (perksTotal + othersTotal) /
        1000000
      ).toFixed(0);
    },
    generateSeries() {
      if (this.data.MonthTotals.length != 12) {
        this.addMonthsWithNoValue();
      }
      let wages = this.data.MonthTotals.map((month) => month["Wage"]);
      let benefits = this.data.MonthTotals.map(
        (month) => month["Perks"] + month["Others"]
      );
      let noDataMarker = [];
      wages.forEach((wage) => {
        if (wage === 0) {
          noDataMarker.push(this.maxMonth);
        } else {
          noDataMarker.push(0);
        }
      });
      this.series = [
        {
          name: "Benefícios",
          data: benefits,
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
          name: "Benefícios",
          data: benefits,
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
            benefits: 0,
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
        "/orgao/totais/" + this.agency.Name + "/" + this.currentYear
      );
      if (resp.data.MonthTotals == null) {
        alert("Não existem dados disponíveis para o ano: " + this.currentYear);
        this.currentYear = this.currentYear - 1;
      } else {
        this.data = resp.data;
        this.yearWithData = this.currentYear;
        this.monthWithData = resp.data.MonthTotals.length;
        this.generateSeries();
        this.sumRemuneration();
      }
    },
    async previousYear() {
      this.currentYear = this.currentYear - 1;
      let resp = await this.$http.get(
        "/orgao/totais/" + this.agency.Name + "/" + this.currentYear
      );
      if (resp.data.MonthTotals == null) {
        alert("Não existem dados disponíveis para o ano: " + this.currentYear);
        this.currentYear = this.currentYear + 1;
      } else {
        this.data = resp.data;
        this.yearWithData = this.currentYear;
        this.monthWithData = resp.data.MonthTotals.length;
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
.agencyName {
  color: #3e5363;
  font-size: 1.3rem;
}

.menuHeader {
  background-color: #3e5363;
  border-bottom: 1px solid #ffffff;
  font-size: 1.3rem;
  color: #ffffff;
}

.graphHeader {
  background-color: #3e5363;
  color: white;
  font-size: 1.3rem;
}

.remunerationMenu {
  background-color: #3e5363;
}

.graphContainer {
  background-color: white;
  align-self: center;
}

a {
  color: #4a4a4a;
}

.buttonContainer {
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



</style>
