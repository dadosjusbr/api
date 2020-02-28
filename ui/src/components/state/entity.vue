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
  font-size: 60px;
  padding-left: 20px;
}

.entity {
  margin-left: 200px;
  margin-right: 200px;
  margin-top: 10px;
}
</style>
