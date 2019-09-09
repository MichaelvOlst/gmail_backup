import {
    GET_SETTINGS,
    SAVE_SETTINGS
  } from './types';
  import axios from 'axios';
  
  const state = {
    settings: {
      storage_options: {},
    },
    active_storage_options: [],
    error: null,
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

    async [SAVE_SETTINGS]({ commit }, settings) {
      try {
          const {data} = await axios.post('/api/settings', settings)
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

      let storage_options = []
      for (let provider in settings.storage_options) {        
        const element = settings.storage_options[provider];
        storage_options.push({option: element.StorageOption.option, name: element.StorageOption.name}); 
      }

      state.active_storage_options = storage_options
    },
  
    set_settings_error(state, error) {
      state.error = error
    },
  
  };
  
  export default {
    state,
    actions,
    mutations,
  };