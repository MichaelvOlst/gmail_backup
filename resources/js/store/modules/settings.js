import {
    GET_SETTINGS,
    SAVE_SETTINGS
  } from './types';
  import axios from 'axios';
  
  const state = {
    settings: {
      storage_options: []
    },
    error: {},
  };
  
  const getters = {
   
  };
  
  const actions = {
    async [GET_SETTINGS]({ commit }) {
        try {
            const {data} = await axios.get('/api/settings')
            commit('set_settings', data.result)
            return data.result
        } catch (error) {
            commit('set_settings_error', error.response.data.error)
            throw error.response.data.error
        }
    },
  };
  
  const mutations = {
  
    set_settings(state, settings) {
      state.settings = settings
    },
  
    set_settings_error(state, error) {
      state.error = error
    },
  
  };
  
  export default {
    state,
    actions,
    mutations,
    getters
  };