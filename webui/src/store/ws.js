export default {
  actions: {
    subscribe () {
      const ws = new WebSocket(`ws://${location.host}/ws`)
      ws.onopen = () => {
      }
      ws.onmessage = (event) => {
        console.log(JSON.parse(event.data))
      }
    },
  },
}
