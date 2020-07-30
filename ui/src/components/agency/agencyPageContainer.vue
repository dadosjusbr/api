<template>
  <div class="agencyContainer">
    <div class="agencyNameContainer">
      <h2 class="agencyName">
        {{ this.agencyFullName }}
        <h2>{{ "(" + this.agencyName.toUpperCase() + ")" }}</h2>
      </h2>
    </div>
    <div v-show="this.chartData.length != 0" class="buttonContainer">
      <md-button
        v-if="this.activateButton.previous"
        v-on:click="previousMonth()"
      >
        <img src="../../assets/previous.png" />
      </md-button>
      <md-button class="deactivatedButton" v-else
        ><img src="../../assets/previousd.png"
      /></md-button>
      <a>
        {{ this.months[this.month] + ", " + this.year }}
      </a>
      <md-button v-if="this.activateButton.next" v-on:click="nextMonth()">
        <img src="../../assets/next.png" />
      </md-button>
      <md-button class="deactivatedButton" v-else
        ><img src="../../assets/nextd.png"
      /></md-button>
    </div>
    <div
      v-show="
        this.Crawling_Timestamp != null &&
          this.agencySummary != null &&
          this.executorLog.cmd == ''
      "
      class="cr"
    >
      Dados Capturados em {{ Crawling_Timestamp }}.
    </div>
    <div>
      <agency-summary
        v-show="this.agencySummary != null && this.executorLog.cmd == ''"
        :agencySummary="agencySummary"
        @disable-members="disableMembers"
        @enable-members="enableMembers"
        @disable-servants="disableServants"
        @enable-servants="enableServants"
      />
    </div>
    <div v-show="this.chartData.length != 0">
      <graph-container :series="chartDataToPlot" />
    </div>
    <error-collecting-data-page
      v-show="this.executorLog.cmd != ''"
      :executorLog="this.executorLog"
    />
    <no-data-available-page
      v-show="this.executorLog.cmd == '' && this.noDataAvailable"
    />
    <div class="downloadAndShare">
      <div
        v-show="this.noDataAvailable != true && this.chartData.length != 0"
        style="order: 1; width: 150px; height: 48px;"
      >
        <a :href="this.fileUrl">
          <button class="buttonDownload">
            Baixar
          </button>
        </a>
        <!--
        <h5 v-show="this.fileHash != ''">
          <b> Hash do arquivo:</b> {{ this.fileHash }}
        </h5> -->
      </div>
      <social-media-share v-show="this.agencySummary != null" />
    </div>
  </div>
</template>

<script>
import agencySummary from "@/components/agency/agencySummary.vue";
import graphContainer from "@/components/agency/graphContainer.vue";
import socialMediaShare from "@/components/agency/socialMediaShare.vue";
import errorCollectingDataPage from "@/components/agency/errorCollectingDataPage.vue";
import noDataAvailablePage from "@/components/agency/noDataAvailablePage.vue";

const formatter = new Intl.NumberFormat("de-DE");

