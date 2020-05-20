<template>
  <div>
    <div v-show="this.executorLog.cmd == '' && this.noDataAvailable">
      <md-empty-state
        md-rounded
        md-icon="highlight_off"
        md-label="Talvez o órgão ainda não tenha disponibilizado os dados ou o dadosjusbr não tentou realizar a coleta."
        md-description="Acha que tem algo errado? Por favor entre em contato conosco abrindo uma issue."
      >
        <a
          style="padding-bottom: 15px; font-size: 16px;  font-weight: normal;"
          href="https://github.com/dadosjusbr/coletores/issues/new"
          target="_blank"
        >
          Abra uma issue aqui</a
        >
      </md-empty-state>
    </div>

    <div v-show="this.executorLog.cmd != ''">
      <md-empty-state
        md-rounded
        md-icon="highlight_off"
        md-label="Tivemos um erro ao coletar os dados. Veja o erro abaixo."
      >
      </md-empty-state>
    </div>

    <div
      v-show="this.noDataAvailable != true && this.series.length != 0"
      class="graphContainer"
    >
      <div class="buttonContainer">
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
      <graph-bar :options="chartOptions" :series="series"></graph-bar>
    </div>
    <div
      style="text-align: center;"
      v-show="this.noDataAvailable != true && this.series.length != 0"
    >
      <h5><b>Faça download do .csv e arquivo: </b></h5>
      <md-button
        style="margin: 0px 0px 0px 0px"
        :href="this.fileUrl"
        target="_blank"
        lass="md-icon-button md-raised"
      >
        <md-icon>cloud_download</md-icon>
      </md-button>
      <h5 v-show="this.fileHash != ''">
        <b> Hash do arquivo:</b> {{ this.fileHash }}
      </h5>
    </div>
    <div v-show="this.executorLog.cmd != ''" class="errorLog">
      <b>Erro no comando: </b>
      <br />
      <textarea wrap="soft" class="textArea">
 {{ this.executorLog.cmd }} </textarea
      >
      <br />
      <b> Saída de erro: </b>
      <br />
      <textarea class="textArea">  {{ this.executorLog.err }}</textarea>
      <br />
      <b v-show="this.executorLog.stdout != ''"> Saída padrão: </b>

      <textarea class="textArea" v-show="this.executorLog.stdout != ''">
        {{ this.executorLog.stdout }}
      </textarea>
      <b> Variáveis de ambiente (env): </b>
      <br />
      <textarea rows="20" class="textArea">
  {{ this.executorLog.env }}</textarea
      >
    </div>
    <div
      style="text-align: center; margin-bottom: 10px"
      v-show="this.executorLog.cmd != ''"
    >
      <h5>
        Acha que tem algo errado? Por favor entre em contato conosco abrindo uma
        issue.
      </h5>
      <a
        style="padding-bottom: 15px; font-size: 16px;  font-weight: normal;"
        href="https://github.com/dadosjusbr/coletores/issues/new"
        target="_blank"
      >
        Abra uma issue aqui</a
      >
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
  data: function() {
    return {
      noDataAvailable: null,
      agencyName: this.$route.params.agencyName,
      year: parseInt(this.$route.params.year, 10),
      month: parseInt(this.$route.params.month, 10),
      executorLog: { cmd: "", err: "", env: [], stdout: "" },
      fileUrl: "",
      fileHash: "",
      activateButton: {
        previous: this.checkPreviousYear(),
        next: this.checkNextYear(),
      },
      series: [],
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
  methods: {
    async checkNextYear() {
      let activateButtonNext = true;
      let { month, year } = this.getNextDate();
      if (year != undefined) {
        await this.$http
          .get("/orgao/salario/" + this.agencyName + "/" + year + "/" + month)
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
          .get("/orgao/salario/" + this.agencyName + "/" + year + "/" + month)
          .catch((err) => {
            activateButtonPrevious = false;
          });
        this.activateButton.previous = activateButtonPrevious;
      }
    },
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
        .then((response) => this.generateSeries(response.data))
        .then(this.checkNextYear())
        .then(this.$emit("change", { year, month }))
        .then(
          this.$router.push({
            name: "agency",
            params: { agencyName: this.agencyName, month: month, year: year },
          })
        );
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
        .then((response) => this.generateSeries(response.data))
        .then(this.checkPreviousYear())
        .then(this.$emit("change", { year, month }))
        .then(
          this.$router.push({
            name: "agency",
            params: { agencyName: this.agencyName, month: month, year: year },
          })
        );
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
    generateSeries(data) {
      this.fileUrl = data.PackageURL;
      this.noDataAvailable = data.PackageHash;
      this.series = [
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
        {
          name: "Inativos",
          data: [
            data.Inactives["-1"],
            data.Inactives["50000"],
            data.Inactives["40000"],
            data.Inactives["30000"],
            data.Inactives["20000"],
            data.Inactives["10000"],
          ],
        },
      ];
    },
  },
  async mounted() {
    var response = await this.$http
      .get(
        "/orgao/salario/" + this.agencyName + "/" + this.year + "/" + this.month
      )
      .catch((err) => {
        this.noDataAvailable = true;
      });

    if (response != undefined && response.status == 206) {
      this.makeExecutorLog(response.data.ProcInfo);
      response = undefined;
    }
    if (response != undefined) this.generateSeries(response.data);
  },
};
</script>

<style scoped>
h5 {
  margin-bottom: 0px;
}
.buttonContainer {
  width: 105%;
  height: 10%;
  margin-top: 8%;
  margin-left: -3%;
}
.graphContainer {
  text-align: center;
  overflow: hidde;
  margin-bottom: 10px;
}
a {
  color: black;
  font-size: 1.4em;
  font-weight: bold;
}

button {
  margin-top: -0.4%;
}

.errorLog {
  padding-left: 5px;
  padding-top: 5px;
  padding-right: 5px;
  padding-bottom: 5px;
  text-align: center;
  margin-bottom: 5px;
}

.textArea {
  border: 1px solid #2ab38b;
  width: 80%;
}
</style>
