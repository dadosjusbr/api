<template>
  <b-container fluid>
    <b-row class=" mt-xl-5 mb-1 d-flex flex-wrap">
      <b-col cols="12" xl="7" class="mb-4">
        <b-row class="wageTotal align-items-center justify-content-center"> 
          <b-col cols="2">
            <img
              style="height: 36px; width:36px; margin-top:-30%; margin-left: -50%"
              src="../../assets/icon-remuneracao.svg"
            />
          </b-col>
          <b-col style="text-align: center;" class="mt-2 mt-xl-0">
              <p>{{this.agencySummary.TotalMembers}} Membros</p>
          </b-col>
          <b-col cols="2">
              <img
                id="tooltip-target-1"
                style="height: 36px; width:36px; margin-top: 0%; margin-left: 50%"
                src="../../assets/icon-info.svg"
              />
          </b-col>
        </b-row> 
        <b-row class="secondBox align-items-center justify-content-center">
          <b-col class="d-none d-xl-block" cols="2">
            <img
              src="../../assets/icon-salario.svg"
            />
          </b-col> 
          <b-col class="mt-1" cols="12" xl="4" >
              <p>Maior Salário: {{ this.agencySummary.MaxWage }}</p>
              <p>Total Salários: {{ this.agencySummary.TotalWage }}</p>
          </b-col>
          <b-col cols="2" class="d-none d-xl-block">
              <img
                style="height: 36px; width:36px"
                src="../../assets/icon-beneficio.svg"
              />
          </b-col> 
          <b-col class="mt-1 p-xl-0" cols="12" xl="4">
              <p>Maior Benefício: {{ this.agencySummary.MaxPerk }}</p>
              <p>Total benefícios: {{ this.agencySummary.TotalPerks }}</p>
          </b-col> 
        </b-row> 
      </b-col>
      <!-- <div class="whiteSpace"></div> -->
      <!-- <b-col cols="12" xl="4" class="mb-2" >
        <b-row class="wageTotal">       
          <b-col cols="2" class="pl-xl-0 ">
              <img
                style="height: 36px; width:36px; "
                src="../../assets/icon-empregados.svg"
              />
          </b-col>
          <b-col cols="8" class="mt-2 mt-xl-0 d-flex align-items-center justify-content-center">              
                <p>Total empregados: {{ this.agencySummary.TotalMembers }}</p>              
          </b-col>
          <b-col cols="2">
                <img
                  id="tooltip-target-2"
                  style="height: 36px; width:36px; margin-top: 50%; margin-left: 20%"
                  src="../../assets/icon-info.svg"
                />
          </b-col>
        </b-row>
        
        <b-row class="secondBox  pt-xl-4">
          <b-col  class="mt-4 mt-xl-0">
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
          <b-col class="mt-4 mt-xl-0">
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
        </b-col> -->
        <b-tooltip target="tooltip-target-1" triggers="hover">
          - Salário: valor recebido de acordo com a prestação de serviços, em
          decorrência do contrato de trabalho.
          <br />
          - Remuneração: é a soma do salário mais outras vantagens (indenizações e
          benefícios).
          <br />
          - Benefício: valores eventuais, por exemplo, auxílios alimentação, saúde,
          escolar...
          - Membro: é o integrante da carreira 'principal' do órgão do sistema de
          justiça. Por exemplo, juízes, desembargadores, ministros, defensores,
          procuradores públicos, promotores de justiça, procuradores de justiça, etc
        </b-tooltip>
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
};
</script>

<style lang="css">

.wageTotal {
  color: white;
  font-size: 1.4em;
  text-align: center;
  background-color: #3e5363;
  border-bottom: 1px solid white;
  height: 4em !important;
  justify-content:center;
}

.secondBox {
  color: white;
  font-size: 1.0rem;
  text-align: center;
  background-color: #3e5363;
  border-bottom: 1px solid white;
  height: 10em !important;
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

@media only screen and (max-width: 700px) {
  .wageTotal {
  font-size: 1.2em;
}

}

</style>
