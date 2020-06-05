<template>
  <div class="summary">
    <div class="wageInfoContainer">
      <div class="wageTotal">
        <p>Total remuneração: {{this.agencySummary.TotalRemuneration}}</p>
      </div>
      <div class="othersTotals">
        <div class="othersTotalsInfo">
          <p>Salário máximo: {{ this.agencySummary.MaxWage }}</p>
          <br />
          <p>Soma total salário: {{ this.agencySummary.TotalWage }}</p>
        </div>
        <div class="othersTotalsInfo">
          <p>Benefício máximo: {{this.agencySummary.MaxPerk}}</p>
          <br />
          <p>
            Soma total benefícios: {{ this.agencySummary.TotalPerks }}
          </p>
        </div>
      </div>
    </div>
    <div class="whiteSpace"></div>
    <div class="employeeInfoContainer">
      <div class="toltalEmployees">
        <p>Total empregados: {{ this.agencySummary.TotalEmployees }}</p>
      </div>
      <div class="employeesClassification">
        <div class="employeeClass">
          <div style="background-color: #c9e4ca;" class="square"></div>
          <p>Membros: {{ this.agencySummary.TotalMembers }}</p>
        </div>
        <div class="employeeClass">
          <div style="background-color: #87bba2;" class="square"></div>
          <p>Servidores: {{ this.agencySummary.TotalServants }}</p>
        </div>
        <div
          class="employeeClass                                                                                                                                                           "
        >
          <div style="background-color: #364958;" class="square"></div>
          <p>Inativos: {{ this.agencySummary.TotalInactives }}</p>
        </div>
      </div>
    </div>
  </div>

  <!-- <info-card
      v-for="(value, name, id) in agencySummary"
      :key="id"
      :info="{ value, name }"
    /> -->
  <!-- </div> -->
</template>

<script>
import infoCard from "@/components/agency/infoCard.vue";

export default {
  name: "agencySummary",
  components: {
    infoCard,
  },
  props: {
    agencySummary: {
      type: Object,
      default: null,
    },
  },
};
</script>

<style lang="css" scoped>
.wageInfoContainer {
  height: 100%;
  display: table-cell;
  width: 66%;
  flex-direction: column;
}
.wageTotal {
  font-size: 17px;
  height: 33%;
  width: 100%;
  line-height: 54.44px;
  text-align: center;
  order: 1;
  background-color: #DCDBDC;
}
.othersTotals {
  width: 100%;
  height: 77%;
  order: 3;
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: center;
}

.othersTotalsInfo {
  height: 70%;
  width: 50%;
  margin-left: 2%;
  border-left: solid 1.5px #DCDBDC;
  padding-left: 3%;
}

.employeeInfoContainer {
  height: 100%;
  display: table-cell;
  width: 33%;
  flex-direction: column;
}

.toltalEmployees {
  font-size: 17px;
  order: 1;
  height: 33%;
  width: 100%;
  background-color:  #DCDBDC;
  line-height: 54.44px;
  text-align: center;
}

.employeesClassification {
  order: 2;
  height: 77%;
  width: 100%;
  flex-direction: row;
  display: flex;
  justify-content: center;
  align-items: center;
}

.summary {
  height: 150px;
  display: table;
  justify-content: center;
  width: 100%;
  margin-top: 3%;
}
.whiteSpace {
  width: 5px;
  height: 100%;
  order: 2;
  display: table-cell;
}

.square {
  height: 55px;
  width: 55px;
  margin-bottom: 15px;
}

.employeeClass {
  font-size: 13px;
  width: 33.33%;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}
</style>
