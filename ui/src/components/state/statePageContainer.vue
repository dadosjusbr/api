<template>
  <b-container fluid class="p-0">
  <b-row
    style="text-align: center; background-color: #3e5363; "
  >
  <b-col >
    <b-form-select v-on:change="fetchData" size="lg" class="dropDownButton" v-model="selected" :options="options"></b-form-select>
  </b-col>
  </b-row>
  <b-row style="text-align: center;">
    <b-col cols="1" style="background-color: #3e5363;"></b-col>
    <b-col :key="selected">
      <agency
        v-for="(agency,state) in agencies"
        :agency="agency"
        :key="state"
        :year="new Date().getFullYear()"
      />
    </b-col>
    <b-col cols="1" style="background-color: #3e5363;"></b-col>
  </b-row>
  </b-container>
</template>

<script>
import agency from "@/components/state/agency.vue";

export default {
  name: "statePageContainer",
  components: {
    agency,
  },
  data() {
    return {
      state: "PB",
      selected: "PB",
      options: [
        { value: "Federal", text: 'Órgãos Federais' },
        {
          label: "Órgãos Estaduais",
          options: [
            { value: "PB", text: "Paraíba" }
          ]
        }
      ],
      agencies: [],
    };
  },
  methods: {
    async fetchData() {
      const { data } = await this.$http.get("/orgao/" + this.selected);
      this.stateData = data;
      this.agencies = this.stateData.Agency;
      this.state = this.selected;
    },
  },
  mounted() {
    this.fetchData();
  },
  head: {
    title: function() {
      return {
        inner: "DadosJusBr",
        complement: this.selected,
      };
    },
    meta: function() {
      return [
        {
          name: "description",
          content:
            "DadosJusBr é uma plataforma que realiza a libertação continua de dados de remuneração do sistema de justiça brasileiro. Esta página mostra dados do estado" +
            this.selected,
          id: "desc",
        },
      ];
    },
  },
};
</script>

<style scoped>
.dropDownButton {
  width: 400px;
  height: 83px;
  margin-top: 50px;
  margin-bottom: 50px;
  border: 1px solid #ffffff;
}

@media only screen and (max-width: 650px) {
  .dropDownButton {
  width: 90%;
  }
}

</style>
