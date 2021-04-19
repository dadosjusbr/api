<template>
  <b-container fluid class="p-0">
  <b-row
    style="text-align: center; background-color: #3e5363; "
  >
  <b-col>
    <!-- watching the select property instead of relying on the onchange method due to this bug https://stackoverflow.com/a/53107234/5822594 -->
    <b-form-select size="lg" class="dropDownButton" v-model="selected" :options="options"></b-form-select>
  </b-col>
  </b-row>
  <b-row style="text-align: center;">
    <b-col cols="1" style="background-color: #3e5363;"></b-col>
    <b-col :key="agencies">
      <agency
        v-for="agency in agencies"
        :agency="agency"
        :key="agency.Name"
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
      selected: "PB",
      options: [
        { value: "Federal", text: 'Órgãos Federais' },
        {
          label: "Órgãos Estaduais",
          options: [
            { value: "BA", text: "Bahia" },
            { value: "CE", text: "Ceará" },
            { value: "ES", text: "Espírito Santo" },
            { value: "GO", text: "Goiás" },
            { value: "MG", text: "Minas Gerais" },
            { value: "PB", text: "Paraíba" },
            { value: "PR", text: "Paraná" },
            { value: "PE", text: "Pernambuco" },
            { value: "RN", text: "Rio Grande do Norte" },
            { value: "RS", text: "Rio Grande do Sul" },
            { value: "RJ", text: "Rio de Janeiro" },
            { value: "SP", text: "São Paulo" }
          ]
        }
      ],
      agencies: [],
    };
  },
  watch: {
    selected(value) {
      this.fetchData(value)
    },
  },
  methods: {
    async fetchData(state) {
      const { data } = await this.$http.get("/orgao/" + state);
      this.agencies = data.Agency;
    },
  },
  mounted() {
    this.fetchData(this.selected);
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
