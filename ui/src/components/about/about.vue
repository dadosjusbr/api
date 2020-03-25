<template>
  <div class="aboutContainer">
    <h1>SOBRE</h1>
    <p>
      A <i>Lei da Transparência</i>, obriga os órgãos públicos a disponibilizar
      seus dados de gastos de forma online e pública. Um dos maiores gastos
      desses órgão é proveniente da folha de pagamento de funcionários. Essas
      folhas de pagamentos são normalmente disponibilizadas em arquivos de
      diversos formatos nos sites de cada um dos órgãos. Devido a essas
      características realizar um controle social sobre esses dados de gastos
      públicos é um tarefa difícil.
    </p>

    <p>
      O <b> DadosJus</b> surgiu com a proposta de realizar a unificação,
      libertação e normalização dos dados de remuneração dos órgãos que
      constituem o Poder Judiciário Brasileiro. O DadosJus é fruto da
      colaboração entre <b>Laboratório Analytics</b>,
      <b>Universidade Federal de Campina Grande</b> e o
      <b>Minitério Público da Paraíba</b>.
    </p>

    <h3>SE INTERESSOU, QUER AJUDAR?</h3>

    <p>
      O Poder Judiciário Brasileiro é composto por muitos órgãos distribuídos
      Por todo o território nacional, devido a essa quantidade muito grande
      torna-se inviável para um único time explorar todos eles. Por isso
      precisamos de contribuidores para criar crawlers e parsers para alimentar
      nossa base de dados. Saiba mais
      <a href="https://github.com/dadosjusbr/coletores/blob/master/TUTORIAL.md"
        >aqui</a
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
  components: { profileCard }
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
  height: 1450px;
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
    height: 1200px;
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
    height: 1200px;
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
    height: 1300px;
  }
}

@media only screen and (max-width: 730px) {
  .aboutContainer {
    height: 1390px;
  }
}

@media only screen and (max-width: 650px) {
  .aboutContainer {
    height: 1220px;
  }
  p {
    font-size: 15px;
  }
}

@media only screen and (max-width: 550px) {
  .aboutContainer {
    margin-left: 2px;
    margin-right: 2px;
    height: 1650px;
  }
  p {
    font-size: 15px;
  }
}

@media only screen and (max-width: 450px) {
  .aboutContainer {
    margin-left: 2px;
    margin-right: 2px;
    height: 1620px;
  }
  p {
    font-size: 15px;
  }
}

@media only screen and (max-width: 320px) {
  .aboutContainer {
    margin-left: 2px;
    margin-right: 2px;
    height: 1900px;
  }
  p {
    font-size: 15px;
  }
}
</style>
