<template>
  <div :class="[equipe ? 'navBarEquipe' : '', 'navBarHome']">
    <div :class="[equipe ? 'auxBarEquipe' : '', 'auxBarHome']">
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
};
</script>

<style scoped>
.dados {
  color: white;
  font-size: 18px;
  margin-right: 90px;
}

.navBarHome {
  background-color: #3e5363;
  height: 185px;
  padding: 50px 100px 0px 100px;
  color: #ffffff;
}

.navBarEquipe {
  background-color: #2fbb96;
  color: #3e5363;
  height: 185px;
  padding: 50px 100px 0px 100px;
}

.auxBarHome {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  border-bottom: 1px solid white;
  padding-right: 100px;
  padding-bottom: 35px;
  align-items: center;
}

.auxBarEquipe {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  border-bottom: 1px solid #3e5363;
  padding-right: 100px;
  padding-bottom: 35px;
  align-items: center;
}
@media only screen and (max-width: 500px) {
  .navBarHome {
    padding: 40px 20px 40px 20px;
  }
  .navBarEquipe {
    padding: 40px 20px 40px 20px;
    background-color: coral;
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
