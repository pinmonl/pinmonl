import Emitter from 'events'
import { doRequest } from '../data'
import loadToken from './loadToken'
import createAutoAuthProvider from './createAutoAuthProvider'
import createCredentialAuthProvider from './createCredentialAuthProvider'

const createAuthProvider = (baseURL, hasDefaultUser) => {
  const emitter = new Emitter()

  const saveToken = ({ token = null, expireAt = null } = {}) => {
    localStorage.setItem('token', JSON.stringify(token))
    localStorage.setItem('expire_at', JSON.stringify(expireAt))
    emitter.emit('token:update', { token, expireAt })
  }

  const addListener = (name, handler) => {
    emitter.on(name, handler)
    return () => emitter.off(name, handler)
  }

  const refreshToken = async () => {
    const { json } = await doRequest(`${baseURL}/api/refresh`, {
      method: 'POST',
    })
    return json.data
  }

  const init = () => {
    const { token, expireAt } = loadToken()
    if (!token) {
      return
    }
    emitter.emit('token:update', { token, expireAt })
  }

  const utils = {
    loadToken,
    saveToken,
    refreshToken,
    addListener,
  }

  const provider = hasDefaultUser
    ? createAutoAuthProvider(emitter, utils)
    : createCredentialAuthProvider(emitter, utils)

  return {
    ...provider,
    init,
    onLogin: (handler) => addListener('login', handler),
    onLogout: (handler) => addListener('logout', handler),
    onTokenUpdate: (handler) => addListener('token:update', handler),
  }
}

export default createAuthProvider
