<template>
  <b-container class=" p-0 footerContainer " fluid v-bind:style="colorFooter" >
    <b-row class="mt-0" >
        <b-col v-bind:style="colorNavBar"  cols="1"></b-col>
        <b-col style="border-top: 1px solid white;"></b-col>
        <b-col v-bind:style="colorNavBar"  cols="1"></b-col>
  </b-row>
    <b-row class="p-xl-5">
      <b-col cols="1" class="d-flex"></b-col>
      <b-col cols="3" class="d-none d-xl-block"
       >
      <img src="../../assets/icon_dadosjusbr_white.svg" alt="" />
      </b-col>
      <b-col xl="7" cols="10" class="mt-5 align-content-center" >
        Alguma dica? Tem algum feedback geral? Se você tiver uma ideia que você
        gostaria de ver no DadosJusBr, envie-nos um e-mail
        (dadosjusbr@gmail.com) ou visite nosso
        <a
          style="color: white; font-size: 1.2em;"
          href="https://github.com/dadosjusbr"
        >
          github
        </a>
      </b-col>
      <b-col cols="1"></b-col>
    </b-row>
  </b-container>
</template>

<script>
export default {
  name: "pageFooter",
  data() {
    return {
      colorFooter: {
        backgroundColor: "#3e5363",
      },
    };
  },
  watch: {
    $route(to, from) {
      if (this.$route.path == "/equipe") {
        this.colorFooter.backgroundColor = "#2fbb96";
      } else if (this.$route.path == "/") {
        this.colorFooter.backgroundColor = "#3e5363";
      } else if (this.$route.path == "/dados") {
        this.colorFooter.backgroundColor = "#3e5363";
      } else {
         this.colorFooter.backgroundColor = "white";
      }
    },
  },
};
</script>

<style scoped>
.footerContainer {
  background-color: #3e5363;
  color: white;
  font-size: 1.3em;
}

</style>
