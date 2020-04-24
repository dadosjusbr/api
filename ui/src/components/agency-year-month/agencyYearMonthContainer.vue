<template>
  <div class="container">
    <h1 class="agency Name text-center">{{ this.agencyName.toUpperCase() }}</h1>
    <h2 class="year text-center">{{this.month + " - " + this.year }}</h2>

    <graph-container
      :agencyNameSimplyComponent="this.agencyName"
      :year="this.year"
      :month="this.month"
      :simplifyComponent="true"
    />
  </div>
</template>

<script>
import graphContainer from "@/components/agency/graphContainer.vue";

export default {
  data() {
    return {
      agencyName: this.$route.params.agencyName,
      year: parseInt(this.$route.params.year, 10),
      month: parseInt(this.$route.params.month, 10),
    };
  },
  components: {
    graphContainer,
  },
};
</script>

<style>
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
