import Vue from "vue";
import VueRouter from "vue-router";
import agencyPageContainer from "@/components/agency/agencyPageContainer.vue";
import about from "@/components/about/about.vue";
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
    path: "/orgao/:agencyName/:year/:month",
    name: "agency",
    component: agencyPageContainer,
  },
  {
    path: "/orgao/:agencyName/:year",
    name: "agencyYearContainer",
    component: agencyYearContainer,
  },
  {
    path: "/sobre",
    name: "sobre",
    component: about,
  },
];

const router = new VueRouter({
  mode: "history",
  routes,
});

export default router;
