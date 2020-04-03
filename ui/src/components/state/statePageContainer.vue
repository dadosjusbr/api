<template>
  <div>
    <div class="header">
      <h1 class="stateName text-left">{{ this.stateName }}</h1>
      <img class="image rounded float-left" :src="this.flagUrl" />
    </div>

    <md-card :class="[hiddenCard ? 'hiddenCard' : 'card']">
      <md-card-content>
        Os gráficos onde as barras não aparecerem é por que não existe dado para
        aquele mês. É provável que o órgão não tenha disponibilizado o dado em
        sua página.
      </md-card-content>
      <md-button v-on:click="hideCard()"> X </md-button>
    </md-card>

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
      stateName: "PARAÍBA",
      stateData: {},
      jAgencies: [],
      mAgencies: [],
      hiddenCard: false
    };
  },
  methods: {
    hideCard() {
      this.hiddenCard = true;
      //es-lint ignore next line
      console.log("apertou no x");
    },
    async fetchData() {
      const { data } = await this.$http.get("/orgao/PB");
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
.card {
  height: 230px;
  width: 200px;
  top: 35%;
  right: 10%;
  text-align: justify;
  position: fixed;
  background-color: #c9e4ca;
  font-weight: 500;
}

.hiddenCard {
  display: none;
}

.stateName {
  font-size: 3rem;
  float: left;
  margin-top: 2%;
  margin-bottom: 0%;
}
.image {
  width: 5%;
  height: 80x;
  margin-left: 2%;
  margin-top: 2.5%;
}
.header {
  margin-left: 10%;
  margin-top: 1%;
  height: 5em;
}

@media only screen and (max-width: 379px) {
  .stateName {
    font-size: 2.2rem;
    float: left;
    margin-top: 4%;
    margin-left: 13%;
    margin-bottom: 0%;
  }

  .image {
    width: 12%;
    margin-left: 6%;
    margin-top: 6%;
  }

  .header {
    margin-left: 10%;
    margin-top: 1%;
    height: 4em;
  }
}

@media only screen and (min-width: 380px) and (max-width: 600px) {
  .stateName {
    font-size: 2.5rem;
    float: left;
    margin-top: 4%;
    margin-left: 12%;
    margin-bottom: 0%;
  }

  .image {
    width: 11%;
    margin-left: 6%;
    margin-top: 6%;
  }
}

@media only screen and (min-width: 601px) and (max-width: 770px) {
  .stateName {
    font-size: 3rem;
    float: left;
    margin-top: 2%;
    margin-left: 3%;
    margin-bottom: 0%;
  }

  .image {
    width: 10%;
  }
}

@media only screen and (min-width: 771px) and (max-width: 1025px) {
  .image {
    width: 6%;
    height: 80x;
    margin-left: 2%;
    margin-top: 3%;
  }
}

@media only screen and (max-width: 600px) {
  .card {
    height: 230px;
    width: 200px;
    top: 60%;
    right: 20%;
    text-align: justify;
    position: fixed;
    background-color: #c9e4ca;
    font-weight: 500;
  }
}
</style>
