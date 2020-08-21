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
          <b style="color: white; font-size: 18px; margin-right: 90px;">
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
  justify-content: space-between;
  border-bottom: 1px solid white;
  padding-right: 100px;
  padding-bottom: 35px;
  align-items: center;
}
.auxBarEquipe {
  display: flex;
  justify-content: space-between;
  border-bottom: 1px solid #3e5363;
  padding-right: 100px;
  padding-bottom: 35px;
  align-items: center;
}
</style>
