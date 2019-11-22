import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import axios from "axios";

Vue.config.productionTip = false;
const base = axios.create({
  baseURL: 'http://dadosjusbr.com/uiapi/v1'
})
Vue.prototype.$http = base

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");