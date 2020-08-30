import React, {
  useCallback,
  useRef,
  useEffect,
} from 'react'
import {
  CRUD_GET_LIST_SUCCESS,
  GET_LIST,
  FETCH_END,
  CRUD_DELETE_SUCCESS,
  CRUD_DELETE,
} from 'react-admin'
import { useDispatch } from 'react-redux'
import PubsubContext from './PubsubContext'
import { baseWs } from '@/utils/constants'
import EventEmitter from 'events'
import { useAuthToken } from '@/auth'

const Pubsub = ({ children }) => {
  const hub = useRef()
  const bus = useRef(new EventEmitter())
  const { token } = useAuthToken()
  const dispatch = useDispatch()

  const onOpen = useCallback((e) => {
    bus.current.emit('open', e)
  }, [])

  const onClose = useCallback((e) => {
    bus.current.emit('close', e)
  }, [])

  const onMessage = useCallback((e) => {
    const { topic, data } = JSON.parse(e.data)

    bus.current.emit(`topic:${topic}`, data)
    bus.current.emit('message', e)
  }, [])

  const onError = useCallback((e) => {
    bus.current.emit('error', e)
  }, [])

  const close = useCallback(() => {
    if (!hub.current) return
    hub.current.close()
    hub.current = null
  }, [])

  const connect = useCallback((token) => {
    const ws = new WebSocket(`${baseWs}?token=${token}`)
    ws.onopen = onOpen
    ws.onclose = onClose
    ws.onerror = onError
    ws.onmessage = onMessage
    hub.current = ws
  }, [onOpen, onClose, onError, onMessage])

  const on = useCallback((name, fn) => {
    bus.current.on(name, fn)
  }, [])

  const off = useCallback((name, fn) => {
    bus.current.off(name, fn)
  }, [])

  const once = useCallback((name, fn) => {
    bus.current.once(name, fn)
  }, [])

  const subscribe = useCallback((topic) => {
    if (!hub.current) return

    const data = { topic }
    hub.current.send(JSON.stringify(data))
  }, [])

  const unsubscribe = useCallback(() => {
  }, [])

  useEffect(() => {
    if (!token) return

    const handlePinlUpdated = (record) => {
      dispatch({
        type: CRUD_GET_LIST_SUCCESS,
        payload: {
          data: [record],
          total: 1,
        },
        meta: {
          resource: 'pinl',
          fetchResponse: GET_LIST,
          fetchStatus: FETCH_END,
        },
      })
    }

    const handlePinlDeleted = (record) => {
      dispatch({
        type: CRUD_DELETE_SUCCESS,
        payload: {
          data: record,
        },
        meta: {
          resource: 'pinl',
          refresh: true,
          fetchResponse: CRUD_DELETE,
          fetchStatus: FETCH_END,
        },
      })
    }

    const handleOpen = () => {
      subscribe('pinl_updated')
      subscribe('pinl_deleted')
      on('topic:pinl_updated', handlePinlUpdated)
      on('topic:pinl_deleted', handlePinlDeleted)
    }

    connect(token)
    once('open', handleOpen)

    return () => {
      close()
      off('open', handleOpen)
      off('topic:pinl_updated', handlePinlUpdated)
      off('topic:pinl_deleted', handlePinlDeleted)
    }
  }, [token, connect, subscribe, once, off])

  const ctx = {
    on,
    off,
    once,
    subscribe,
    unsubscribe,
  }

  return (
    <PubsubContext.Provider value={ctx}>
      {children}
    </PubsubContext.Provider>
  )
}

export default Pubsub
