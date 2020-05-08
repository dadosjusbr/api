<template>
  <div class="agencyContainer">
    <div class="agencyNameContainer">
      <h1 class="agencyName">{{ agencyName.toUpperCase() }}</h1>
    </div>
    <div>
      <agency-summary
        v-show="this.agencySummary != null"
        :agencySummary="agencySummary"
      />
    </div>
    <div>
      <graph-container @change="date" />
    </div>
    <div
      v-show="this.Crawling_Timestamp != null && this.agencySummary != null"
      class="cr"
    >
      Dados Capturados em {{ Crawling_Timestamp | formatDate }}, horário de
      Brasília.
    </div>

    <div class="socialMidiaShare">
      {{ this.url }}
      <h5><b>Compartilhe essa informação: </b></h5>
      <facebook
        style="margin-right: 5px"
        :url="this.url"
        scale="2"
        :title="this.socialMidiaMsg"
      ></facebook>
      <whats-app
        style="margin-right: 5px"
        :url="this.url"
        :title="this.socialMidiaMsg"
        scale="2"
      ></whats-app>
      <twitter
        style="margin-right: 5px"
        :url="this.url"
        :title="this.socialMidiaMsg"
        scale="2"
      ></twitter>
      <email
        style="margin-right: 5px"
        :url="this.url"
        :subject="this.socialMidiaMsg"
        scale="2"
      ></email>
    </div>
  </div>
</template>

<script>
import { Facebook, Twitter, WhatsApp, Email } from "vue-socialmedia-share";
import agencySummary from "@/components/agency/agencySummary.vue";
import graphContainer from "@/components/agency/graphContainer.vue";

const formatter = new Intl.NumberFormat("de-DE");

export default {
  name: "agencyPageContainer",
  components: {
    agencySummary,
    graphContainer,
    Facebook,
    Twitter,
    WhatsApp,
    Email,
  },
  data() {
    return {
      socialMidiaMsg:
        "Descubra como é a distribuição das remunerações dos funcionários do " +
        this.$route.params.agencyName.toUpperCase() +
        " no ano e mês " +
        this.$route.params.year +
        "/" +
        this.$route.params.month,
      url:
        "https://dadosjusbr.org/orgao/" +
        this.$route.params.agencyName +
        "/" +
        this.$route.params.year +
        "/" +
        this.$route.params.month,
      agencyName: this.$route.params.agencyName,
      year: this.$route.params.year,
      month: this.$route.params.month,
      agencySummary: null,
      Crawling_Timestamp: null,
    };
  },
  methods: {
    date(date) {
      this.year = date.year;
      this.month = date.month;
      this.fetchData();
    },
    async fetchData() {
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
      if (response != undefined && response.data.TotalEmployees != 0) {
        this.agencySummary = {
          Total_Empregados: formatter.format(
            Math.trunc(response.data.TotalEmployees)
          ),
          Total_Salários:
            "R$ " + formatter.format(response.data.TotalWage.toFixed(2)),
          Total_Indenizações:
            "R$ " + formatter.format(response.data.TotalPerks.toFixed(2)),
          Salário_Máximo:
            "R$ " + formatter.format(response.data.MaxWage.toFixed(2)),
        };
        this.Crawling_Timestamp = response.data.CrawlingTime;
      } else {
        this.agencySummary = null;
      }
    },
  },
  mounted() {
    this.fetchData();
  },
};
</script>

<style scoped>
.agencyName {
  font-weight: bold;
}

.socialMidiaShare {
  text-align: center;
  margin-top: 10px;
  margin-bottom: 5px;
}

.agencyContainer {
  margin-left: 14%;
  margin-right: 14%;
}

.agencyNameContainer {
  padding-top: 2.4%;
  text-align: center;
}

.cr {
  text-align: center;
  font-size: 1.1em;
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
