<template>
  <div>
    <entity-summary :entityName="entityName" :entitySummary="entitySummary" />
    <graph-container />
  </div>
</template>

<script>
import entitySummary from "@/components/entity/entitySummary.vue";
import graphContainer from "@/components/entity/graphContainer.vue";

export default {
  name: "entityPageContainer",
  components: {
    entitySummary,
    graphContainer
  },
  data() {
    return {
      entityName: "TJPB",
      entitySummary: ["info 1", "info2", "info 3", "info4"]
    };
  }
};
</script>

<style lang="scss" scoped></style>
