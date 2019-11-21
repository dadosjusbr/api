<template>
  <div>
    <h1>{{ this.stateName }}</h1>
    <img :src="this.flagUrl" width="10%" />
    <entity />
    <entity />
  </div>
</template>

<script>
import entity from "@/components/state/entity.vue";

export default {
  name: "statePageContainer",
  components: {
    entity
  },
  data() {
    return {
      flagUrl:
        "https://1.bp.blogspot.com/-422XO8VbnkM/WFwr1v6yeoI/AAAAAAACRBM/0wtdW0JfArwQQMucxHxRrLSoHTsy7_6OwCEw/s1600/paraibano%2B2%2Bbandeira.png",
      stateName: "Paraíba - PB"
    };
  }
};
</script>

<style scoped></style>