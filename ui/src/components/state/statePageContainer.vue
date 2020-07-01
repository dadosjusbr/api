<template>
  <div style="min-height: 500px; text-align: center;">
    <b-dropdown split size="lg" text="Paraíba" class="dropDownButton">
    </b-dropdown>
    <entity :entityName="'Ministério Público'" :agencies="mAgencies" />
    <entity :entityName="'Judiciário'" :agencies="jAgencies" />
  </div>
</template>

<script>
import entity from "@/components/state/entity.vue";

export default {
  name: "statePageContainer",
  components: {
    entity,
  },
  data() {
    return {
      flagUrl:
        "https://1.bp.blogspot.com/-422XO8VbnkM/WFwr1v6yeoI/AAAAAAACRBM/0wtdW0JfArwQQMucxHxRrLSoHTsy7_6OwCEw/s1600/paraibano%2B2%2Bbandeira.png",
      stateName: "PARAÍBA",
      stateData: {},
      jAgencies: [],
      mAgencies: [],
    };
  },
  methods: {
    async fetchData() {
      const { data } = await this.$http.get("/orgao/PB");
      this.stateData = data;
      this.setjAgencies(data);
      this.setmAgencies(data);
    },
    setjAgencies(stateData) {
      let jAgencies = [];
      if (stateData !== {}) {
        stateData.Agency.forEach((agency) => {
          if (agency.AgencyCategory == "J") {
            jAgencies.push(agency);
          }
        });
      }
      this.jAgencies = jAgencies;
    },
    setmAgencies(stateData) {
      let mAgencies = [];
      if (stateData !== {}) {
        stateData.Agency.forEach((agency) => {
          if (agency.AgencyCategory == "M") {
            mAgencies.push(agency);
          }
        });
      }
      this.mAgencies = mAgencies;
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
