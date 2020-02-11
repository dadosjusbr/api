<template>
  <div class="cards">
    <info-card
      v-for="(value, name, id) in agencySummary"
      :key="id"
      :info="{ value, name }"
    />
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
    }
  }
};
</script>

<style>
.cards {
  justify-content: center;
  display: flex;
  border: #6a757a 1px solid;
  opacity: 10;
  width: 100%;
  align-self: center;
  align-content: center;
  margin-top: 5px;
  overflow: scroll;
}
</style>
