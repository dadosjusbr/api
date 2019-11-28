import Vue from "vue";
import VueRouter from "vue-router";
import App from "@/App.vue";
import statePageContainer from "@/components/state/statePageContainer.vue";
import agencyPageContainer from "@/components/agency/agencyPageContainer.vue";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "home",
    component: App
  },
  {
    path: "/paraiba",
    name: "state",
    component: statePageContainer
  },
  {
    path: "/tjpb",
    name: "agency",
    component: agencyPageContainer
  }
];

const router = new VueRouter({
  routes
});

export default router;
