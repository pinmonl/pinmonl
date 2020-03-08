import { fetch } from '@/pkgs/fetch'

export default {
  namespaced: true,
  state: {
  },
  getters: {
  },
  mutations: {
  },
  actions: {
    async find ({ rootGetters }, data) {
      const { user, name } = data
      const req = rootGetters.newRequest(`/api/sharing/${user}/${name}`)
      const resp = await fetch(req)
      if (resp.status >= 400) {
        throw resp
      }
      if (resp.ok) {
        return await resp.json()
      }
    },
    async listPinls ({ rootGetters }, data) {
      const { user, name } = data
      const req = rootGetters.newRequest(`/api/sharing/${user}/${name}/pinl`)
      const resp = await fetch(req)
      if (resp.status >= 400) {
        throw resp
      }
      if (resp.ok) {
        return await resp.json()
      }
    },
    async listTags ({ rootGetters }, data) {
      const { user, name } = data
      const req = rootGetters.newRequest(`/api/sharing/${user}/${name}/tag`)
      const resp = await fetch(req)
      if (resp.status >= 400) {
        throw resp
      }
      if (resp.ok) {
        return await resp.json()
      }
    },
    async findPinl ({ rootGetters }, data) {
      const { user, name, id } = data
      const req = rootGetters.newRequest(`/api/sharing/${user}/${name}/pinl/${id}`)
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
