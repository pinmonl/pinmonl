import { fetch } from '@/pkgs/fetch'

export default {
  namespaced: true,
  state: {
    tags: [],
    parents: {},
  },
  getters: {
    tags: (state) => {
      return state.tags
    },
    new: () => () => {
      return {
        name: '',
        parentId: '',
        sort: 0,
      }
    },
    find: () => (tags, id) => {
      return tags.find(tag => tag.id == id)
    },
    findByName: () => (tags, name) => {
      return tags.find(tag => tag.name == name)
    },
    getByName: (state, getters) => (tags, names) => {
      return names.reduce((carry, name) => {
        const tag = getters.findByName(tags, name)
        if (tag) {
          return [ ...carry, tag ]
        }
        return carry
      }, [])
    },
    mapName: () => (tags) => {
      return tags.map(tag => tag.name)
    },
    search: () => (tags, q) => {
      const escaped = q.replace(/[-[\]{}()*+?.,\\^$|#\s]/g, '\\$&')
      const pattern = new RegExp(escaped, 'i')
      return tags.filter((tag) => {
        return pattern.test(tag.name)
      })
    },
    getByParent: () => (tags, parentId) => {
      return tags.filter(tag => tag.parentId == parentId)
    },
    parents: (state) => {
      return state.parents
    },
  },
  mutations: {
    SET_TAGS (state, tags) {
      state.tags = tags
    },
    UPDATE_TAG (state, tag) {
      const i = state.tags.findIndex(t => t.id == tag.id)
      if (i < 0) {
        return
      }
      state.tags = [
        ...state.tags.slice(0, i),
        tag,
        ...state.tags.slice(i + 1),
      ]
    },
    ADD_TAG (state, tag) {
      state.tags = [ ...state.tags, tag ]
    },
    DELETE_TAG (state, tag) {
      const i = state.tags.findIndex(t => t.id == tag.id)
      if (i < 0) {
        return
      }
      state.tags = [
        ...state.tags.slice(0, i),
        ...state.tags.slice(i + 1),
      ]
    },
    SET_PARENTS (state, { tag = {}, parents = [] }) {
      state.parents[tag.id] = parents
    },
  },
  actions: {
    async fetchAll ({ dispatch, commit }) {
      const tags = await dispatch('list')
      commit('SET_TAGS', tags)
    },
    async list ({ rootGetters }) {
      const req = rootGetters.newRequest('/api/tag')
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
      const req = rootGetters.newRequest(`/api/tag/${id}`)
      const resp = await fetch(req)
      if (resp.status >= 400) {
        throw resp
      }
      if (resp.ok) {
        const tag = await resp.json()
        commit('UPDATE_TAG', tag)
        return tag
      }
    },
    async create ({ rootGetters, commit }, data) {
      const req = rootGetters.newRequest('/api/tag', {
        method: 'POST',
        body: JSON.stringify(data),
      })
      const resp = await fetch(req)
      if (resp.status >= 400) {
        throw resp
      }
      if (resp.ok) {
        const tag = await resp.json()
        commit('ADD_TAG', tag)
        return tag
      }
    },
    async update ({ rootGetters, commit }, data) {
      const { id } = data
      const req = rootGetters.newRequest(`/api/tag/${id}`, {
        method: 'PUT',
        body: JSON.stringify(data),
      })
      const resp = await fetch(req)
      if (resp.status >= 400) {
        throw resp
      }
      if (resp.ok) {
        const tag = await resp.json()
        commit('UPDATE_TAG', tag)
        return tag
      }
    },
    async delete ({ rootGetters, commit }, data) {
      const { id } = data
      const req = rootGetters.newRequest(`/api/tag/${id}`, {
        method: 'DELETE',
      })
      const resp = await fetch(req)
      if (resp.status >= 400) {
        throw resp
      }
      if (resp.ok) {
        commit('DELETE_TAG', data)
        return
      }
    },
    getParents ({ getters, commit }, tag) {
      let parents = getters.parents[tag.id]
      if (!parents) {
        let parentId = tag.parentId
        parents = []
        while (parentId) {
          const parent = getters.find(getters.tags, parentId)
          parents.unshift(parent)
          parentId = parent.parentId
        }
        commit('SET_PARENTS', { tag, parents })
      }
      return parents
    },
  },
}
