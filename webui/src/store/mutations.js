const mutations = {
  SET_TOKEN (state, token) {
    state.token = token
    state.client.setToken(token)
  },
  SET_USER (state, user) {
    state.user = user
  },
  SET_SEARCH (state, search) {
    state.search = search
  },
  SET_PINS (state, pins) {
    state.pins = pins
  },
  SET_TAGS (state, tags) {
    state.tags = tags
  },
}

export default mutations
