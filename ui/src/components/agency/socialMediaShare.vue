<template>
  <b-row>
    <b-col class="buttonDiv">
      <img
        style="height: 83px; width:295px"
        src="../../assets/button_compartilhar.svg"
        v-on:click="showMenu()"
      />
    </b-col>
    <div class="shareMenu" v-show="this.displayMenu">
      <button class="closeButton" v-on:click="hideMenu()">X</button>
      <h3 style="margin-top: 50px;">Compartilhar</h3>
      <facebook
        v-on:click="hideMenu()"
        style="margin-right: 5px"
        :url="this.url"
        scale="3"
        :title="this.socialMidiaMsg"
      ></facebook>
      <whats-app
        style="margin-right: 5px"
        :url="this.url"
        :title="this.socialMidiaMsg"
        scale="3"
      ></whats-app>
      <twitter
        style="margin-right: 5px"
        :url="this.url"
        :title="this.socialMidiaMsg"
        scale="3"
      ></twitter>
      <email
        style="margin-right: 5px"
        :url="this.url"
        :subject="this.socialMidiaMsg"
        scale="3"
      ></email>
    </div>
    <div class="overlay"></div>
  </b-row>
</template>

<script>
import { Facebook, Twitter, WhatsApp, Email } from "vue-socialmedia-share";
export default {
  name: "socialMediaShare",
  data() {
    return {
      displayMenu: false,
      socialMidiaMsg:
        "Descubra como é a distribuição das remunerações dos funcionários do " +
        this.$route.params.agencyName.toUpperCase() +
        " no ano e mês " +
        this.$route.params.year +
        "/" +
        this.$route.params.month,
      url:
        "https://dadosjusbr.org/orgao/" +
        this.$route.params.agencyName +
        "/" +
        this.$route.params.year +
        "/" +
        this.$route.params.month,
    };
  },
  methods: {
    showMenu() {
      this.displayMenu = true;
    },
    hideMenu() {
      this.displayMenu = false;
    },
  },
  components: {
    Facebook,
    Twitter,
    WhatsApp,
    Email,
  },
};
</script>

<style>

.button {
  width: 150px;
  height: 48px;
  background-color: #545454;
  border: solid #545454;
  color: white;
  font-size: 17px;
}

.closeButton {
  float: right;
  width: 30px;
  height: 30px;
  background-color: #545454;
  border: solid #545454;
  color: white;
}

.shareMenu {
  width: 285px;
  height: 218px;
  background-color: #e0dfe0;
  text-align: center;
  position: fixed;
  left: 50%;
  top: 50%;
  margin-left: -142.5px;
  margin-top: -109px;
}
.overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 200px;
  background-color: tomato;
}
</style>
