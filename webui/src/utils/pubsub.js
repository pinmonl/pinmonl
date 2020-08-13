import { wsBaseURL, Topics } from './constants'
import Events from 'events'

class Pubsub {
  constructor (token) {
    this.events = new Events()
    this.connect(token)
  }

  connect (token) {
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
    if (!token) {
      return
    }

    const ws = new WebSocket(wsBaseURL + `?token=${token}`)
    ws.onopen = this.handleOpen
    ws.onmessage = this.handleMessage
    this.ws = ws
  }

  handleOpen = () => {
    // console.log('ws opened', event)
    this.subscribe(Topics.PINL_UPDATED)
    this.subscribe(Topics.PINL_DELETED)
  }

  handleMessage = (event) => {
    // console.log('ws message', event)
    const payload = JSON.parse(event.data)
    const data = {
      topic: payload.topic,
      data: payload.data,
    }

    this.events.emit(payload.topic, data)
  }

  subscribe = (topic) => {
    const data = { topic, subscribe: true }
    this.ws.send(JSON.stringify(data))
  }

  unsubscribe = (topic) => {
    const data = { topic, unsubscribe: true }
    this.ws.send(JSON.stringify(data))
  }

  on = (topic, fn) => {
    this.events.on(topic, fn)
  }

  off = (topic, fn) => {
    this.events.off(topic, fn)
  }

  once = (topic, fn) => {
    this.events.once(topic, fn)
  }
}

export default Pubsub
