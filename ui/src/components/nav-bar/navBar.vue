<template>
  <b-row v-bind:style="colorNavBar" class="navBar d-flex align-items-center ">
    <b-col cols="7" v-bind:style="auxBarColor" class="p-xl-5 mt-xl-0 mt-4  ml-1 ml-xl-5"  >
      <router-link to="/">
        <img
          v-on:click="toHome()"
          src="../../assets/icon_dadosjusbr_colored.svg"
        />
      </router-link>
    </b-col>
    <b-col cols="2" class="ml-n2 ml-xl-5" >
      <router-link to="/dados">
        <b v-bind:style="textDecorationDados" class="textNav">
          Dados
        </b>
      </router-link>
    </b-col>
    <b-col cols="2" class="ml-1 ">
      <router-link to="/equipe">
        <b
          v-bind:style="textDecorationEquipe"
          class="textNav"
        >
          Equipe
        </b>
      </router-link>
    </b-col>
  </b-row>
</template>

<script>
export default {
  name: "navBar",
  data() {
    return {
      colorNavBar: {
        backgroundColor: "#3e5363",
        color: "white",
      },
      textDecorationDados: {
        textDecoration: "none",
        color: "white",
      },
      textDecorationEquipe: {
        textDecoration: "none",
        color: "white",
      },
      auxBarColor: {
       
      },
      imgPath: "../../assets/icon_dadosjusbr_colored.svg",
    };
  },
  methods: {
    toEquipe() {
      this.equipe = true;
    },
    toHome() {
      this.equipe = false;
    },
  },
  watch: {
    $route(to, from) {
      if (this.$route.path == "/equipe") {
        this.colorNavBar.backgroundColor = "#2fbb96";
        this.textDecorationEquipe.textDecoration = "underline";
        this.textDecorationDados.textDecoration = "none";
        this.colorNavBar.color = "#3E5363";
        this.textDecorationEquipe.color = "#3e5363";
        this.textDecorationDados.color = "#3e5363";
        this.auxBarColor.borderBottom = "1px solid #3e5363";
        this.imgPath = "../../assets/icon_dadosjusbr_equipe.svg";
      } else if (this.$route.path == "/") {
        this.colorNavBar.backgroundColor = "#3e5363";
        this.textDecorationEquipe.textDecoration = "none";
        this.textDecorationDados.textDecoration = "none";
        this.auxBarColor.borderBottom = "1px solid white";
        this.textDecorationEquipe.color = "white";
        this.textDecorationDados.color = "white";
        this.imgPath = "../../assets/icon_dadosjusbr_colored.svg";
      } else if (this.$route.path == "/dados") {
        this.colorNavBar.backgroundColor = "#3e5363";
        this.textDecorationDados.textDecoration = "underline";
        this.textDecorationEquipe.textDecoration = "none";
        this.textDecorationEquipe.color = "white";
        this.textDecorationDados.color = "white";
        this.auxBarColor.borderBottom = "1px solid white";
        this.imgPath = "../../assets/icon_dadosjusbr_colored.svg";
      } else {
        this.colorNavBar.backgroundColor = "#3e5363";
        this.imgPath = "../../assets/icon_dadosjusbr_colored.svg";
      }
    },
  },
};
</script>

<style scoped>

.textNav {
  color: white;
  font-size: 1.2rem;
}
.navBar {
  background-color: #3e5363;
}

</style>
