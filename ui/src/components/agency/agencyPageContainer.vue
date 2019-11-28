<template>
  <div>
    <h1 class="agencyName">{{ agencyName }}</h1>
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
        SalárioMaximo: data.MaxWage
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
</style>
