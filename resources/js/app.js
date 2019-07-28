import Vue from 'vue'
import App from './App.vue'

import Vuetify from 'vuetify'
Vue.use(Vuetify)
import 'vuetify/dist/vuetify.min.css'

window.axios = require('axios');
window.axios.defaults.baseURL = window.location.origin

import store from './store';

import { createRouter } from './router'
const router = createRouter()

const app = new Vue({
    el: '#app',
    router,
    store,
    render: h => h(App),   
})