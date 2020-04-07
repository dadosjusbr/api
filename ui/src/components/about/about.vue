<template>
  <div class="aboutContainer">
    <h1>SOBRE</h1>
    <p>
      Os agentes públicos do sistema de justiça brasileiro recebem
      outras verbas, além de seus salários, para exercerem seus cargos. Dentre
      elas encontramos auxílio moradia, despesas com saúde, auxílio transporte,
      gratificações, diárias, entre outros benefícios. O objetivo do DadosJusBr
      é apresentar de forma detalhada e unificada os dados de gastos com
      remuneração dos órgãos que constituem o poder judiciário brasileiro, assim
      facilitando o acesso e promovendo o controle social sobre esses gastos do
      poder judiciário.
    </p>

    <p>
      A
      <a
        href="http://www.planalto.gov.br/ccivil_03/_ato2011-2014/2011/lei/l12527.htm"
        ><b>Lei Federal 12.527/2011</b></a
      >, ou mais comumente conhecida como Lei de Acesso à Informação (LAI), diz que é “<i
        >obrigatória a divulgação em sítios oficiais da rede mundial de
        computadores (internet)</i
      >” d  os dados de gastos públicos. Porém, a lei da transparência pouco diz
      sobre a forma como esses dados devem ser disponibilizados. Por isso, nos
      sites dos órgãos são encontrados arquivos em diversos formatos (pdf, html,
      planilhas eletrônicas, json e etc), além disso nomenclaturas e formatação
      são muitas vezes diferentes para cada órgão. Devido a essas
      características realizar um controle social e financeiro sobre essa enorme
      quantidade de dados de gastos públicos é um tarefa difícil para o humano.
      O DadosJusBr utiliza a inteligência de dados para a ação cidadã,
      promovendo um acesso mais democrático e fácil aos dados de remuneração do
      poder judiciário brasileiro.
    </p>

    <p>
      No DadosJusBr podemos entender como cada juiz, promotor e desembargador
      são remunerados. Quais auxílios recebem? Quais os valores destes auxílios?
      Quanto além do salário um funcionário recebeu em determinado mês? Quanto
      um órgão gastou em determinado mês? Todas essas perguntas podem ser
      respondidas através do DadosJusBr.
    </p>

    <h3>SE INTERESSOU, QUER AJUDAR?</h3>

    <p>
      O Judiciário Brasileiro é composto por uma grande quantidade de órgãos
      distribuídos por todo o território nacional. Por isso, precisamos de
      contribuidores para criar crawlers e parsers para alimentar nossa base de
      dados. Saiba mais
      <a href="https://github.com/dadosjusbr/coletores/blob/master/TUTORIAL.md">
        <b> aqui </b></a
      >.
    </p>

    <h2>EQUIPE</h2>
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
        :image="'https://avatars3.githubusercontent.com/u/33490097?s=460&v=4'"
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
    </div>

    <div class="teamPhotosContainer">
      <profile-card
        :name="'João Brunet'"
        :role="'Professor'"
        :image="'https://avatars0.githubusercontent.com/u/5099338?s=460&v=4'"
        :github="'https://github.com/joaoarthurbm'"
      >
      </profile-card>
      <profile-card
        :name="'Marcos Barros'"
        :role="'Desenvolvedor'"
        :image="'https://avatars3.githubusercontent.com/u/34271729?s=460&v=4'"
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

    <!-- <div class="logoContainer">
      <img class="img1" src="../../assets/logo_ufcg.png" />
      <img class="img3" src="../../assets/logo_mppb.png" />
      <img class="img2" src="../../assets/logo_analytics.png" />
    </div> -->
  </div>
</template>

<script>
import profileCard from "@/components/about/profileCard.vue";

export default {
  name: "about",
  components: { profileCard },
};
</script>

<style scoped>
.teamPhotosContainer {
  padding-top: 5px;
  padding-right: 5px;
  padding-bottom: 5px;
  padding-left: 5px;
  margin-top: 20px;
  height: 300px;
  justify-content: center;
  display: flex;
  align-self: center;
  align-content: center;
}

.logoContainer {
  border: 3px black;
  text-align: center;
}

.aboutContainer {
  border: 1px solid #6a757a;
  margin-left: 200px;
  margin-right: 200px;
  margin-top: 10px;
  background-color: #6a757a;
  margin-bottom: 10px;
}

h1 {
  color: white;
  text-align: center;
  padding-top: 50px;
  font-size: 50px;
}

h2 {
  color: white;
  text-align: center;
  padding-top: 50px;
  font-size: 35px;
}

h3 {
  color: white;
  text-align: center;
  padding-top: 50px;
  font-size: 30px;
}

.img1 {
  width: 200px;
  height: 200px;
}

.img2 {
  width: 200px;
  height: 200px;
}

.img3 {
  width: 30%;
  height: 20%;
}

p {
  padding-left: 60px;
  padding-right: 60px;
  padding-top: 25px;
  text-align: justify;
  color: white;
  font-size: 20px;
}

a {
  color: white;
}

@media only screen and (max-width: 1276px) {
  .teamPhotosContainer {
    padding-top: 5px;
    padding-right: 5px;
    padding-bottom: 5px;
    padding-left: 5px;
    margin-top: 20px;
    height: 150px;
    justify-content: center;
    display: flex;
    align-self: center;
    align-content: center;
  }

  .aboutContainer {
    border: 1px solid #6a757a;
    margin-left: 100px;
    margin-right: 100px;
    margin-top: 10px;
    background-color: #6a757a;
    margin-bottom: 10px;
  }
}

@media only screen and (max-width: 1150px) {
  .teamPhotosContainer {
    padding-top: 5px;
    padding-right: 5px;
    padding-bottom: 5px;
    padding-left: 5px;
    margin-top: 20px;
    height: 150px;
    justify-content: center;
    display: flex;
    align-self: center;
    align-content: center;
  }

  .aboutContainer {
    border: 1px solid #6a757a;
    margin-left: 20px;
    margin-right: 20px;
    margin-top: 10px;
  }
}

@media only screen and (max-width: 950px) {
  .teamPhotosContainer {
    padding-top: 5px;
    padding-right: 5px;
    padding-bottom: 5px;
    padding-left: 5px;
    margin-top: 20px;
    height: 150px;
    justify-content: center;
    display: flex;
    align-self: center;
    align-content: center;
  }

  .aboutContainer {
    margin-left: 10px;
    margin-right: 10px;
  }
}

@media only screen and (max-width: 650px) {
  p {
    font-size: 15px;
  }
}

@media only screen and (max-width: 550px) {
  .aboutContainer {
    margin-left: 2px;
    margin-right: 2px;
  }
  p {
    font-size: 15px;
  }
}

@media only screen and (max-width: 450px) {
  .aboutContainer {
    margin-left: 2px;
    margin-right: 2px;
  }
  p {
    font-size: 15px;
  }
}

@media only screen and (max-width: 320px) {
  .aboutContainer {
    margin-left: 2px;
    margin-right: 2px;
  }
  p {
    font-size: 15px;
  }
}
</style>
