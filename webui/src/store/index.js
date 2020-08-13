import Vue from 'vue'
import Vuex from 'vuex'
import getters from './getters'
import mutations from './mutations'
import actions from './actions'
import APIClient from '@/api/client'
import * as auth from '@/utils/auth'
import { SearchParams } from '@/utils/search'
import Pubsub from '@/utils/pubsub'

Vue.use(Vuex)

const client = new APIClient()
client.errorHandler = (error) => {
  if (error.status == 401) {
    auth.logout()
    return false
  }
}

const ws = new Pubsub()

const state = {
  user: null,
  token: '',
  search: new SearchParams(),
  client: client,
  socket: ws,
}

const store = new Vuex.Store({
  state,
  getters,
  mutations,
  actions,
})

store.commit('SET_TOKEN', auth.getToken())

export default store
