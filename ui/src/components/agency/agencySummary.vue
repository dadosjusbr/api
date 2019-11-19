<template>
  <div class="summaryContainer">
    <h1>{{ agencyName }}</h1>
    <div class="cards">
      <info-card v-for="(info, id) in infos" :key="id" :info="info" />
    </div>
  </div>
</template>

<script>
import infoCard from "@/components/agency/infoCard.vue";

export default {
  name: "agencySummary",
  components: {
    infoCard
  },
  props: {
    agencySummary: {
      type: Object,
      default: null
    },
    agencyName: {
      type: String,
      default: ""
    }
  },
  computed: {
    infos() {
      return Object.values(this.agencySummary);
    }
  }
};
</script>

<style>
.cards {
  display: flex;
  flex-direction: row;
  border: #82b4b4 2px solid;
  opacity: 10;
  width: 1700px;
  text-align: center;
}

.summaryContainer {
  text-align: center;
}
h1 {
  color: #111;
  font-family: "Open Sans Condensed", sans-serif;
  font-size: 50px;
  font-weight: 700;
  line-height: 64px;
  margin: 0 0 0;
  padding: 20px 30px;
  text-align: center;
  text-transform: uppercase;
}
</style>
