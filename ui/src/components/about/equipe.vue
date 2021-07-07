<template>
  <b-container class="p-0" fluid>
    <b-row class="equipe">
      <b-col cols="1"> </b-col>
      <b-col class="mt-5">
        <b class="equipe">TIME</b> 
      </b-col>
      <b-col cols="1"> </b-col>
    </b-row>
    <b-row class="equipe d-flex flex-wrap align-items-center">
      <b-col class="p-xl-5 mt-3" cols="6" xl="4">
        <profile-card
          :name="'Joeberth Souza'"
          :role="'Desenvolvedor'"
          :image="'https://avatars3.githubusercontent.com/u/33490097?s=460&v=4'"
          :github="'https://github.com/joeberth'"
        >
        </profile-card>
      </b-col>
      <b-col class="p-xl-5 mt-3" cols="6" xl="4">
         <profile-card
          :name="'Daniel Fireman'"
          :role="'Professor'"
          :image="
            'https://avatars2.githubusercontent.com/u/8951363?s=400&u=acdedbbd00a15d2f913e3e1bde6b14f2c7451f90&v=4'
          "
          :github="'https://github.com/danielfireman'"
        >
        </profile-card>
      </b-col>
      <b-col class="p-xl-5 mt-3" cols="6" xl="4">
        <profile-card
          :name="'Lorena Pereira'"
          :role="'Desenvolvedora'"
          :image="
            'https://avatars1.githubusercontent.com/u/9660774?s=460&u=ac2f817f4ca07be4b34d9534759409a2f52be269&v=4'
          "
          :github="'https://github.com/Lorenaps'"
        >
        </profile-card>
      </b-col>
      <b-col class="p-xl-5 mt-3" cols="6" xl="4">
        <profile-card
          :name="'Marcos Barros'"
          :role="'Desenvolvedor'"
          :image="'https://avatars3.githubusercontent.com/u/34271729?s=460&v=4'"
          :github="'https://github.com/marcosbmf'"
        >
        </profile-card>
      </b-col>
      <b-col class="p-xl-5 mt-3" cols="6" xl="4">
        <profile-card
          :name="'Mariana Souto'"
          :role="'UI/UX'"
          :image="
            'https://avatars3.githubusercontent.com/u/15800324?s=400&u=11b49d52de39d773024ad1894cdbf33c1637532a&v=4'
          "
          :github="'https://github.com/soutoam'"
        >
        </profile-card>
        </b-col>
        <b-col class="p-xl-5 mt-3" cols="6" xl="4">
        <profile-card
          :name="'Vinicius Agostini'"
          :role="'Desenvolvedor'"
          :image="
            'https://avatars1.githubusercontent.com/u/12215291?s=460&u=313fcc28f44cb90215a50a7d512abf2b61768155&v=4'
          "
          :github="'https://github.com/viniagostini'"
        >
        </profile-card>
        </b-col>
      <b-col class="p-xl-5 mt-3" cols="6" xl="4">
        <profile-card
          :name="'David Eduardo'"
          :role="'Desenvolvedor'"
          :image="
            'https://avatars3.githubusercontent.com/u/25506334?s=460&u=46b571dfa3e0ed7c134f21a787954c2230ab9e68&v=4'
          "
          :github="'https://github.com/davidedup'"
        >
        </profile-card>
        </b-col>
      <b-col class="p-xl-5 mt-3" cols="6" xl="4">
        <profile-card
          :name="'Nazareno Andrade'"
          :role="'Professor'"
          :image="
            'https://avatars0.githubusercontent.com/u/1053850?s=460&u=e1cbfe4f7cba06cf5908f8967c2db24e142ff7c7&v=4  '
          "
          :github="'https://github.com/nazareno'"
        >
        </profile-card>
        </b-col>
      <b-col class="p-xl-5 mt-3" cols="6" xl="4">
        
        <profile-card
          :name="'João Brunet'"
          :role="'Professor'"
          :image="'https://avatars0.githubusercontent.com/u/5099338?s=460&v=4'"
          :github="'https://github.com/joaoarthurbm'"
        >
        </profile-card>
       </b-col>
    </b-row>
    <b-row class="contribua d-flex flex-wrap align-items-center">
       <b-col cols="1"></b-col>
      <b-col cols="11" xl="6" class="mt-5 ml-3 ml-xl-0 mb-5 contribuaText align-items-center justify-content-center">
        <b>SE INTERESSOU, QUER AJUDAR?</b>
        <br />
        <br />
        O sistema de justiça brasileiro é composto por uma grande quantidade de
        órgãos federais e estaduais distribuídos por todo o território nacional.
        <br />
        <br />
        Por isso, precisamos de contribuidores para criar crawlers e parsers
        para alimentar nossa base de dados.
      </b-col>
      <b-col cols="11" xl="5" class=" mt-5 mb-5  ml-3 ml-xl-0 d-flex align-items-center justify-content-center">
      <a
        style="text-decoration: none;"
        href="https://github.com/dadosjusbr/coletores/blob/master/TUTORIAL.md"
        ><div class="contribuaBotao">
          <b>
            SAIBA MAIS AQUI
          </b>
        </div>
      </a>
      </b-col>
       <b-col cols="1"></b-col>
    </b-row>
    <b-row class="parceiros"> 
      <b-col class="ml-3 mt-3">
        <b> REALIZAÇÂO </b>
        <br />
        <br />
      </b-col>
    </b-row>
    <b-row class="parceiros d-flex flex-wrap">
      <b-col cols="1" class="d-xl-none p-0"></b-col>
      <b-col xl="3" cols="12" class="   " ><img src="../../assets/logo_ufcg.svg" class="logo" /></b-col>
      <b-col xl="3" cols="12" class="   " ><img src="../../assets/logo_tb.svg" class="logo" /></b-col>
      <b-col cols="1" xl="0" class="d-xl-none p-0"></b-col>
    </b-row>

    <b-row class="parceiros"> 
      <b-col class="ml-3 mt-3">
        <b> PARCEIROS </b>
        <br />
        <br />
      </b-col>
    </b-row>
    <b-row class="parceiros d-flex flex-wrap">
      <b-col cols="1" class="d-xl-none p-0"></b-col>
      <b-col xl="3" cols="12" class="   " ><img src="../../assets/logo_analytics.svg" class="logo" /></b-col>
      <b-col xl="3" cols="12" class="p-0"><img src="../../assets/logo_shuttleworth.svg" class="logo" /></b-col>
      <b-col xl="3" cols="12" class="p-0"><img src="../../assets/logo_APBMP.svg" class="logo" /></b-col>
      <b-col xl="3" cols="12" class="p-0"><img src="../../assets/logo_mppb.svg" class="logo" /></b-col>
      <b-col xl="3" cols="12" class="p-0"><img src="../../assets/logo_inst_betty_jacob_laffer.svg" class="logo" /></b-col>
      <b-col cols="1" xl="0" class="d-xl-none p-0"></b-col>
    </b-row>
  </b-container>
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
  font-size: 1.3em;
  color: #3e5363;
  background-color: #2fbb96;
  
}

.parceiros {
  background-color: #2fbb96;
  color: #3e5363;
  font-size: 1.3em;
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
}

.contribuaBotao:hover {
  background-color: #ffffff;
  color: #3e5363;
}


a {
  font-size: 1.3em;
}

.contribua {  
  background-color: #3e5363;
  font-size: 1.3em;
  color: white;
}

.logo {
  height: 245px;
  width: 295px;
}

</style>
