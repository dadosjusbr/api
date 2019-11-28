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
  line-height: 40px;
  padding-left: 15px;
}

.entity {
  border: 1px solid #6a757a;
}
</style>
