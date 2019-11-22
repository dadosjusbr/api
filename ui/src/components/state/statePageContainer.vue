<template>
  <div>
    <div>{{ this.stateData.Agency }}</div>
    <div>{{ this.mAgencies }}</div>
    <div>{{ this.jAgencies }}</div>
    <h1>{{ this.stateName }}</h1>
    <img :src="this.flagUrl" width="10%" />
    <entity :entityName="'Ministério Público'" :agencies="mAgencies" />
    <entity :entityName="'Judiciário'" :agencies="jAgencies" />
  </div>
</template>

<script>
import entity from "@/components/state/entity.vue";

export default {
  name: "statePageContainer",
  components: {
    entity
  },
  data() {
    return {
      flagUrl:
        "https://1.bp.blogspot.com/-422XO8VbnkM/WFwr1v6yeoI/AAAAAAACRBM/0wtdW0JfArwQQMucxHxRrLSoHTsy7_6OwCEw/s1600/paraibano%2B2%2Bbandeira.png",
      stateName: "Paraíba - PB",
      stateData: {}
    };
  },
  computed: {
    jAgencies() {
      let jAgencies = [];
      this.stateData.Agency.forEach(function(agency) {
        if (agency.AgencyCategory == "J") {
          jAgencies.push(agency.Name);
        }
      });
      return jAgencies;
    },
    mAgencies() {
      let mAgencies = [];
      this.stateData.Agency.forEach(function(agency) {
        if (agency.AgencyCategory == "M") {
          mAgencies.push(agency.Name);
        }
      });
      return mAgencies;
    }
  },
  mounted() {
    this.$http
      .get("/entidades/resumo/PB")
      .then(response => (this.stateData = response.data));
  }
};
</script>

<style scoped></style>