export default {
  name: "agencyPageContainer",
  components: {
    agencySummary,
    graphContainer,
    socialMediaShare,
    errorCollectingDataPage,
    noDataAvailablePage,
  },
  data() {
    return {
      year: parseInt(this.$route.params.year, 10),
      month: parseInt(this.$route.params.month, 10),
      activateButton: {
        previous: this.checkPreviousYear(),
        next: this.checkNextYear(),
      },
      months: {
        1: "Jan",
        2: "Fev",
        3: "Mar",
        4: "Abr",
        5: "Mai",
        6: "Jun",
        7: "Jul",
        8: "Ago",
        9: "Set",
        10: "Out",
        11: "Nov",
        12: "Dez",
      },
      fileUrl: "",
      fileHash: "",
      executorLog: { cmd: "", err: "", env: [], stdout: "" },
      noDataAvailable: false,
      agencyName: this.$route.params.agencyName,
      agencyFullName: "",
      agencySummary: null,
      chartData: [],
      chartDataToPlot: [],
      Crawling_Timestamp: null,
    };
  },
  methods: {
    getNextDate() {
      let month = this.month;
      let year = this.year;
      if (month === 12) {
        month = 1;
        year = year + 1;
      } else {
        month = month + 1;
      }
      return { month, year };
    },
    getPreviousDate() {
      let month = this.month;
      let year = this.year;
      if (month === 1) {
        month = 12;
        year = year - 1;
      } else {
        month = month - 1;
      }
      return { month, year };
    },
    disableMembers() {
      this.chartDataToPlot.splice(0, 1);
      this.chartDataToPlot.splice(0, 0, { data: [], name: "" });
    },
    enableMembers() {
      this.chartDataToPlot.splice(0, 1);
      this.chartDataToPlot.splice(0, 0, this.chartData[0]);
    },
    disableServants() {
      this.chartDataToPlot.splice(1, 1);
      this.chartDataToPlot.splice(1, 0, { data: [], name: "" });
    },
    enableServants() {
      this.chartDataToPlot.splice(1, 1);
      this.chartDataToPlot.splice(1, 0, this.chartData[1]);
    },
    async checkNextYear() {
      let activateButtonNext = true;
      let { month, year } = this.getNextDate();
      if (year != undefined) {
        await this.$http
          .get(
            "/orgao/salario/" +
              this.$route.params.agencyName +
              "/" +
              year +
              "/" +
              month
          )
          .catch((err) => {
            activateButtonNext = false;
          });
        this.activateButton.next = activateButtonNext;
      }
    },
    async checkPreviousYear() {
      let activateButtonPrevious = true;
      var { month, year } = this.getPreviousDate();
      if (year != undefined) {
        await this.$http
          .get(
            "/orgao/salario/" +
              this.$route.params.agencyName +
              "/" +
              year +
              "/" +
              month
          )
          .catch((err) => {
            activateButtonPrevious = false;
          });
        this.activateButton.previous = activateButtonPrevious;
      }
    },
    async nextMonth() {
      var { month, year } = this.getNextDate();
      this.month = month;
      this.year = year;
      this.activateButton.previous = true;
      await this.$http
        .get(
          "/orgao/salario/" +
            this.agencyName +
            "/" +
            this.year +
            "/" +
            this.month
        )
        .then((response) => {
          this.chartData = this.generateSeries(response.data);
          this.chartDataToPlot = this.generateSeries(response.data);
        })
        .then(this.fetchSummaryData())
        .then(this.checkNextYear())
        .then(
          this.$router.push({
            name: "agency",
            params: { agencyName: this.agencyName, month: month, year: year },
          })
        );
      this.$router.go();
    },
    async previousMonth() {
      var { month, year } = this.getPreviousDate();
      this.month = month;
      this.year = year;
      this.activateButton.next = true;
      await this.$http
        .get(
          "/orgao/salario/" +
            this.agencyName +
            "/" +
            this.year +
            "/" +
            this.month
        )
        .then((response) => {
          this.chartData = this.generateSeries(response.data);
          this.chartDataToPlot = this.generateSeries(response.data);
        })
        .then(this.fetchSummaryData())
        .then(this.checkPreviousYear())
        .then(
          this.$router.push({
            name: "agency",
            params: { agencyName: this.agencyName, month: month, year: year },
          })
        );
      this.$router.go();
    },
    generateSeries(data) {
      return [
        {
          name: "Membros",
          data: [
            data.Members["-1"],
            data.Members["50000"],
            data.Members["40000"],
            data.Members["30000"],
            data.Members["20000"],
            data.Members["10000"],
          ],
        },
        {
          name: "Servidores",
          data: [
            data.Servers["-1"],
            data.Servers["50000"],
            data.Servers["40000"],
            data.Servers["30000"],
            data.Servers["20000"],
            data.Servers["10000"],
          ],
        },
      ];
    },
    makeExecutorLog(procInfo) {
      this.executorLog.cmd = procInfo.cmd;
      this.executorLog.err = procInfo.stderr;
      this.executorLog.stdout = procInfo.stdout;
      var envString = "";
      procInfo.env.forEach((env) => {
        envString = envString + env + "\n";
      });
      this.executorLog.env = envString.trim();
    },
    async fetchChartData() {
      this.checkPreviousYear();
      this.checkNextYear();
      var response = await this.$http
        .get(
          "/orgao/salario/" +
            this.agencyName +
            "/" +
            this.year +
            "/" +
            this.month
        )
        .catch((err) => {
          this.noDataAvailable = true;
        });

      if (response != undefined && response.status == 206) {
        this.makeExecutorLog(response.data.ProcInfo);
        response = undefined;
      }
      if (response != undefined) {
        this.chartData = this.generateSeries(response.data);
        this.chartDataToPlot = this.generateSeries(response.data);
        this.fileUrl = response.data.PackageURL;
        this.fileHash = response.data.PackageHash;
      }
    },
    async fetchSummaryData() {
      const response = await this.$http
        .get(
          "/orgao/resumo/" +
            this.agencyName +
            "/" +
            this.year +
            "/" +
            this.month
        )
        .catch((err) => {});
      if (response != undefined) {
        this.agencyFullName = response.data.FullName;
        this.agencySummary = {
          TotalEmployees: formatter.format(
            Math.trunc(response.data.TotalMembers + response.data.TotalServants)
          ),
          TotalWage:
            "R$ " + formatter.format(response.data.TotalWage.toFixed(2)),
          TotalPerks:
            "R$ " + formatter.format(response.data.TotalPerks.toFixed(2)),
          MaxWage: "R$ " + formatter.format(response.data.MaxWage.toFixed(2)),
          TotalMembers: response.data.TotalMembers,
          TotalServants: response.data.TotalServants,
          MaxPerk: "R$ " + formatter.format(response.data.MaxPerk.toFixed(2)),
          TotalRemuneration:
            "R$ " +
            formatter.format(response.data.TotalRemuneration.toFixed(2)),
        };
        const date = new Date(response.data.CrawlingTime);
        this.Crawling_Timestamp =
          date.getUTCDay() +
          " de " +
          this.months[date.getUTCMonth()] +
          " de " +
          date.getFullYear();
      } else {
        this.agencySummary = null;
      }
    },
  },
  mounted() {
    this.fetchSummaryData();
    this.fetchChartData();
  },
  head: {
    title: function() {
      return {
        inner: "DadosJusBr",
        complement:
          this.agencyName.toUpperCase() + " " + this.month + "/" + this.year,
      };
    },
    meta: function() {
      return [
        {
          name: "description",
          content:
            "DadosJusBr é uma plataforma que realiza a libertação continua de dados de remuneração do sistema de justiça brasileiro. Esta página mostra dados do orgão" +
            this.agencyName.toUpperCase(),
          id: "desc",
        },
        // Twitter
        { name: "twitter:card", content: "summary_large_image" },
        { name: "twitter:site", content: "@dadosjusbr" },
        { name: "twitter:creator", content: "@dadosjusbr" },
        {
          name: "twitter:url",
          content:
            "https://dadosjusbr.org/orgao/" +
            this.agencyName +
            "/" +
            this.year +
            "/" +
            this.month,
        },
        { name: "twitter:title", content: "DadosJusBr" },
        {
          name: "twitter:description",
          content:
            "Descubra como foram as remunerações dos funcionários do órgão" +
            this.agencyName +
            " em " +
            this.month +
            "/" +
            this.year,
        },
        {
          name: "twitter:image",
          content: "https://dadosjusbr.org/logo.png",
        },
        {
          name: "twitter:image:alt",
          content: "logo do dadojus",
        },
      ];
    },
  },
};
</script>

