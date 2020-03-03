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
        Total_Empregados: data.TotalEmployees,
        Total_Salários: data.TotalWage,
        Total_Indenizações: data.TotalPerks,
        Salário_Maximo: data.MaxWage
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
  font-size: 5rem;
  float: left;
  margin: 30px auto;
}

.agencyContainer {
  margin-left: 200px;
  margin-right: 200px;
}

.agencyNameContainer {
  margin-top: 5px;
  height: 150px;
}
</style>
