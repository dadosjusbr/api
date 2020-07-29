<template>
  <div class="summary">
    <div class="wageInfoContainer">
      <div class="wageTotal">
        <div style="width: 92%">
          <p>Total remuneração: {{ this.agencySummary.TotalRemuneration }}</p>
        </div>
        <div style="width: 60px;">
          <md-icon id="tooltip-target-1">info</md-icon>
        </div>
      </div>
      <div class="othersTotals">
        <div class="othersTotalsInfo">
          <p>Maior Salário: {{ this.agencySummary.MaxWage }}</p>
          <br />
          <p>Total Salários: {{ this.agencySummary.TotalWage }}</p>
        </div>
        <div class="othersTotalsInfo">
          <p>Maior Benefício: {{ this.agencySummary.MaxPerk }}</p>
          <br />
          <p>Total benefícios: {{ this.agencySummary.TotalPerks }}</p>
        </div>
      </div>
    </div>
    <div class="whiteSpace"></div>
    <div class="employeeInfoContainer">
      <div class="toltalEmployees">
        <div style="width: 90%">
          <p>Total empregados: {{ this.agencySummary.TotalEmployees }}</p>
        </div>
        <div style="width: 60px;">
          <md-icon id="tooltip-target-2">info</md-icon>
        </div>
      </div>
      <div class="employeesClassification">
        <div class="employeeClass">
          <div
            style="background-color: #c9a0d0;"
            :class="[!members ? 'squareOpac' : '', 'square']"
            v-on:click="membersClick()"
          ></div>
          <p>Membros: {{ this.agencySummary.TotalMembers }}</p>
        </div>
        <div class="employeeClass">
          <div
            style="background-color: #513658;"
            :class="[!servants ? 'squareOpac' : '', 'square']"
            v-on:click="serventsClick()"
          ></div>
          <p>Servidores: {{ this.agencySummary.TotalServants }}</p>
        </div>
      </div>
    </div>
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
  </div>
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
.wageInfoContainer {
  height: 100%;
  display: table-cell;
  width: 66%;
  flex-direction: column;
}
.wageTotal {
  display: flex;
  flex-direction: row;
  font-size: 17px;
  height: 33%;
  width: 100%;
  line-height: 54.44px;
  text-align: center;
  order: 1;
  background-color: #dcdbdc;
}
.othersTotals {
  width: 100%;
  height: 77%;
  order: 3;
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: center;
  background-color: white;
}

.othersTotalsInfo {
  height: 70%;
  width: 50%;
  margin-left: 2%;
  border-left: solid 1.5px #dcdbdc;
  padding-left: 3%;
}

.employeeInfoContainer {
  height: 100%;
  display: table-cell;
  width: 33%;
  flex-direction: column;
}

.toltalEmployees {
  display: flex;
  flex-direction: row;
  font-size: 17px;
  order: 1;
  height: 33%;
  width: 100%;
  background-color: #dcdbdc;
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
  background-color: white;
}

.summary {
  height: 150px;
  display: table;
  justify-content: center;
  width: 100%;
  margin-top: 3%;
  margin-bottom: 40px;
}

.whiteSpace {
  width: 5px;
  height: 100%;
  order: 2;
  display: table-cell;
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
  margin-bottom: 15px;
  opacity: 0.2;
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

@media only screen and (max-width: 650px) {
  .summary {
    height: 470px;
    display: table;
    justify-content: center;
    width: 100%;
    margin-top: 3%;
    margin-bottom: 16px;
    display: flex;
    flex-direction: column;
  }

  .wageInfoContainer {
    height: 277px;
    display: table-cell;
    width: 100%;
    flex-direction: column;
    /* margin-bottom: 16px; */
  }
  .wageTotal {
    height: 59px;
    width: 100%;
    line-height: 54.44px;
    text-align: center;
    order: 1;
    background-color: #dcdbdc;
  }
  .employeeInfoContainer {
    height: 50%;
    display: table-cell;
    width: 100%;
    flex-direction: column;
  }

  .othersTotals {
    margin: 0 auto;
    flex-direction: column;
    height: 250px;
  }
  .othersTotalsInfo {
    padding-top: 2px;
    width: 250px;
    height: 80px;
    margin-bottom: 15px;
  }
  .whiteSpace {
    width: 100%;
    height: 16px;
    background-color: #f4f4f4;
    order: 2;
  }
  br {
    display: block;
    content: "";
    margin-top: 1px;
  }

  .toltalEmployees {
    height: 59px;
  }
  .employeeInfoContainer {
    height: 180px;
    order: 3;
  }

  .employeesClassification {
    height: 119px;
    padding-top: 16px;
  }
}
</style>
