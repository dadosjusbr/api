<template>
  <div class="container">
    <h1 class="agency Name text-center">{{ this.agencyName.toUpperCase() }}</h1>
    <h2 class="year text-center">{{ this.month + " - " + this.year }}</h2>
    <agency-summary
      v-show="!this.noSummaryData"
      :agencySummary="this.agencySummary"
    />
    <graph-container
      :agencyNameSimplyComponent="this.agencyName"
      :year="this.year"
      :month="this.month"
      :simplifyComponent="true"
    />
  </div>
</template>

<script>
import graphContainer from "@/components/agency/graphContainer.vue";
import agencySummary from "@/components/agency/agencySummary.vue";
const formatter = new Intl.NumberFormat("de-DE");

export default {
  data() {
    return {
      agencyName: this.$route.params.agencyName,
      year: parseInt(this.$route.params.year, 10),
      month: parseInt(this.$route.params.month, 10),
      agencySummary: null,
      noSummaryData: false,
    };
  },
  components: {
    graphContainer,
    agencySummary,
  },

  methods: {
    async fetchData() {
      const response = await this.$http.get(
        "/orgao/resumo/" + this.agencyName + "/" + this.year + "/" + this.month
      ).catch((err) => {});
      if (response == undefined) {
        //eslint-ignore
        console.log("oii");

        this.noSummaryData = true;
        return;
      }
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
    },
  },
  mounted() {
    this.fetchData();
  },
};
</script>

<style>
.agencyName {
  font-size: 3rem;
  margin-top: 2%;
  margin-bottom: 0%;
}
.container {
  text-align: center;
  align-items: center;
  margin-top: 1%;
}
.year {
  font-size: 1.5em;
  margin-top: 2%;
  font-weight: bold;
}
</style>
