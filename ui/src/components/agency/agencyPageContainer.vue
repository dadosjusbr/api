<template>
  <div class="agencyContainer">
    <div class="agencyNameContainer">
      <h1 class="agencyName">{{ agencyName.toUpperCase() }}</h1>
    </div>
    <div>
      <agency-summary v-show="this.agencySummary != null" :agencySummary="agencySummary" />
    </div>
    <div>
      <graph-container @change="date" />
    </div>
    <div v-show="this.Crawling_Timestamp != null" class="cr">
      Dados Capturados em {{ Crawling_Timestamp | formatDate }}, horário de
      Brasília.
    </div>
  </div>
</template>

<script>
import agencySummary from "@/components/agency/agencySummary.vue";
import graphContainer from "@/components/agency/graphContainer.vue";

const formatter = new Intl.NumberFormat("de-DE");

export default {
  name: "agencyPageContainer",
  components: {
    agencySummary,
    graphContainer,
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
      const { data } = await this.$http.get(
        "/orgao/resumo/" + this.agencyName + "/" + this.year + "/" + this.month
      );
      this.agencySummary = {
        Total_Empregados: formatter.format(Math.trunc(data.TotalEmployees)),
        Total_Salários: "R$ " + formatter.format(data.TotalWage.toFixed(2)),
        Total_Indenizações:
          "R$ " + formatter.format(data.TotalPerks.toFixed(2)),
        Salário_Máximo: "R$ " + formatter.format(data.MaxWage.toFixed(2)),
      };
      this.Crawling_Timestamp = data.CrawlingTime;
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
