<template>
  <b-container class="sobre p-0"  fluid>

    <b-row class="infoCard p-xl-5 d-flex flex-wrap align-items-center">
      <b-col class="ml-4 mt-5" xl="7">
        <b>VOCÊ JÁ TENTOU ACESSAR DADOS EM SITES DE ÓRGÃOS PÚBLICOS? </b>
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
        etc). Além disso nomenclaturas e formatação são muitas vezes diferentes
        para cada órgão.
        <br />
        <br />
        Devido a essas características, realizar um controle social e financeiro
        sobre essa enorme quantidade de dados de gastos públicos é uma tarefa
        difícil para uma pessoa.
      </b-col>
      <b-col
        class="d-flex justify-content-xl-end justify-content-center mt-xl-0 mt-5 mb-xl-0 mb-4 "
        xl="4"
      >
        <img
          src="../../assets/Esquema_3.svg"
          style="height: 106px; width: 105px;"
        />
      </b-col>
    </b-row>

    <b-row class="sloganContainer d-flex align-items-center justify-content-center   mt-5 mb-xl-0 mb-4 "
        xl="6">
        <b-col cols="1"></b-col>
      <b-col>O DADOSJUSBR EXISTE PARA DENUNCIAR E LIBERTAR ESSES DADOS.</b-col>
      <b-col cols="1"></b-col>
    </b-row>

    <b-row class="infoCard p-xl-5 d-flex flex-wrap align-items-center">
        <b-col class="ml-4 mt-5" xl="7">
        <b> COMO NÓS FAZEMOS ISSO? </b>
        <br />
        <br />    
        Os agentes públicos do sistema de justiça brasileiro recebem outras
        verbas, além de seus salários, para exercerem seus cargos. Dentre elas
        encontramos auxílio moradia, despesas com saúde, auxílio transporte,
        gratificações, diárias, entre outros benefícios.
        </b-col>
        <b-col class="d-flex  justify-content-center mt-xl-0 mt-5 mb-xl-0 mb-4 "
        xl="4">
      <img
        src="../../assets/Esquema_1.svg"
        style="height: 242px; width: 376px;"
      />
      </b-col>
    </b-row>

    <b-row class="infoCard p-xl-5 d-flex flex-wrap align-items-center" >
      <b-col class="ml-4 mt-5" xl="7">
        Inspirado em projetos como o
        <a href="https://serenata.ai/"> Serenata de amor </a> e
        <a href="https://brasil.io/"> Brasil.io</a>, o DadosJusBr surge com o
        objetivo de apresentar de forma detalhada, organizada e unificada os
        dados de gastos com remuneração dos órgãos que constituem o sistema de
        justiça brasileiro, assim facilitando o acesso e promovendo o controle
        social sobre esses gastos do poder Judiciário, Ministério Público, 
        Defensoria Pública e Procuradorias.
      </b-col>
      <b-col class="d-flex justify-content-center mt-xl-0 mt-5 mb-xl-0 mb-4 "
        xl="4">
      <img
        src="../../assets/Esquema_2.svg"
        style="width: 112px; height: 151px;"
      />
      </b-col>
    </b-row>
    
    <b-row class="infoCard p-xl-5 d-flex flex-wrap align-items-center">
      <b-col  class="ml-4 mt-5" xl="7">
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
      </b-col>
      <b-col class="d-flex justify-content-center mt-xl-0 mt-5 mb-xl-0 mb-4 "
        xl="4">
      <img
        src="../../assets/Esquema_4.svg"
        style="heigh: 142px; width: 452px"
      />
      </b-col>
    </b-row>
  </b-container>
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
  height: 5em !important;
  color: #ffffff;
  font-size: 1.5rem;
  background-color: #b361c6;
  text-align: center;
}

.sobre {
  background-color: #3e5363;
}

a {
  color: white;
  font-size: 22px;
}

.infoCard {
  color: #ffffff;
  font-size: 22px;
  background-color: #3e5363;
}

</style>
