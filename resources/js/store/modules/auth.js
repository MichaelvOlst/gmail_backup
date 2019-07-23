import {
  LOGIN,
  LOGOUT,
  CHECK_AUTH
} from './types';
// import { SET_AUTH, PURGE_AUTH, SET_ERROR } from './mutations.type';

const state = {
  errors: null,
  isAuthenticated: true
};

const getters = {
  isAuthenticated(state) {
    return state.isAuthenticated;
  }
};

const actions = {
  [LOGIN](context, credentials) {
    return new Promise(resolve => {
    //   ApiService.post("users/login", { user: credentials })
    //     .then(({ data }) => {
    //       context.commit(SET_AUTH, data.user);
    //       resolve(data);
    //     })
    //     .catch(({ response }) => {
    //       context.commit(SET_ERROR, response.data.errors);
    //     });
    });
  },
  [LOGOUT](context) {
    // context.commit(PURGE_AUTH);
  },
 
  [CHECK_AUTH](context) {
    // if (JwtService.getToken()) {
    //   ApiService.setHeader();
    //   ApiService.get("user")
    //     .then(({ data }) => {
    //       context.commit(SET_AUTH, data.user);
    //     })
    //     .catch(({ response }) => {
    //       context.commit(SET_ERROR, response.data.errors);
    //     });
    // } else {
    //   context.commit(PURGE_AUTH);
    // }
  },
};

const mutations = {
//   [SET_ERROR](state, error) {
//     state.errors = error;
//   },
//   [SET_AUTH](state, user) {
//     state.isAuthenticated = true;
//     state.errors = {};
//   },
//   [PURGE_AUTH](state) {
//     state.isAuthenticated = false;
//     state.errors = {};
//   }
};

export default {
  state,
  actions,
  mutations,
  getters
};