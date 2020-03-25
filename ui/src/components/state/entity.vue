<template>
  <div class="entity">
    <h1 class="entityName">{{ this.entityName.toUpperCase() }}</h1>
    <agency v-for="(name, i) in agencies" :agencyName="name" :key="i" />
  </div>
</template>

<script>
import agency from "@/components/state/agency.vue";

export default {
  name: "entity",
  components: {
    agency
  },
  props: {
    entityName: {
      type: String,
      default: ""
    },
    agencies: {
      type: Array,
      default: function() {
        return [];
      }
    }
  }
};
</script>

<style scoped>
.entityName {
  font-size: 1.5em;
  width: 20%;
  padding-left: 1%;
  margin-top: 2%;
  font-weight: bold;
  margin-bottom: 2%;
  margin-left: 2%;
}

.entity {
  margin-left: 10%;
  margin-right: 10%;
  margin-top: 2%;
}

@media only screen and (min-width: 380px) and (max-width: 600px) {
  
  .entityName {
    text-align: center;
    font-size: 1.5rem;
    width: 80%;
    padding-left: 7%;
    margin-top: 0%;
    font-weight: bold;
    margin-bottom: 5%;
    margin-left: 5%;
  }

  .entity {
    margin-left: 2%;
    margin-right: 0%;
    margin-top: 2%;
  }


}

@media only screen and (min-width: 601px) and (max-width: 770px) {
  
  .entityName {
      font-size: 1.5em;
      width: 60%;
      padding-left: 1%;
      margin-top: 0%;
      font-weight: bold;
      margin-bottom: 3%;
      margin-left: 2%;
  }
}
</style>
