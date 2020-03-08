import { fetch } from '@/pkgs/fetch'

export default {
  state: {
    user: null,
  },
  getters: {
    authed (state, getters) {
      return !!getters.user
    },
    user (state) {
      return state.user
    },
  },
  mutations: {
    SET_USER (state, user) {
      state.user = user
    },
  },
  actions: {
    async login ({ getters, dispatch }, data) {
      const req = getters.newRequest('/api/session', {
        method: 'POST',
        body: JSON.stringify(data),
      })
      const resp = await fetch(req)
      if (resp.status >= 400) {
        throw resp
      }
      if (resp.ok) {
        await dispatch('getMe')
      }
    },
    async logout ({ getters, commit }) {
      const req = getters.newRequest('/api/session', {
        method: 'DELETE',
      })
      const resp = await fetch(req)
      if (resp.status >= 400) {
        throw resp
      }
      if (resp.ok) {
        commit('SET_USER', null)
      }
    },
    async getMe ({ getters, commit }) {
      const req = getters.newRequest('/api/me')
      const resp = await fetch(req)
      if (resp.status >= 400) {
        throw resp
      }
      if (resp.ok) {
        const user = await resp.json()
        commit('SET_USER', user)
      }
    },
    async updateMe ({ getters, commit }, data) {
      const req = getters.newRequest('/api/me', {
        method: 'PUT',
        body: JSON.stringify(data),
      })
      const resp = await fetch(req)
      if (resp.status >= 400) {
        throw resp
      }
      if (resp.ok) {
        const user = await resp.json()
        commit('SET_USER', user)
      }
    },
    async signup ({ dispatch }, data) {
      await dispatch('createUser', data)
      await dispatch('login', data)
    },
    async createUser ({ getters }, data) {
      const req = getters.newRequest('/api/user', {
        method: 'POST',
        body: JSON.stringify(data),
      })
      const resp = await fetch(req)
      if (resp.status >= 400) {
        throw resp
      }
      if (resp.ok) {
        return await resp.json()
      }
    },
  },
}