<style scoped>
.agencyName {
  font-weight: bold;
}
.downloadAndShare {
  display: flex;
  flex-direction: row;
  justify-content: flex-end;
  margin-bottom: 10px;
}

.buttonContainer {
  width: 105%;
  height: 10%;
  margin-top: 62px;
  margin-left: -3%;
  text-align: center;
}

button {
  margin-top: -0.2%;
}

.agencyContainer {
  margin-left: 11%;
  margin-right: 11%;
}

.agencyNameContainer {
  margin-top: 62px;
  text-align: center;
}

.buttonDownload {
  width: 150px;
  height: 48px;
  background-color: #545454;
  border: solid #545454;
  color: white;
  font-size: 17px;
}
.cr {
  text-align: center;
  font-size: 1.1em;
}
@media only screen and (max-width: 650px) {
  .agencyContainer {
    margin-left: 3%;
    margin-right: 3%;
  }
}

@media only screen and (max-width: 379px) {
  .agencyName {
    margin-left: 0%;
    margin-top: 0%;
  }

  .cr {
    text-align: center;
    font-size: 0.85em;
  }

  .agencyContainer {
    margin-left: 1%;
    margin-right: 1%;
  }
}

@media only screen and (min-width: 380px) and (max-width: 600px) {
  .agencyName {
    margin-left: 0%;
    margin-top: 0%;
  }

  .agencyContainer {
    margin-left: 3%;
    margin-right: 3%;
  }
}
</style>
