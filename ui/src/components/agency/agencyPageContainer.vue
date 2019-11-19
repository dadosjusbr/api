<template>
  <div>
    <agency-summary :agencyName="agencyName" :agencySummary="agencySummary" />
    <graph-container />
  </div>
</template>

<script>
import agencySummary from "@/components/agency/agencySummary.vue";
import graphContainer from "@/components/agency/graphContainer.vue";

export default {
  name: "agencyPageContainer",
  components: {
    agencySummary,
    graphContainer
  },
  data() {
    return {
      agencyName: "TJPB",
      agencySummary: ["info 1", "info2", "info 3", "info4"]
    };
  }
};
</script>

<style lang="scss" scoped></style>
