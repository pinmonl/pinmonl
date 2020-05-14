export default {
  state: {
    globalSearch: false,
  },
  getters: {
    globalSearch (state) {
      return state.globalSearch
    },
  },
  mutations: {
    SET_GLOBAL_SEARCH (state, val) {
      state.globalSearch = val
    },
  },
}