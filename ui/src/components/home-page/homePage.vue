<template>
  <div>
    <div class="homePageContainer">
      <div class="sobre">
        <div class="infoCard">
          <div>
            <b> SOBRE O DADOSJUS </b>
            <br />
            <br />
            <div style="width: 750px">
              Os agentes públicos do sistema de justiça brasileiro recebem
              outras verbas, além de seus salários, para exercerem seus cargos.
              Dentre elas encontramos auxílio moradia, despesas com saúde,
              auxílio transporte, gratificações, diárias, entre outros
              benefícios.
            </div>
          </div>
          <img src="../../assets/icones1.jpeg" />
        </div>

        <div class="infoCard">
          <div style="width: 750px">
            Inspirados em projetos como o
            <a href="https://serenata.ai/"> Serenata de amor </a> e
            <a href="https://brasil.io/"> Brasil.io</a>, o DadosJusBr surge com
            o objetivo de apresentar de forma detalhada, organizada e unificada
            os dados de gastos com remuneração dos órgãos que constituem o
            sistema de justiça brasileiro, assim facilitando o acesso e
            promovendo o controle social sobre esses gastos do poder judiciário,
            ministério público, defensoria pública e procuradorias.
          </div>
          <img src="../../assets/icones2.jpeg" />
        </div>

        <div class="infoCard" style="height: 450px">
          <div style="width:750px ">
            A
            <a
              href="http://www.planalto.gov.br/ccivil_03/_ato2011-2014/2011/lei/l12527.htm"
              data-toggle="tooltip"
              title="É obrigação dos órgãos, independentemente de requerimentos, a divulgação em local de fácil acesso, dos seus registros de despesas. Também assegurado pela Lei Complementar 101/2000 (Lei de Responsabilidade Fiscal) a transparência é assegurada quando houver a liberação, em tempo real, de informações pormenorizadas sobre a execução orçamentária e financeira, em meios eletrônicos de acesso público"
              >Lei Federal 12.527/2011</a
            >, ou mais comumente conhecida como Lei de Acesso à Informação
            (LAI), diz que é “<i
              >obrigatória a divulgação em sítios oficiais da rede mundial de
              computadores (internet)</i
            >” dos dados de gastos públicos.
            <br />
            <br />
            Porém, a LAI pouco diz sobre a forma como esses dados devem ser
            disponibilizados. Por isso, nos sites dos órgãos são encontrados
            arquivos em diversos formatos (pdf, html, planilhas eletrônicas,
            json e etc), além disso nomenclaturas e formatação são muitas vezes
            diferentes para cada órgão.
            <br />
            <br />
            Devido a essas características realizar um controle social e
            financeiro sobre essa enorme quantidade de dados de gastos públicos
            é uma tarefa difícil para uma pessoa.
          </div>
          <img src="../../assets/icones3.jpeg" />
        </div>

        <div class="infoCard">
          <div style="width: 750px">
            O DadosJusBr utiliza a inteligência de dados para a ação cidadã,
            promovendo um acesso mais democrático e fácil aos dados de
            remuneração do sistema de justiça brasileiro.
            <br />
            <br />
            No DadosJusBr podemos entender como cada juiz, promotor e
            desembargador são remunerados. Quais auxílios recebem? Quais os
            valores destes auxílios? Quanto além do salário um funcionário
            recebeu em determinado mês? Quanto um órgão gastou em determinado
            mês? Todas essas perguntas podem ser respondidas através do
            DadosJusBr.
          </div>
        </div>
      </div>
    </div>
    <div class="contribua">
      <div style="width:750px;">
        <b>SE INTERESSOU, QUER AJUDAR?</b>
        <br />
        <br />
        O sistema de justiça brasileiro é composto por uma grande quantidade de
        órgãos federais e estaduais distribuídos por todo o território nacional.
        Por isso, precisamos de contribuidores para criar crawlers e parsers
        para alimentar nossa base de dados.
      </div>
      <div class="contribuaBotao">
        <a
          href="https://github.com/dadosjusbr/coletores/blob/master/TUTORIAL.md"
        >
          <b> SAIBA MAIS AQUI </b></a
        >
      </div>
    </div>
    <div class="parceiros">
      <b> PARCEIROS </b>
      <div
        style="display: flex; flex-direction: row; justify-content: space-between;"
      >
        <img src="../../assets/logo_analytics.svg" />
        <img src="../../assets/logo_shuttleworth.svg" />
        <img src="../../assets/logo_ufcg.svg" />
        <img src="../../assets/logo_mppb.svg" />
      </div>
    </div>
    <div class="equipe">
      <b> EQUIPE </b>
      <div>
        <div class="equipeLinha">
          <profile-card
            :name="'Marcos Barros'"
            :role="'Desenvolvedor'"
            :image="
              'https://avatars3.githubusercontent.com/u/34271729?s=460&v=4'
            "
            :github="'https://github.com/marcosbmf'"
          >
          </profile-card>
          <profile-card
            :name="'Nazareno Andrade'"
            :role="'Professor'"
            :image="
              'https://avatars0.githubusercontent.com/u/1053850?s=460&u=e1cbfe4f7cba06cf5908f8967c2db24e142ff7c7&v=4  '
            "
            :github="'https://github.com/nazareno'"
          >
          </profile-card>
          <profile-card
            :name="'Vinicius Agostini'"
            :role="'Desenvolvedor'"
            :image="
              'https://avatars1.githubusercontent.com/u/12215291?s=460&u=313fcc28f44cb90215a50a7d512abf2b61768155&v=4'
            "
            :github="'https://github.com/viniagostini'"
          >
          </profile-card>
        </div>
        <div class="equipeLinha">
          <profile-card
            :name="'David Eduardo'"
            :role="'Desenvolvedor'"
            :image="
              'https://avatars3.githubusercontent.com/u/25506334?s=460&u=46b571dfa3e0ed7c134f21a787954c2230ab9e68&v=4'
            "
            :github="'https://github.com/davidedup'"
          >
          </profile-card>
          <profile-card
            :name="'Lorena Pereira'"
            :role="'Desenvolvedora'"
            :image="
              'https://avatars1.githubusercontent.com/u/9660774?s=460&u=ac2f817f4ca07be4b34d9534759409a2f52be269&v=4'
            "
            :github="'https://github.com/Lorenaps'"
          >
          </profile-card>
          <profile-card
            :name="'Joeberth Souza'"
            :role="'Desenvolvedor'"
            :image="
              'https://avatars3.githubusercontent.com/u/33490097?s=460&v=4'
            "
            :github="'https://github.com/joeberth'"
          >
          </profile-card>
        </div>
        <div class="equipeLinha">
          <profile-card
            :name="'Mariana Souto'"
            :role="'UI/UX'"
            :image="
              'https://avatars3.githubusercontent.com/u/15800324?s=400&u=11b49d52de39d773024ad1894cdbf33c1637532a&v=4'
            "
            :github="'https://github.com/soutoam'"
          >
          </profile-card>
          <profile-card
            :name="'João Brunet'"
            :role="'Professor'"
            :image="
              'https://avatars0.githubusercontent.com/u/5099338?s=460&v=4'
            "
            :github="'https://github.com/joaoarthurbm'"
          >
          </profile-card>
          <profile-card
            :name="'Daniel Fireman'"
            :role="'Professor'"
            :image="
              'https://avatars2.githubusercontent.com/u/8951363?s=400&u=acdedbbd00a15d2f913e3e1bde6b14f2c7451f90&v=4'
            "
            :github="'https://github.com/danielfireman'"
          >
          </profile-card>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import profileCard from "@/components/about/profileCard.vue";

