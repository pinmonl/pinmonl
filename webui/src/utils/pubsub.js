// import { baseURL } from './constants'

const wsBaseURL = `ws://localhost:3399/ws`

class Pubsub {
  constructor (token) {
    let wsURL = wsBaseURL
    if (token) {
      wsURL += `?token=${token}`
    }

    const ws = new WebSocket(wsURL)
    ws.onopen = this.handleOpen
    ws.onmessage = this.handleMessage
    this.ws = ws
  }

  handleOpen = (...args) => {
    console.log('ws opened', args)
    this.subscribe('pinl_updated')
  }

  handleMessage = (...args) => {
    console.log('ws message', args)
  }

  subscribe = (topic) => {
    const data = { topic }
    this.ws.send(JSON.stringify(data))
  }
}

export default Pubsub
