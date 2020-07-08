<template>
  <div>
    <div class="infoContainer">
      <div class="about">
        <h3 style="font-weight: 600;">SOBRE O DADOSJUSBR</h3>
        <p>
          Os agentes públicos do sistema de justiça brasileiro recebem outras
          verbas, além de seus salários, para exercerem seus cargos. Dentre elas
          encontramos auxílio moradia, despesas com saúde, auxílio transporte,
          gratificações, diárias, entre outros benefícios. Inspirados em
          projetos como o
          <a href="https://serenata.ai/"> Serenata de amor </a> e
          <a href="https://brasil.io/"> Brasil.io</a>, o DadosJusBr surge com o
          objetivo de apresentar de forma detalhada, organizada e unificada os
          dados de gastos com remuneração dos órgãos que constituem o sistema de
          justiça brasileiro, assim facilitando o acesso e promovendo o controle
          social sobre esses gastos do poder judiciário, ministério público,
          defensoria pública e procuradorias.
        </p>

        <p>
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
          >” dos dados de gastos públicos. Porém, a LAI pouco diz sobre a forma
          como esses dados devem ser disponibilizados. Por isso, nos sites dos
          órgãos são encontrados arquivos em diversos formatos (pdf, html,
          planilhas eletrônicas, json e etc), além disso nomenclaturas e
          formatação são muitas vezes diferentes para cada órgão. Devido a essas
          características realizar um controle social e financeiro sobre essa
          enorme quantidade de dados de gastos públicos é uma tarefa difícil
          para uma pessoa. O DadosJusBr utiliza a inteligência de dados para a
          ação cidadã, promovendo um acesso mais democrático e fácil aos dados
          de remuneração do sistema de justiça brasileiro.
        </p>

        <p style="margin-bottom: 40px">
          No DadosJusBr podemos entender como cada juiz, promotor e
          desembargador são remunerados. Quais auxílios recebem? Quais os
          valores destes auxílios? Quanto além do salário um funcionário recebeu
          em determinado mês? Quanto um órgão gastou em determinado mês? Todas
          essas perguntas podem ser respondidas através do DadosJusBr.
        </p>
      </div>
      <div class="contribute">
        <h3 style="font-weight: 600;">SE INTERESSOU, QUER AJUDAR?</h3>

        <p style="margin-bottom: 40px">
          O sistema de justiça brasileiro é composto por uma grande quantidade
          de órgãos federais e estaduais distribuídos por todo o território
          nacional. Por isso, precisamos de contribuidores para criar crawlers e
          parsers para alimentar nossa base de dados. Saiba mais
          <a
            style="color: white;"
            href="https://github.com/dadosjusbr/coletores/blob/master/TUTORIAL.md"
          >
            <b> aqui </b></a
          >.
        </p>
      </div>
    </div>

    <div style="background-color: #7f3d8b;  padding-top: 67px;">
      <h3
        style="margin: 0px 0px 30px 67px; font-weight: 600; color: #ffffff; height: 50px;"
      >
        EQUIPE
      </h3>
      <div class="cardsContainer">
        <div class="teamPhotosContainer">
          <profile-card
            :name="'Daniel Fireman'"
            :role="'Professor'"
            :image="
              'https://avatars2.githubusercontent.com/u/8951363?s=400&u=acdedbbd00a15d2f913e3e1bde6b14f2c7451f90&v=4'
            "
            :github="'https://github.com/danielfireman'"
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
            :name="'David Eduardo'"
            :role="'Desenvolvedor'"
            :image="
              'https://avatars3.githubusercontent.com/u/25506334?s=460&u=46b571dfa3e0ed7c134f21a787954c2230ab9e68&v=4'
            "
            :github="'https://github.com/davidedup'"
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
        </div>
        <div class="teamPhotosContainer">
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
          <profile-card
            :name="'Mariana Souto'"
            :role="'Designer'"
            :image="
              'https://avatars3.githubusercontent.com/u/15800324?s=400&u=11b49d52de39d773024ad1894cdbf33c1637532a&v=4'
            "
            :github="'https://github.com/soutoam'"
          >
          </profile-card>
          <div class="auxFiller"></div>
        </div>
      </div>
    </div>

    <div
      style="background-color: #2fbb96; padding-top: 60px; padding-bottom: 30px"
    >
      <h3
        style="margin: 0px 0px 0px 67px; font-weight: 600; color: #ffffff; height: 50px;"
      >
        PARCEIROS
      </h3>
      <div class="partnersLogosContainer">
        <img src="../../assets/logo_analytics.svg" />
        <img src="../../assets/logo_shuttleworth.svg" />
        <img src="../../assets/logo_ufcg.svg" />
        <img src="../../assets/logo_mppb.svg" />
      </div>
    </div>
  </div>
</template>

<script>
import profileCard from "@/components/about/profileCard.vue";

export default {
  name: "about",
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
.partnersLogosContainer {
  display: flex;
  justify-content: space-evenly;
  padding: 10px 100px 50px 100px;
}

.infoContainer {
  display: flex;
  flex-direction: row;
}

.about {
  width: 60%;
  background-color: #ffffff;
  padding: 80px 70px 70px 70px;
  font-size: 15px;
  text-align: justify;
}
.auxFiller {
  height: 250px;
  width: 250px;
}

.contribute {
  width: 40%;
  background-color: #3e5363;
  color: #ffffff;
  padding: 240px 70px 70px 70px;
  text-align: justify;
  font-size: 15px;
}

.partnerCard {
  text-align: center;
  background-color: white;
  color: black;
}

.teamPhotosContainer {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-evenly;
}

a {
  color: black;
  font-weight: 500;
  font-size: 15px;
}

@media only screen and (max-width: 630px) {
  .infoContainer {
    flex-direction: column;
  }
  .about {
    width: 100%;
    min-height: 700px;
    padding: 20px 10px 10px 10px;
  }
  .contribute {
    width: 100%;
    padding: 20px 10px 10px 10px;
    text-align: justify;
    font-size: 15px;
  }
  .teamPhotosContainer {
    display: flex;
    flex-wrap: wrap;
    min-height: 600px;
    justify-content: space-evenly;
  }
  .partnersLogosContainer {
    flex-direction: column;
    padding: 10px 0px 0px 0px;
  }
  img {
    margin-bottom: 45px;
  }
  .auxFiller {
    height: 0px;
    width: 0px;
  }
  .cardsContainer {
    display: flex;
    flex-direction: column-reverse;
  }
}
</style>
