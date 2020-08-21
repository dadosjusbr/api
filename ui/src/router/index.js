import Vue from "vue";
import VueRouter from "vue-router";
import agencyPageContainer from "@/components/agency/agencyPageContainer.vue";
import equipe from "@/components/about/equipe.vue";
import agencyYearContainer from "@/components/agency-year-month/agencyYearContainer.vue";
import homePage from "@/components/home-page/homePage.vue";
import statePageContainer from "@/components/state/statePageContainer.vue";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "home",
    component: homePage,
  },
  {
    path: "/dados",
    name: "statePage",
    component: statePageContainer,
  },
  {
    path: "/equipe",
    name: "equipe",
    component: equipe,
  },
  {
    path: "/orgao/:agencyName/:year/:month",
    name: "agency",
    component: agencyPageContainer,
  },
  {
    path: "/orgao/:agencyName/:year",
    name: "agencyYearContainer",
    component: agencyYearContainer,
  },
];

const router = new VueRouter({
  mode: "history",
  routes,
});

export default router;
