import Vue from 'vue'
import Router from 'vue-router'
import store from './../store'

import { CHECK_AUTH } from './../store/modules/types';

Vue.use(Router)

import Login from '../pages/Login.vue';
import Dashboard from '../pages/Dashboard.vue';

const router = new Router({
  mode: 'history',
  fallback: false,
  routes: [
    { path: '/login', name: 'login', component: Login },
    { path: '/dashboard', name: 'dashboard', component: Dashboard, meta: { requiresAuth: true } },
    { path: '/', redirect: '/dashboard' }
  ]
})

router.beforeEach((to, from, next) => {
  if (to.matched.some(record => record.meta.requiresAuth)) {

    store.dispatch(CHECK_AUTH)
      .then(next)
      .catch(() => {
        next({
          path: '/login',
          query: { redirect: to.fullPath }
        })
      })
  } else {
    next()
  }
})

export function createRouter () {
  return router
}