import { Request } from '@/pkgs/fetch'

function newRequest(input, opts) {
  return new Request(input, opts)
}

export default {
  state: {
  },
  getters: {
    newRequest: () => (input, opts) => {
      return newRequest(input, opts)
    },
  },
  mutations: {
  },
  actions: {
  },
}
