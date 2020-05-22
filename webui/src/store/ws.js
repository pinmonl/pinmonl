export default {
  actions: {
    subscribe ({ dispatch }) {
      const ws = new WebSocket(`ws://${location.host}/ws`)
      ws.onopen = () => {
      }
      ws.onmessage = (e) => {
        const msg = JSON.parse(e.data)
        dispatch('handleTopic', msg)
      }
    },
    handleTopic ({ getters, commit }, { topic, data }) {
      if (['pinl.create', 'pinl.update'].includes(topic)) {
        const found = getters['pinl/find'](getters['pinl/pinls'], data.id)
        if (found) {
          commit('pinl/UPDATE_PINL', data)
        } else {
          commit('pinl/ADD_PINL', data)
        }
      }
    },
  },
}
