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
      O Poder Judiciário Brasileiro é composto por xxx órgãos, uma quantidade
      muito grande para um único time explorar todos eles. Por isso precisamos
      de contribuidores para criar crawlers e parsers para alimentar nossa base
      de dados. Saiba mais <a>aqui</a>.
    </p>

    <h2>EQUIPE</h2>
    <div class="teamPhotosContainer">
      <profile-card :name="'Maria Souza'" :role="'Dev'"> </profile-card>
      <profile-card :name="'Maria Souza'" :role="'Dev'"> </profile-card
      ><profile-card :name="'Maria Souza'" :role="'Dev'"> </profile-card
      ><profile-card :name="'Maria Souza'" :role="'Dev'"> </profile-card>
    </div>

    <div class="teamPhotosContainer">
      <profile-card :name="'Maria Souza'" :role="'Dev'"> </profile-card
      ><profile-card :name="'Maria Souza'" :role="'Dev'"> </profile-card
      ><profile-card :name="'Maria Souza'" :role="'Dev'"> </profile-card
      ><profile-card :name="'Maria Souza'" :role="'Dev'"> </profile-card>
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
  padding-top: 50px;
  padding-right: 50px;
  padding-bottom: 50px;
  padding-left: 50px;
  margin-top: 20px;
  margin-left: 100px;
  margin-right: 100px;
  height: 300;
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
  height: 1400px;
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
</style>