export default {
  name: "homePage",
  components: { profileCard },
  head: {
    title: function() {
      return {
        inner: "DadosJusBr",
        complement: "Sobre",
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
.equipeLinha {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  margin-top: 120px ;
}

.equipe {
  font-size: 22px;
  color: white;
  background-color: #2fbb96;
  padding: 145px 100px 145px 100px;
}

.parceiros {
  background-color: #7f3d8b;
  padding: 150px 100px 150px 100px;
  color: white;
  font-size: 22px;
}

img {
  height: 230px;
  width: 250px;
  padding: 20px 20px 20px 20px;
}

.contribuaBotao {
  width: 295px;
  height: 83px;
  border: solid 2px #ffffff;
  background-color: rgba(62, 83, 99, 0);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-left: 70px;
  margin-top: 20px;
}

.contribua {
  padding-bottom: 150px;
  padding-top: 150px;
  background-color: #3e5363;
  font-size: 22px;
  color: white;
  padding-left: 100px;
  display: flex;
  flex-direction: row;
  align-items: center;
}

.sobre {
  background-color: #7f3d8b;
  padding-top: 145px;
  padding-bottom: 145px;
}

a {
  color: white;
  font-size: 22px;
}

.infoCard {
  color: #ffffff;
  height: 151px;
  margin-left: 100px;
  font-size: 22px;
  width: 90%;
  display: flex;
  flex-direction: row;
  height: 250px;
  margin-bottom: 100px;
}
</style>
