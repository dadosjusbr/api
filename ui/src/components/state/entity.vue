<template>
  <div class="entity">
    <!-- <h1 class="entityName">{{ this.entityName.toUpperCase() }}</h1> -->
    <agency
      v-for="(agency, i) in agencies"
      :agency="agency"
      :key="i"
      :year="new Date().getFullYear()"
    />
  </div>
</template>

<script>
import agency from "@/components/state/agency.vue";

export default {
  name: "entity",
  components: {
    agency,
  },
  props: {
    entityName: {
      type: String,
      default: "",
    },
    agencies: {
      type: Array,
      default: function() {
        return [];
      },
    },
  },
};
</script>

<style scoped>
.entityName {
  font-size: 1.5em;
  font-weight: bold;
}
.entity {
  text-align: center;
}
</style>
