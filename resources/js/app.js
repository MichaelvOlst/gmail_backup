import Vue from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify';


window.axios = require('axios');
window.axios.defaults.baseURL = window.location.origin

import store from './store';

import { createRouter } from './router'
const router = createRouter()

Vue.config.productionTip = false


const app = new Vue({
    el: '#app',
    router,
    store,
    vuetify,
    render: h => h(App),   
})