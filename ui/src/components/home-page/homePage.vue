<template>
  <div class="sobre">
    <div class="infoCard infoCard1">
      <div class="infoCardAux">
        <b>VOCÊ JÁ TENTOU ACESSAR DADOS EM SITES DE ÓRGÃO PÚBLICOS? </b>
        <br />
        <br />
        A
        <a
          href="http://www.planalto.gov.br/ccivil_03/_ato2011-2014/2011/lei/l12527.htm"
          data-toggle="tooltip"
          title="É obrigação dos órgãos, independentemente de requerimentos, a divulgação em local de fácil acesso, dos seus registros de despesas. Também assegurado pela Lei Complementar 101/2000 (Lei de Responsabilidade Fiscal) a transparência é assegurada quando houver a liberação, em tempo real, de informações pormenorizadas sobre a execução orçamentária e financeira, em meios eletrônicos de acesso público"
          >Lei Federal 12.527/2011</a
        >, ou mais comumente conhecida como Lei de Acesso à Informação (LAI),
        diz que é “<i
          >obrigatória a divulgação em sítios oficiais da rede mundial de
          computadores (internet)</i
        >” dos dados de gastos públicos.
        <br />
        <br />
        Porém, a LAI pouco diz sobre a forma como esses dados devem ser
        disponibilizados. Por isso, nos sites dos órgãos são encontrados
        arquivos em diversos formatos (pdf, html, planilhas eletrônicas, json e
        etc), além disso nomenclaturas e formatação são muitas vezes diferentes
        para cada órgão.
        <br />
        <br />
        Devido a essas características realizar um controle social e financeiro
        sobre essa enorme quantidade de dados de gastos públicos é uma tarefa
        difícil para uma pessoa.
      </div>
      <img
        src="../../assets/Esquema_3.svg"
        style="height: 106px; width: 105px;"
      />
    </div>

    <div class="sloganContainer">
      <b>O DADOSJUSBR EXISTE PARA DENUNCIAR E LIBERTAR ESSES DADOS.</b>
    </div>

    <div class="infoCard" style="margin-top: 145px;">
      <div>
        <b> COMO NOS FAZEMOS ISSO? </b>
        <br />
        <br />
        <div style="infoCardAux">
          Os agentes públicos do sistema de justiça brasileiro recebem outras
          verbas, além de seus salários, para exercerem seus cargos. Dentre elas
          encontramos auxílio moradia, despesas com saúde, auxílio transporte,
          gratificações, diárias, entre outros benefícios.
        </div>
      </div>
      <img
        src="../../assets/Esquema_1.svg"
        style="height: 242px; width: 376px;"
      />
    </div>

    <div class="infoCard" style="padding-right: 10%;">
      <div style="infoCardAux">
        Inspirados em projetos como o
        <a href="https://serenata.ai/"> Serenata de amor </a> e
        <a href="https://brasil.io/"> Brasil.io</a>, o DadosJusBr surge com o
        objetivo de apresentar de forma detalhada, organizada e unificada os
        dados de gastos com remuneração dos órgãos que constituem o sistema de
        justiça brasileiro, assim facilitando o acesso e promovendo o controle
        social sobre esses gastos do poder judiciário, ministério público,
        defensoria pública e procuradorias.
      </div>
      <img
        src="../../assets/Esquema_2.svg"
        style="width: 112px; height: 151px;"
      />
    </div>

    <div class="infoCard">
      <div style="infoCardAux">
        O DadosJusBr utiliza a inteligência de dados para a ação cidadã,
        promovendo um acesso mais democrático e fácil aos dados de remuneração
        do sistema de justiça brasileiro.
        <br />
        <br />
        No DadosJusBr podemos entender como cada juiz, promotor e desembargador
        são remunerados. Quais auxílios recebem? Quais os valores destes
        auxílios? Quanto além do salário um funcionário recebeu em determinado
        mês? Quanto um órgão gastou em determinado mês? Todas essas perguntas
        podem ser respondidas através do DadosJusBr.
      </div>
      <img
        src="../../assets/Esquema_4.svg"
        style="heigh: 142px; width: 452px"
      />
    </div>
  </div>
</template>

<script>
export default {
  name: "homePage",
  head: {
    title: function() {
      return {
        inner: "DadosJusBr",
        complement: "Home",
      };
    },
    meta: function() {
      return [
        {
          name: "description",
          content:
            "DadosJusBr é uma plataforma que realiza a libertação continua de dados de remuneração do sistema de justiça brasileiro. Esta página mostra dados do orgão. No DadosJusBr podemos entender como cada juiz, promotor e desembargador são remunerados. Quais auxílios recebem? Quais os valores destes auxílios? Quanto além do salário um funcionário recebeu em determinado mês? Quanto um órgão gastou em determinado mês? Todas essas perguntas podem ser respondidas através do DadosJusBr.",
          id: "desc",
        },
      ];
    },
  },
};
</script>

<style scoped>
.sloganContainer {
  display: flex;
  justify-content: center;
  align-items: center;
  color: #ffffff;
  font-size: 22px;
  height: 250px;
  background-color: #b361c6;
}

.sobre {
  background-color: #3e5363;
  padding-top: 145px;
  height: 2200px;
}

a {
  color: white;
  font-size: 22px;
}

.infoCard {
  color: #ffffff;
  height: 151px;
  padding-left: 100px;
  padding-right: 40px;
  font-size: 22px;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  height: 250px;
  margin-bottom: 100px;
  background-color: #3e5363;
}

.infoCard1 {
  height: 500px;
  padding-right: 10%;
}

.infoCardAux {
  width: 700px;
}

@media only screen and (max-width: 500px) {
  .sobre {
    padding: 60px 20px 60px 20px;
    height: 3800px;
  }
  .infoCard {
    padding: 0px 0px 0px 0px;
    flex-direction: column;
    height: 750px;
  }
  .infoCardAux {
    width: 100%;
  }
  .sloganContainer {
    margin-right: -20px;
    margin-inline-start: -20px;
    text-align: center;
  }
}
</style>
