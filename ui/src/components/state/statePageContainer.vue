<template>
  <div style="min-height: 500px; text-align: center;">
    <b-dropdown split size="lg" text="Paraíba" class="dropDownButton">
    </b-dropdown>
    <agency
      v-for="(agency, i) in agencies"
      :agency="agency"
      :key="i"
      :year="new Date().getFullYear()"
    />
  </div>
</template>

<script>
import agency from "@/components/state/agency.vue";

export default {
  name: "statePageContainer",
  components: {
    agency,
  },
  data() {
    return {
      flagUrl:
        "https://1.bp.blogspot.com/-422XO8VbnkM/WFwr1v6yeoI/AAAAAAACRBM/0wtdW0JfArwQQMucxHxRrLSoHTsy7_6OwCEw/s1600/paraibano%2B2%2Bbandeira.png",
      stateName: "PARAÍBA",
      agencies: [],
    };
  },
  methods: {
    async fetchData() {
      const { data } = await this.$http.get("/orgao/PB");
      this.stateData = data;
      this.agencies = this.stateData.Agency;
    },
  },
  mounted() {
    this.fetchData();
  },
  head: {
    title: function() {
      return {
        inner: "DadosJusBr",
        complement: this.stateName,
      };
    },
    meta: function() {
      return [
        {
          name: "description",
          content:
            "DadosJusBr é uma plataforma que realiza a libertação continua de dados de remuneração do sistema de justiça brasileiro. Esta página mostra dados do estado" +
            this.stateName,
          id: "desc",
        },
      ];
    },
  },
};
</script>

<style scoped>
.dropDownButton {
  width: 284px;
  height: 59px;
  background-color: #545454;
  margin-top: 50px;
  margin-bottom: 50px;
}
</style>
