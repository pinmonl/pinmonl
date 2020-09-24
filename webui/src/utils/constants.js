const Pinmonl = window.Pinmonl || {}

export const baseURL = Pinmonl.baseURL

export const baseWs = function () {
  const url = new URL(baseURL)
  url.protocol = url.protocol.replace(/([^\w])/, '') === 'http' ? 'ws' : 'wss'
  url.pathname = '/ws'
  return url.toString()
}()

export const hasDefaultUser = Pinmonl.hasDefaultUser
