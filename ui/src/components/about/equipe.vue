<template>
  <div>
    <div class="equipe">
      <b> TIME </b>
      <div class="time">
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
          :image="'https://avatars3.githubusercontent.com/u/33490097?s=460&v=4'"
          :github="'https://github.com/joeberth'"
        >
        </profile-card>
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
          :image="'https://avatars0.githubusercontent.com/u/5099338?s=460&v=4'"
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
    <div class="contribua">
      <div class="contribuaText">
        <b>SE INTERESSOU, QUER AJUDAR?</b>
        <br />
        <br />
        O sistema de justiça brasileiro é composto por uma grande quantidade de
        órgãos federais e estaduais distribuídos por todo o território nacional.
        <br />
        <br />
        Por isso, precisamos de contribuidores para criar crawlers e parsers
        para alimentar nossa base de dados.
      </div>
      <a
        style="text-decoration: none;"
        href="https://github.com/dadosjusbr/coletores/blob/master/TUTORIAL.md"
        ><div class="contribuaBotao">
          <b>
            SAIBA MAIS AQUI
          </b>
        </div>
      </a>
    </div>

    <div class="parceiros">
      <b> PARCEIROS </b>
      <br />
      <br />
      <div class="parceirosLogo">
        <img src="../../assets/logo_analytics.svg" class="logo" />
        <img src="../../assets/logo_shuttleworth.svg" class="logo" />
        <img src="../../assets/logo_APBMP.svg" class="logo" />
        <img src="../../assets/logo_mppb.svg" class="logo" />
      </div>
    </div>
  </div>
</template>

<script>
import profileCard from "@/components/about/profileCard.vue";

export default {
  name: "equipe",
  components: { profileCard },
  head: {
    title: function() {
      return {
        inner: "DadosJusBr",
        complement: "Equipe",
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
.equipe {
  font-size: 22px;
  color: #3e5363;
  background-color: #2fbb96;
  padding: 145px 100px 145px 100px;
}

.parceiros {
  background-color: #2fbb96;
  padding: 150px 100px 150px 100px;
  color: #3e5363;
  font-size: 22px;
}

.parceirosLogo {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
}

.contribuaBotao {
  width: 295px;
  height: 83px;
  border: solid 2px #ffffff;
  background-color: rgba(62, 83, 99, 0);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-left: 70px;
  margin-top: 20px;
}

.contribuaBotao:hover {
  background-color: #ffffff;
  color: #3e5363;
}

.contribuaText {
  width: 700px;
}

a {
  font-size: 22px;
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

.time {
  width: 100%;
  margin-top: 40px;
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  grid-template-rows: 350px 350px 350px;
  justify-items: center;
  align-items: center;
}

.logo {
  height: 245px;
  width: 295px;
}

@media only screen and (max-width: 500px) {
  .equipe {
    padding: 100px 20px 70px 20px;
  }
  .contribua {
    padding: 60px 20px 60px 20px;
    flex-direction: column;
    align-items: space-around;
    justify-content: center;
  }
  .contribuaBotao {
    width: 100%;
    margin: 0px 0px 0px 0px;
    text-decoration: none;
  }
  a {
    width: 100%;
    height: 83px;
  }
  .contribuaText {
    height: 370px;
    width: 100%;
  }

  .time {
    width: 100%;
    margin-top: 80px;
    display: grid;
    grid-template-columns: 1fr 1fr;
    grid-template-rows: 300px 300px 300px 300px 300px;
    justify-items: center;
    align-items: center;
  }
  .parceirosLogo {
    flex-wrap: wrap;
  }
  .logo {
    height: 135px;
    width: 162px;
  }
  .parceiros {
    padding: 75px 20px 75px 20px;
  }
  br {
    height: 15px;
    display: block !important;
  }
}
@media only screen and (max-width: 320px) {
  .logo {
    height: 120px;
    width: 140px;
  }
  .contribuaText {
    height: auto;
    width: 100%;
    margin-bottom: 30px;
  }
}
</style>
