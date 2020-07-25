import Vue from 'vue'
import Vuex from 'vuex'
import { isProd } from '@/utils/constants'
import getters from './getters'
import mutations from './mutations'
import actions from './actions'
import APIClient from '@/api/client'
import * as auth from '@/utils/auth'
import { SearchParams } from '@/utils/search'

Vue.use(Vuex)

const client = new APIClient()
client.errorHandler = (error) => {
  if (error.status == 401) {
    auth.logout()
    return false
  }
}

const state = {
  user: null,
  token: '',
  search: new SearchParams(),
  client: client,
}

const store = new Vuex.Store({
  strict: !isProd,
  state,
  getters,
  mutations,
  actions,
})

store.commit('SET_TOKEN', auth.getToken())

export default store
