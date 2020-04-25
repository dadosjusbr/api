<template>
  <div class="container">
    <h1 class="agency Name text-center">{{ this.agencyName.toUpperCase() }}</h1>
    <h2 class="year text-center">{{ this.year }}</h2>
    <agency :agencyName="agencyName" :simplifyComponent="true" :year="year" />
  </div>
</template>

<script>
import agency from "@/components/state/agency.vue";

export default {
  data() {
    return {
      agencyName: this.$route.params.agencyName,
      year: parseInt(this.$route.params.year, 10),
    };
  },
  components: {
    agency,
  },
};
</script>

<style scoped>
.agencyName {
  font-size: 3rem;
  margin-top: 2%;
  margin-bottom: 0%;
}
.container {
  text-align: center;
  align-items: center;
  margin-top: 1%;
}
.year {
  font-size: 1.5em;
  margin-top: 2%;
  font-weight: bold;
}
</style>
