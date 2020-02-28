<template>
  <div class="agencyContainer">
    <div class="agencyNameContainer">
      <h1 class="agencyName">{{ agencyName }}</h1>
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
      agencyName: this.$route.params.agencyName.toUpperCase(),
      agencySummary: null
    };
  },
  methods: {
    async fetchData() {
      const { data } = await this.$http.get("/orgao/resumo/a");
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
  font-family: "Montserrat", sans-serif;
  font-size: 50px;
  line-height: 40px;
  padding-left: 15px;
}

.agencyContainer {
  margin-left: 200px;
  margin-right: 200px;
}

.agencyNameContainer {
  border: 1px solid #6a757a;
  margin-top: 5px;
}
</style>
