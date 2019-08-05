import {GOOGLE_URL} from './types';
import axios from 'axios';

const state = {
    accessTokenURL: ""
};

const getters = {
    getAccessTokenURL(state) {
        return state.accessTokenURL;
    }
};

const actions = {
    [GOOGLE_URL]({ commit }) {
        return new Promise(resolve => {
            axios.get('/api/google-url')
                .then(({ data }) => {                                        
                    commit('accesstoken_url', data.result)
                    resolve(data)
                })
                .catch(() => {
                    commit('accesstoken_error')
                })
        });
    },
    
};

const mutations = {

    accesstoken_url(state, url) {
        state.accessTokenURL = url
    },

    accesstoken_error(state) {
        state.accessTokenURL = ""
    },
};

export default {
    state,
    actions,
    mutations,
    getters
};