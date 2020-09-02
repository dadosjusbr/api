<template>
  <div v-bind:style="colorNavBar" class="navBar">
    <div class="auxBar">
      <router-link to="/">
        <img
          v-on:click="toHome()"
          src="../../assets/icon_dadosjusbr_colored.svg"
        />
      </router-link>
      <div>
        <router-link to="/dados">
          <b class="dados">
            Dados
          </b>
        </router-link>
        <router-link to="/equipe">
          <b v-on:click="toEquipe()" style="color: white; font-size: 18px;">
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
      equipe: false,
      colorNavBar: {
        backgroundColor: "#3e5363",
      },
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
      // alert(to, from);
      if (this.$route.path == "/equipe") {
        this.colorNavBar.backgroundColor = "#2fbb96";
      } else if (this.$route.path == "/") {
        this.colorNavBar.backgroundColor = "#3e5363";
      } else if (this.$route.path == "/dados") {
        this.colorNavBar.backgroundColor = "#3e5363";
      } else {
        this.colorNavBar.backgroundColor = "white";
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
  color: #ffffff;
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
  .navBarHome {
    padding: 20px 20px 40px 20px;
  }
  .navBarEquipe {
    padding: 20px 20px 40px 20px;
  }
  .auxBarEquipe {
    padding: 40px 20px 40px 20px;
  }
  .auxBarHome {
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
