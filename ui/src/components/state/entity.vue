<template>
  <div class="entity">
    <h1 class="entityName">{{ this.entityName }}</h1>
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
  font-family: "Montserrat", sans-serif;
  font-size: 30px;
  padding-left: 10px;
}

.entity {
  border: 1px solid #6a757a;
  margin-left: 200px;
  margin-right: 200px;
  margin-top: 10px;
}
</style>
