import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

import auth from './auth'
import http from './http'
import keybinding from './keybinding'
import modules from './modules'
import ws from './ws'

export default new Vuex.Store({
  state: {
    ready: false,
    showNav: false,
  },
  mutations: {
    SET_READY (state, ready) {
      state.ready = ready
    },
    SET_NAV (state, show) {
      state.showNav = show
    },
  },
  actions: {
  },
  modules: {
    auth,
    http,
    keybinding,
    ws,
    ...modules,
  },
})
