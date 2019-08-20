import {NOTIFY} from './types';

const state = {
    message: ""
}

const mutations = {
   [NOTIFY](state, message) {
       state.message = message
   }
}

export default {
    state,
    mutations
}