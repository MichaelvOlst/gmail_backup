import Vue from "vue";
import Vuex from "vuex";

import auth from "./modules/auth";
import accounts from "./modules/accounts";
import settings from "./modules/settings";
import notification from "./modules/notification";


Vue.use(Vuex);

export default new Vuex.Store({
  modules: {
    auth,
    accounts,
    settings,
    notification,
  }
});