<template>
  <div>
    <div>{{ this.stateData.Agency }}</div>
    <div>{{ this.mAgencies }}</div>
    <div>{{ this.jAgencies }}</div>
    <div class="header">
      <h1 class="stateName">{{ this.stateName }}</h1>
      <img class="image" :src="this.flagUrl" />
    </div>
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
      stateName: "Paraíba",
      stateData: {}
    };
  },
  methods: {
    async fetchData() {
      await this.$http
        .get("/entidades/resumo/PB")
        .then(response => (this.stateData = response.data));
    }
  },
  created() {
    this.fetchData();
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
  }
};
</script>

<style scoped>
.stateName {
  font-family: "Montserrat", sans-serif;
  font-size: 50px;
  line-height: 40px;
  padding-left: 15px;
}
.image {
  width: 7%;
  position: absolute;
  top: 140px;
  right: 17px;
}
.header {
  border: 2px solid #d00000;
}
</style>
