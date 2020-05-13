export default {
  actions: {
    subscribe () {
      const ws = new WebSocket(`ws://${location.host}/ws`)
      ws.onopen = () => {
      }
      ws.onmessage = () => {
      }
    },
  },
}
