<template>
  <div v-bind:style="colorNavBar" class="navBar">
    <div v-bind:style="auxBarColor" class="auxBar">
      <router-link to="/">
        <img
          v-on:click="toHome()"
          src="../../assets/icon_dadosjusbr_colored.svg"
        />
      </router-link>
      <div>
        <router-link to="/dados">
          <b v-bind:style="textDecorationDados" class="dados">
            Dados
          </b>
        </router-link>
        <router-link to="/equipe">
          <b
            v-bind:style="textDecorationEquipe"
            style="color: white; font-size: 18px;"
          >
            Equipe
          </b>
        </router-link>
      </div>
    </div>
  </div>
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
        borderBottom: "1px solid white",
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
.dados {
  color: white;
  font-size: 18px;
  margin-right: 90px;
}

.navBar {
  background-color: #3e5363;
  height: 185px;
  padding: 50px 100px 0px 100px;
}

.auxBar {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  border-bottom: 1px solid white;
  padding-right: 100px;
  padding-bottom: 35px;
  align-items: center;
}

@media only screen and (max-width: 500px) {
  .navBar {
    padding: 20px 20px 40px 20px;
  }
  .auxBar {
    padding: 40px 20px 40px 20px;
  }
  .dados {
    margin-right: 20px;
  }
}
@media only screen and (max-width: 330px) {
  .dados {
    margin-right: 20px;
  }
}
</style>
