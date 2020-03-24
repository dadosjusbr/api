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

@media only screen and (max-width: 770px) {
  
  .entityName[data-v-44f4c2b1] {
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