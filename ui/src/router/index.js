import Vue from "vue";
import VueRouter from "vue-router";
import App from "@/App.vue";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "home",
    component: App
  }
];

const router = new VueRouter({
  routes
});

export default router;
