import {GOOGLE_URL, SAVE_ACCOUNT, ALL_ACCOUNTS, GET_ACCOUNT, DELETE_ACCOUNT} from './types';
import axios from 'axios';

const state = {
    accessTokenURL: "",
    accounts: [],
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

    [SAVE_ACCOUNT]({commit}, data) {
        return new Promise(resolve => {
            axios.post('/api/accounts', data)
                .then(({ data }) => {                                        
                    commit('account_saved', data.result)
                    resolve(data)
                })
                .catch((error) => {
                    commit('account_error', error)
                })
        });
    },

    [ALL_ACCOUNTS]({commit}, data) {
        return new Promise(resolve => {
            axios.get('/api/accounts')
                .then(({ data }) => {                                                            
                    commit('account_get', data.result)
                    resolve(data)
                })
                .catch((error) => {
                    commit('account_error', error)
                })
        });
    },

    [GET_ACCOUNT]({commit}, id) {
        return new Promise(resolve => {
            axios.get(`/api/accounts/${id}`)
                .then(({ data }) => {
                    resolve(data.result)
                })
                .catch((error) => {
                    commit('account_error', error)
                })
        });
    },

    [DELETE_ACCOUNT]({commit}, id) {
        return new Promise(resolve => {
            axios.delete(`/api/accounts/${id}`)
                .then(({ data }) => {
                    resolve(data.result)
                })
                .catch((error) => {
                    commit('account_error', error)
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

    account_saved(state, account) {
        state.accounts = [...state.accounts, account];
    },

    account_get(state, accounts) {
        state.accounts = accounts
    },

    account_error(state, error) {
        console.log(error)
    }
};

export default {
    state,
    actions,
    mutations,
    getters
};