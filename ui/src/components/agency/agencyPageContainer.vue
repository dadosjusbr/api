<template>
  <div class="agencyContainer">
    <div class="agencyNameContainer">
      <h1 class="agencyName">{{ agencyName.toUpperCase() }}</h1>
    </div>
    <agency-summary :agencySummary="agencySummary" />
    <graph-container />
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
    graphContainer
  },
  data() {
    return {
      agencyName: this.$route.params.agencyName,
      agencySummary: null
    };
  },
  methods: {
    async fetchData() {
      const { data } = await this.$http.get("/orgao/resumo/" + this.agencyName);
      this.agencySummary = {
        Total_Empregados: formatter.format(Math.trunc(data.TotalEmployees)),
        Total_Salários: "R$ " + formatter.format(data.TotalWage.toFixed(2)),
        Total_Indenizações:
          "R$ " + formatter.format(data.TotalPerks.toFixed(2)),
        Salário_Maximo: "R$ " + formatter.format(data.MaxWage.toFixed(2))
      };
    }
  },
  mounted() {
    this.fetchData();
  }
};
</script>

<style scoped>
.agencyName {
  font-size: 3.5rem;
  float: left;
  margin: 30px auto;
}

.agencyContainer {
  margin-left: 200px;
  margin-right: 200px;
}

.agencyNameContainer {
  margin-top: 1px;
  height: 100px;
}
</style>
