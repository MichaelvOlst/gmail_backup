import Vue from 'vue'
import Router from 'vue-router'
import store from './../store'


Vue.use(Router)

import Login from '../components/pages/Login.vue';
import Dashboard from '../components/pages/Dashboard.vue';

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
    if (!store.getters.isAuthenticated) {
      next({
        path: '/login',
        query: { redirect: to.fullPath }
      })
    } else {
      next()
    }
  } else {
    next()
  }
})

export function createRouter () {
  return router
}