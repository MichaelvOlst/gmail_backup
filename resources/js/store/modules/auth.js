import {
  LOGIN,
  LOGOUT,
  CHECK_AUTH
} from './types';
// import { SET_AUTH, PURGE_AUTH, SET_ERROR } from './mutations.type';
import axios from 'axios';

const state = {
  errors: null,
  isAuthenticated: false
};

const getters = {
  isAuthenticated(state) {
    return state.isAuthenticated;
  }
};

const actions = {
  [LOGIN]({commit}, data) {
    return new Promise(resolve => {
      console.log(data);
      axios.post('/auth/login', data)
        .then(({data}) => {
          commit('auth_success')
          resolve(data)
        })
        .catch( () => {
          commit('auth_error')
        })
    });
  },
  [LOGOUT](context) {
    // context.commit(PURGE_AUTH);
  },
 
  [CHECK_AUTH]({commit}) {    
    return new Promise((resolve, reject) => {
      // commit('auth_request')
      axios.get('/auth/session')
        .then(response => {          
          if(response.data.result === false) {
            commit('auth_error')
            reject(response)
            return
          }

          localStorage.setItem('loggedIn', true)
          resolve(response)
          commit('auth_success')
        })
        .catch(err => {
          commit('auth_error')
          localStorage.removeItem('loggedIn')
          reject(err)
        })
    })
  },
};

const mutations = {

  auth_succes(state) {
    state.isAuthenticated = true;
  },

  auth_error(state) {
    state.isAuthenticated = false;
    state.errors = {};
  }
};

export default {
  state,
  actions,
  mutations,
  getters
};