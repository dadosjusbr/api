<template>
  <div class="summaryContainer">
    <h1>{{ orgaoName }}</h1>
    <div class="cards">
      <info-card v-for="(info, id) in infos" :key="id" :info="info" />
    </div>
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
