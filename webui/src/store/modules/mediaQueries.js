import Vue from 'vue'

export default {
  namespaced: true,
  state: {
    queries: {},
    matches: {},
  },
  actions: {
    match ({ state, dispatch }, query) {
      if (!state.queries[query]) {
        const m = window.matchMedia(query)
        Vue.set(state.queries, query, m)
        dispatch('handleMatch', { query, m })
        m.addListener((m) => {
          dispatch('handleMatch', { query, m })
        })
      }
      return state.queries[query]
    },
    handleMatch ({ state }, { query, m }) {
      Vue.set(state.matches, query, m.matches)
    },
  },
}
