<template>
  <div class="agencyContainer">
    <div class="agencyNameContainer">
      <h1 class="agencyName">{{ agencyName.toUpperCase() }}</h1>
    </div>
    <div
      v-show="this.Crawling_Timestamp != null && this.agencySummary != null"
      class="cr"
    >
      Dados Capturados em {{ Crawling_Timestamp | formatDate }}, horário de
      Brasília.
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
    <social-media-share v-show="this.agencySummary != null" />
  </div>
</template>

<script>
import agencySummary from "@/components/agency/agencySummary.vue";
import graphContainer from "@/components/agency/graphContainer.vue";
import socialMediaShare from "@/components/agency/socialMediaShare.vue";

const formatter = new Intl.NumberFormat("de-DE");

export default {
  name: "agencyPageContainer",
  components: {
    agencySummary,
    graphContainer,
    socialMediaShare,
  },
  data() {
    return {
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
        { name: "twitter:title", content: "DadosJusBr" },
        {
          name: "twitter:image",
          content: "https://dadosjusbr.org/img/white_logo.16edf55b.png",
        },
        { name: "twitter:card", content: "summary" },
        { name: "twitter:site", content: "@dadosjusbr" },
        {
          name: "twitter:description",
          content:
            "Descubra como é a distribuição das remunerações dos funcionários do sistema judiciário ",
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
