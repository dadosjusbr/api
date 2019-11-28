import Vue from "vue";
import VueRouter from "vue-router";
import Home from "@/Home.vue";
import agencyPageContainer from "@/components/agency/agencyPageContainer.vue";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "home",
    component: Home
  },
  {
    path: "/orgao/:agencyName",
    name: "agency",
    component: agencyPageContainer
  }
];

const router = new VueRouter({
  routes
});

export default router;
