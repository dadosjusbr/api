<template>
  <div>
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
      stateData: {},
      jAgencies: [],
      mAgencies: []
    };
  },
  methods: {
    async fetchData() {
      const { data } = await this.$http.get("/entidades/resumo/PB");
      this.stateData = data;
      this.setjAgencies(data);
      this.setmAgencies(data);
    },
    setjAgencies(stateData) {
      let jAgencies = [];
      if (stateData !== {}) {
        stateData.Agency.forEach(agency => {
          if (agency.AgencyCategory == "J") {
            jAgencies.push(agency.Name);
          }
        });
      }
      this.jAgencies = jAgencies;
    },
    setmAgencies(stateData) {
      let mAgencies = [];
      if (stateData !== {}) {
        stateData.Agency.forEach(agency => {
          if (agency.AgencyCategory == "M") {
            mAgencies.push(agency.Name);
          }
        });
      }
      this.mAgencies = mAgencies;
    }
  },
  mounted() {
    this.fetchData();
  }
};
</script>

<style scoped>
.stateName {
  font-family: "Montserrat", sans-serif;
  font-size: 50px;
  line-height: 40px;
  padding-left: 15px;
  float: left;
}
.image {
  width: 7%;
  margin-top: 5px;
  margin-left: 30px;
}
.header {
  border: 2px solid #6a757a;
  margin-left: 200px;
  margin-right: 200px;
  margin-top: 10px;
}
</style>
