<template>
  <div class="socialMidiaShare">
    <button class="button" :href="this.fileUrl" target="_blank">
      Compartilhar
    </button>
    <!-- <facebook
      style="margin-right: 5px"
      :url="this.url"
      scale="2"
      :title="this.socialMidiaMsg"
    ></facebook>
    <whats-app
      style="margin-right: 5px"
      :url="this.url"
      :title="this.socialMidiaMsg"
      scale="2"
    ></whats-app>
    <twitter
      style="margin-right: 5px"
      :url="this.url"
      :title="this.socialMidiaMsg"
      scale="2"
    ></twitter>
    <email
      style="margin-right: 5px"
      :url="this.url"
      :subject="this.socialMidiaMsg"
      scale="2"
    ></email> -->
  </div>
</template>

<script>
import { Facebook, Twitter, WhatsApp, Email } from "vue-socialmedia-share";
export default {
  name: "socialMediaShare",
  data() {
    return {
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
  components: {
    Facebook,
    Twitter,
    WhatsApp,
    Email,
  },
};
</script>

<style>
.socialMidiaShare {
  text-align: center;
  margin-left: 10px;
  width: 150px;
  height: 48px;
  background-color: red;
  order: 2;
}
.button {
  width: 150px;
  height: 48px;
  background-color: #545454;
  border: solid #545454;
  color: white;
}
</style>
