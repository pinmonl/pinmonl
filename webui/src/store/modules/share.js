import { fetch } from '@/pkgs/fetch'

export default {
  namespaced: true,
  state: {
    shares: [],
  },
  getters: {
    shares: (state) => {
      return state.shares
    },
    new: () => () => ({
      name: '',
      description: '',
      mustTags: [],
      anyTags: [],
    }),
    find: () => (shares, id) => {
      return shares.find(share => share.id == id)
    },
  },
  mutations: {
    SET_SHARES (state, shares) {
      state.shares = shares
    },
    UPDATE_SHARE (state, share) {
      const i = state.shares.findIndex(s => s.id == share.id)
      if (i < 0) {
        return
      }
      state.shares = [
        ...state.shares.slice(0, i),
        share,
        ...state.shares.slice(i + 1),
      ]
    },
    ADD_SHARE (state, share) {
      state.shares = [ ...state.shares, share ]
    },
    DELETE_SHARE (state, share) {
      const i = state.shares.findIndex(s => s.id == share.id)
      if (i < 0) {
        return
      }
      state.shares = [
        ...state.shares.slice(0, i),
        ...state.shares.slice(i + 1),
      ]
    },
  },
  actions: {
    async fetchAll ({ dispatch, commit }) {
      const shares = await dispatch('list')
      commit('SET_SHARES', shares)
    },
    async list ({ rootGetters }) {
      const req = rootGetters.newRequest(`/api/share`)
      const resp = await fetch(req)
      if (resp.status >= 400) {
        throw resp
      }
      if (resp.ok) {
        return await resp.json()
      }
    },
    async find ({ rootGetters, commit }, data) {
      const { id } = data
      const req = rootGetters.newRequest(`/api/share/${id}`)
      const resp = await fetch(req)
      if (resp.status >= 400) {
        throw resp
      }
      if (resp.ok) {
        const share = await resp.json()
        commit('UPDATE_SHARE', share)
        return share
      }
    },
    async create ({ rootGetters, commit }, data) {
      const req = rootGetters.newRequest(`/api/share`, {
        method: 'POST',
        body: JSON.stringify(data),
      })
      const resp = await fetch(req)
      if (resp.status >= 400) {
        throw resp
      }
      if (resp.ok) {
        const share = await resp.json()
        commit('ADD_SHARE', share)
        return share
      }
    },
    async update ({ rootGetters, commit }, data) {
      const { id } = data
      const req = rootGetters.newRequest(`/api/share/${id}`, {
        method: 'PUT',
        body: JSON.stringify(data),
      })
      const resp = await fetch(req)
      if (resp.status >= 400) {
        throw resp
      }
      if (resp.ok) {
        const share = await resp.json()
        commit('UPDATE_SHARE', share)
        return share
      }
    },
    async delete ({ rootGetters, commit }, data) {
      const { id } = data
      const req = rootGetters.newRequest(`/api/share/${id}`, {
        method: 'DELETE',
      })
      const resp = await fetch(req)
      if (resp.status >= 400) {
        throw resp
      }
      if (resp.ok) {
        commit('DELETE_SHARE', data)
        return
      }
    },
  },
}
