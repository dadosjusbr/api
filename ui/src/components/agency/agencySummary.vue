<template>
  <b-container fluid>
    <b-row class=" mt-xl-5 mb-5">
      <b-col cols="7">
        <b-row class="wageTotal align-items-center justify-content-center"> 
          <b-col cols="2">
            <img
              style="height: 36px; width:36px; margin-top:-30%; margin-left: -50%"
              src="../../assets/icon-remuneracao.svg"
            />
          </b-col>
          <b-col style="text-align: center;">
              <p>Total remuneração: {{ this.agencySummary.TotalRemuneration }}</p>
          </b-col>
          <b-col cols="2">
              <img
                id="tooltip-target-1"
                style="height: 36px; width:36px; margin-top: 0%; margin-left: 50%"
                src="../../assets/icon-info.svg"
              />
          </b-col>
        </b-row> 
        <b-row class="secondBox align-items-center justify-content-center pt-xl-3">
          <b-col class="p-xl-0" cols="2">
            <img
              src="../../assets/icon-salario.svg"
            />
          </b-col> 
          <b-col cols="4" >
              <p>Maior Salário: {{ this.agencySummary.MaxWage }}</p>
              <br />
              <p>Total Salários: {{ this.agencySummary.TotalWage }}</p>
          </b-col>
          <b-col cols="2">
              <img
                style="height: 36px; width:36px"
                src="../../assets/icon-beneficio.svg"
              />
          </b-col> 
          <b-col class="p-xl-0" cols="4">
              <p>Maior Benefício: {{ this.agencySummary.MaxPerk }}</p>
              <br />
              <p>Total benefícios: {{ this.agencySummary.TotalPerks }}</p>
          </b-col> 
        </b-row> 
      </b-col>
      <div class="whiteSpace"></div>
      <b-col >
        <b-row class="wageTotal">       
          <b-col cols="2" class="pl-xl-0 pt-xl-0">
              <img
                style="height: 36px; width:36px; margin-top:-30%"
                src="../../assets/icon-empregados.svg"
              />
          </b-col>
          <b-col cols="8" class="ml-xl-n4 pl-xl-0">              
                <p>Total empregados: {{ this.agencySummary.TotalEmployees }}</p>
              
          </b-col>
          <b-col cols="2">
                <img
                  id="tooltip-target-2"
                  style="height: 36px; width:36px; margin-top: 0%; margin-left: 70%"
                  src="../../assets/icon-info.svg"
                />
          </b-col>
        </b-row>
        
        <b-row class="secondBox  pt-xl-4">
          <b-col  >
            <b-row class="d-flex align-items-center justify-content-center">
              <div
                style="background-color: #c9a0d0;"
                :class="[!members ? 'squareOpac' : '', 'square']"
                v-on:click="membersClick()"
              >
              <img style="height: 100%; width: 100%" src="../../assets/icon_membros.svg"/></div>
            </b-row>
            <b-row class="d-flex align-items-center justify-content-center">
              <p>Membros: {{ this.agencySummary.TotalMembers }}</p>
            </b-row>
          </b-col>
          <b-col class="employeeClass">
            <b-row class="d-flex align-items-center justify-content-center">
              <div
                style="background-color: #513658;"
                :class="[!servants ? 'squareOpac' : '', 'square']"
                v-on:click="serventsClick()"
              >
              <img style="height: 100%; width: 100%" src="../../assets/icon_servidores.svg"/>
              </div>
            </b-row>
            <b-row class="d-flex align-items-center justify-content-center">
            <p>Servidores: {{ this.agencySummary.TotalServants }}</p>
            </b-row>
          </b-col>
        </b-row>
        <b-tooltip target="tooltip-target-1" triggers="hover">
          - Salário: valor recebido de acordo com a prestação de serviços, em
          decorrência do contrato de trabalho.
          <br />
          - Remuneração: é a soma do salário mais outras vantagens (indenizações e
          benefícios).
          <br />
          - Benefício: valores eventuais, por exemplo, auxílios alimentação, saúde,
          escolar...
        </b-tooltip>
        <b-tooltip target="tooltip-target-2" triggers="hover">
          - Membro: é o integrante da carreira 'principal' do órgão do sistema de
          justiça. Por exemplo, juízes, desembargadores, ministros, defensores,
          procuradores públicos, promotores de justiça, procuradores de justiça, etc
          <br />
          - Servidor: é todo integrante da carreira 'auxiliar', ou seja, são os
          analistas, técnicos, oficiais de justiça, etc.
        </b-tooltip>
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
export default {
  name: "agencySummary",
  props: {
    agencySummary: {
      type: Object,
      default: null,
    },
  },
  data() {
    return {
      members: true,
      servants: true,
    };
  },
  methods: {
    membersClick() {
      if (this.members) this.$emit("disable-members");
      else this.$emit("enable-members");
      this.members = !this.members;
    },
    serventsClick() {
      if (this.servants) this.$emit("disable-servants");
      else this.$emit("enable-servants");
      this.servants = !this.servants;
    },
  },
};
</script>

<style lang="css">

.wageTotal {
  color: white;
  font-size: 1.5rem;
  text-align: center;
  background-color: #3e5363;
  border-bottom: 1px solid white;
  line-height: 60px;
}

.secondBox {
  color: white;
  font-size: 1.0rem;
  text-align: center;
  background-color: #3e5363;
  border-bottom: 1px solid white;
}

.employeesClassification {
  justify-content: center;
  align-items: center;
}

.whiteSpace {
  width: 15px;
  height: 100%;
}

.square {
  cursor: pointer;
  height: 55px;
  width: 55px;
  margin-bottom: 15px;
}
.squareOpac {
  cursor: pointer;
  height: 45px;
  width: 45px;
  opacity: 0.2;
}

</style>
