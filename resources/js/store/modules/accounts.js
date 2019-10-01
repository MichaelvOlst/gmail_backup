import {GOOGLE_URL, SAVE_ACCOUNT, ALL_ACCOUNTS, GET_ACCOUNT, DELETE_ACCOUNT, BACKUP_ACCOUNT} from './types';
import axios from 'axios';

const state = {
    tokenURL: null,
    accounts: [],
    errors: null
};

const getters = {
    getTokenURL(state) {
        return state.tokenURL;
    }
};
 
const actions = {
    async [GOOGLE_URL]({ commit }) {
        try {
            const {data} = await axios.get('/api/google-url')
            commit('google_token_url', data.result)
            return data.result
        } catch (error) {
            commit('google_token_error', error.response.data.error)
            throw error.response.data.error
        }
    },

    async [SAVE_ACCOUNT]({commit}, data) {
        try {
            const response = await axios.post('/api/accounts', data)
            commit('account_saved', response.result)
            return response
        } catch(error) {
            commit('account_error', error.response.data.error)
            throw error.response.data.error
        }
    },

    async [BACKUP_ACCOUNT]({commit}, id) {

        try {
            const response = await axios.post(`/api/backup/${id}`)
            // commit('account_saved', response.result)
            return response
        } catch(error) {
            // commit('account_error', error.response.data.error)
            throw error.response.data.error
        }
    },

    async [ALL_ACCOUNTS]({commit}, data) {
        try {
            const {data} = await axios.get('/api/accounts')
            commit('account_get', data.result)
            return data.result
        } catch(error) {
            commit('account_error', error.response.data.error)
            throw error.response.data.error
        }
    },

    async [GET_ACCOUNT]({commit}, id) {

        try {
            const {data} = await axios.get(`/api/accounts/${id}`)
            return data.result
        } catch (error) {
            commit('account_error', error.response.data.error)
            throw error.response.data.error
        }
    },

    async [DELETE_ACCOUNT]({commit}, id) {

        try {
            const {data} = await axios.delete(`/api/accounts/${id}`)
            return data.result
        } catch (error) {
            commit('account_error', error.response.data.error)
            throw error.response.data.error
        }
    }
    
};

const mutations = {

    google_token_url(state, url) {
        state.tokenURL = url
    },

    google_token_error(state) {
        state.tokenURL = ""
    },

    account_saved(state, account) {
        // state.accounts = [...state.accounts, account];
    },

    account_get(state, accounts) {
        state.accounts = accounts
    },

    account_error(state, errors) {
        state.errors = errors
    }
};

export default {
    state,
    actions,
    mutations,
    getters
};