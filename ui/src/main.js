import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import axios from "axios";
import VueMaterial from "vue-material";
import VueGtag from "vue-gtag";
import "vue-material/dist/vue-material.min.css";
import "vue-material/dist/vue-material.css";
import "vue-material/dist/vue-material.min.css";
import "bootstrap/dist/css/bootstrap.min.css";

Vue.config.productionTip = false;
Vue.use(VueMaterial);
Vue.use(VueGtag, {
  config: { id: "UA-143064237-1" }
});

const base = axios.create({
  baseURL: process.env.VUE_APP_API_BASE_URL
});
Vue.prototype.$http = base;

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");
