<template>
  <div class="summaryContainer">
    <h1>{{ orgaoName }}</h1>
    <info-card v-for="(info, id) in infos" :key="id" :info="info" />
  </div>
</template>

<script>
import infoCard from "@/components/orgao/infoCard.vue";

export default {
  name: "orgaoSummary",
  components: {
    infoCard
  },
  props: {
    orgaoName: {
      type: String,
      default: "Erro ao carregar o nome do orgão"
    },
    infos: {
      type: Array,
      dafault: []
    }
  },
  data() {
    return {
      info: "Hello I am A Circle"
    };
  }
};
</script>

<style></style>
