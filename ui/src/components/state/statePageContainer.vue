<template>
  <div>
    <div class="header">
      <h1 class="stateName text-left">{{ this.stateName }}</h1>
      <img class="image rounded float-left" :src="this.flagUrl" />
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
      const { data } = await this.$http.get("entidades/resumo/a");
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
  font-size: 5rem;
  float: left;
  margin: 30px auto;
}
.image {
  width: 200px;
  height: 130px;
  margin: 9px auto;
  margin-left: 20px;
}
.header {
  margin-left: 200px;
  margin-right: 200px;
  margin-top: 10px;
  height: 150px;
}
</style>
