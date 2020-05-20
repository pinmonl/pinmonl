import { fetch } from '@/pkgs/fetch'

export default {
  namespaced: true,
  state: {
    pinls: [],
  },
  getters: {
    pinls: (state) => {
      return state.pinls
    },
    new: () => () => ({
      url: '',
      title: '',
      description: '',
      readme: '',
      tags: [],
      imageId: '',
    }),
    getByTag: () => (pinls, tags) => {
      return pinls.filter(p => {
        return tags.reduce((matched, tag) => {
          return matched && p.tags.includes(tag)
        }, true)
      })
    },
    find: () => (pinls, id) => {
      return pinls.find(p => p.id == id)
    },
    searchByTitle: () => (pinls, q) => {
      return pinls.filter(pinl => pinl.title.includes(q))
    },
    parseSearch: (state, getters, rootState, rootGetters) => (str) => {
      const params = new URLSearchParams(str)
      const search = { input: '', tags: [] }
      if (params.has('input')) {
        search.input = params.get('input')
      }
      if (params.has('tags')) {
        const rawTags = params.get('tags')
        if (rawTags != '') {
          const allTags = rootGetters['tag/tags']
          const tagNames = rawTags.split(',')
          const tags = rootGetters['tag/getByName'](allTags, tagNames)
          search.tags = tags
        }
      }
      return search
    },
    composeSearch: (state, getters, rootState, rootGetters) => (search) => {
      const { input, tags = [] } = search
      const tagNames = rootGetters['tag/mapName'](tags)
      const params = new URLSearchParams({
        input: input,
        tags: tagNames.join(','),
      })
      return params.toString()
    },
  },
  mutations: {
    SET_PINLS (state, pinls) {
      state.pinls = pinls
    },
    UPDATE_PINL (state, pinl) {
      const i = state.pinls.findIndex(p => p.id == pinl.id)
      if (i < 0) {
        return
      }
      state.pinls = [
        ...state.pinls.slice(0, i),
        pinl,
        ...state.pinls.slice(i + 1),
      ]
    },
    ADD_PINL (state, pinl) {
      state.pinls = [ ...state.pinls, pinl ]
    },
    DELETE_PINL (state, pinl) {
      const i = state.pinls.findIndex(p => p.id == pinl.id)
      if (i < 0) {
        return
      }
      state.pinls = [
        ...state.pinls.slice(0, i),
        ...state.pinls.slice(i + 1),
      ]
    },
  },
  actions: {
    async fetchAll ({ dispatch, commit }) {
      const pinls = await dispatch('list')
      commit('SET_PINLS', pinls)
    },
    async list ({ rootGetters }) {
      const req = rootGetters.newRequest('/api/pinl')
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
      const req = rootGetters.newRequest(`/api/pinl/${id}`)
      const resp = await fetch(req)
      if (resp.status >= 400) {
        throw resp
      }
      if (resp.ok) {
        const pinl = await resp.json()
        commit('UPDATE_PINL', pinl)
        return pinl
      }
    },
    async create ({ rootGetters, commit }, data) {
      const req = rootGetters.newRequest('/api/pinl', {
        method: 'POST',
        body: JSON.stringify(data),
      })
      const resp = await fetch(req)
      if (resp.status >= 400) {
        throw resp
      }
      if (resp.ok) {
        const pinl = await resp.json()
        commit('ADD_PINL', pinl)
        return pinl
      }
    },
    async update ({ rootGetters, commit }, data) {
      const { id } = data
      const req = rootGetters.newRequest(`/api/pinl/${id}`, {
        method: 'PUT',
        body: JSON.stringify(data),
      })
      const resp = await fetch(req)
      if (resp.status >= 400) {
        throw resp
      }
      if (resp.ok) {
        const pinl = await resp.json()
        commit('UPDATE_PINL', pinl)
        return pinl
      }
    },
    async delete ({ rootGetters, commit }, data) {
      const { id } = data
      const req = rootGetters.newRequest(`/api/pinl/${id}`, {
        method: 'DELETE',
      })
      const resp = await fetch(req)
      if (resp.status >= 400) {
        throw resp
      }
      if (resp.ok) {
        commit('DELETE_PINL', data)
        return
      }
    },
    openLink (ctx, { pinl, newTab = true }) {
      const { url } = pinl
      window.open(url, newTab ? '_blank' : '')
    },
  },
}