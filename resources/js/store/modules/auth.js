import {
  LOGIN,
  LOGOUT,
  CHECK_AUTH
} from './types';
// import { SET_AUTH, PURGE_AUTH, SET_ERROR } from './mutations.type';
import axios from 'axios';

const state = {
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
      axios.post('/auth/login', data)
        .then(({data}) => {
          commit('authLogin')
          resolve(data)
        })
        .catch( () => {
          commit('authError')
        })
    });
  },
  [LOGOUT]({commit}) {
    return new Promise(resolve => {
      axios.post('/auth/logout')
        .then(({data}) => {
          commit('authLogout')
          resolve(data)
        })
        .catch( () => {
          commit('authError')
        })
    });
  },
 
  [CHECK_AUTH]({commit}) {    
    return new Promise((resolve, reject) => {
      // commit('auth_request')
      axios.get('/auth/session')
        .then(response => {
          if(response.data.result === false) {
            commit('authError')
            reject(response)
            return
          }

          localStorage.setItem('loggedIn', true)
          resolve(response)
          commit('authLogin')
        })
        .catch(err => {
          commit('authError')
          localStorage.removeItem('loggedIn')
          reject(err)
        })
    })
  },
};

const mutations = {

  authLogin(state) {
    state.isAuthenticated = true
  },

  authLogout(state) {
    state.isAuthenticated = false
  },

  authError(state) {
    state.isAuthenticated = false
  }
};

export default {
  state,
  actions,
  mutations,
  getters
};