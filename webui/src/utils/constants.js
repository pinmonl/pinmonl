const baseURL = window.Pinmonl.baseURL
const basePrefix = window.Pinmonl.basePrefix
const isProd = process.env.NODE_ENV == 'production'

const wsBaseURL = (function() {
  const u = new URL(baseURL)
  u.protocol = (u.protocol == 'http:') ? 'ws:' : 'wss:'
  u.pathname = '/ws'
  return u.toString()
}())

const Topics = {
  PINL_UPDATED: 'pinl_updated',
  PINL_DELETED: 'pinl_deleted',
}

export {
  baseURL,
  basePrefix,
  wsBaseURL,
  isProd,
  Topics,
}
