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
        Salário_Máximo: "R$ " + formatter.format(data.MaxWage.toFixed(2))
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
  font-weight: bold;
  color: white;
}

.agencyContainer {
  margin-left: 200px;
  margin-right: 200px;
}

.agencyNameContainer {
  padding: 25px;
  border: 1px solid #666;
  margin-top: 8px;
  height: 100px;
  background-color: rgb(4, 4, 173);
  text-align: center;
}

@media only screen and (max-width: 379px) {
  .agencyName {
    margin-left: 16%;
    margin-top: 2%;
  }
}

@media only screen and (min-width: 380px) and (max-width: 600px) {
  .agencyName {
    margin-left: 16%;
    margin-top: 2%;
  }
}
</style>
