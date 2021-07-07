<template>
  <b-container fluid  style="background-color: #ffffff;" >
    <b-row class="mt-4"></b-row>
    <b-row v-show="this.chartData.length != 0" class="mt-5 mb-3 buttonContainer">
      <b-col  class="buttonContainer d-flex align-items-center justify-content-center">
      <b class="agencyName ">
        {{ this.agencyFullName + " (" + this.agencyName.toUpperCase() + ")" }}
      </b>
      <br />
      </b-col>
    </b-row>
     <b-row v-show="this.chartData.length != 0" class="buttonContainer mb-3">
      <b-col  class="buttonContainer d-flex align-items-center justify-content-center">
      <md-button
        v-if="this.activateButton.previous"
        v-on:click="previousMonth()"
      >
        <img
          style="height: 30px; width: 30px;"
          src="../../assets/previous.svg"
        />
      </md-button>
      <md-button class="deactivatedButton" v-else
        ><img style="height: 30px; width:30px;" src="../../assets/previousd.png"
      /></md-button>
      <a>
        {{ this.months[this.month] + ", " + this.year }}
      </a>
      <md-button v-if="this.activateButton.next" v-on:click="nextMonth()">
        <img style="height: 30px; width: 30px;" src="../../assets/next.svg" />
      </md-button>
      <md-button class="deactivatedButton" v-else
        ><img style="height: 30px; width:30px;" src="../../assets/nextd.svg"
      /></md-button>
       </b-col >
    </b-row>
    <b-row
      v-if="
        this.Crawling_Timestamp != null &&
          this.agencySummary != null &&
          this.executorLog.cmd == ''
      "
      class="cr  d-flex align-items-center justify-content-center mb-3"
    >
      Dados Capturados em {{ Crawling_Timestamp }}.
    </b-row>
    <b-row>
      <b-col cols="1"></b-col>
      <b-col>
        <agency-summary
          v-if="this.agencySummary != null && this.executorLog.cmd == ''"
          :agencySummary="agencySummary"
        />
      </b-col>
      <b-col cols="1"></b-col>
    </b-row>
    <b-row>
      <b-col cols="1"></b-col>
      <b-col >
        <graph-container
          :series="chartData"
          :date="{ month: this.months[this.month], year: this.year }"
          v-show="this.chartData.length != 0"
        />
      </b-col>
      <b-col cols="1"></b-col>
    </b-row>
    <error-collecting-data-page
      v-if="this.executorLog.cmd != ''"
      :executorLog="this.executorLog"
    />
    <no-data-available-page
      v-if="this.executorLog.cmd == '' && this.noDataAvailable"
    />
    <b-row class="buttonContainer2 mb-5 ml-1 ">
      <b-col cols="1" class="d-none d-xl-block"></b-col>
      <b-col cols="12" xl="3" class="mb-2" >
        <router-link to="/dados">
          <img
            style="height: 83px; width:295px"
            src="../../assets/button-explorar-anos.svg"
          />
        </router-link>
      </b-col>
      <b-col cols="0" xl="1"></b-col>
      <b-col cols="12" xl="3" class="mb-2">
        <social-media-share v-if="this.agencySummary != null" />
      </b-col>
      <b-col cols="12" xl="3" class="mb-2">
        <a :href="this.fileUrl">
          <img
            style="height: 83px; width:295px"
            src="../../assets/button-baixar.svg"
          />
        </a>
      </b-col>
      <b-col cols="1"></b-col>
    </b-row>
  </b-container>
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
          TotalWage:
            "R$ " + formatter.format(response.data.TotalWage.toFixed(2)),
          TotalPerks:
            "R$ " + formatter.format(response.data.TotalPerks.toFixed(2)),
          MaxWage: "R$ " + formatter.format(response.data.MaxWage.toFixed(2)),
          TotalMembers: response.data.TotalMembers,
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
.buttonContainer2 {
  height: 82px;
}

.agencyName {
  color: #3e5363;
  font-size: 1.4em;
}

.cr {
  text-align: center;
  font-size: 1.0em;
  color: #3e5363;
}

.buttonContainer {
  color: #3e5363;
  font-size: 1.3em;
}


@media only screen and (max-width: 700px) {
  .buttonContainer2 {
  height: 82px;
  margin-bottom: 66% !important;
  }
  .agencyName {
    color: #3e5363;
    font-size: 1.2em;
  }

}

</style>
