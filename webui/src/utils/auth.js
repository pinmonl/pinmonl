import store from '@/store'
import router from '@/router'

export function getToken () {
  return localStorage.getItem('token') || ''
}

export function setToken (token) {
  localStorage.setItem('token', token)
}

export function getExpireAt () {
  return localStorage.getItem('expire_at') || null
}

export function setExpireAt (at) {
  localStorage.setItem('expire_at', at)
}

export function restore () {
  store.commit('SET_TOKEN', getToken())
}

export function logout () {
  store.commit('SET_TOKEN', '')
  setToken('')
  setExpireAt(null)
  router.push('/login')
}

export function authed (accessInfo) {
  save(accessInfo)
  router.push('/')
}

export function save ({ token, expireAt }) {
  store.commit('SET_TOKEN', token)
  setToken(token)
  setExpireAt(expireAt)
}
