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
  padding-left: 4%;
  margin-top: 2%;
  font-weight: bold;
}

.entity {
  margin-left: 200px;
  margin-right: 200px;
  margin-top: 10px;
}
</style>