import { baseURL } from '@/utils/constants'
import { doRequest } from '@/dataProvider'
import dayjs from 'dayjs'

export const TOKEN = 'token'
export const EXPIRE_AT = 'expire_at'

const getToken = () => localStorage.getItem(TOKEN)
const getExpireAt = () => localStorage.getItem(EXPIRE_AT)

const refreshToken = async () => {
  return await doRequest(`${baseURL}/api/refresh`, {
    method: 'POST',
  })
}

export const authProvider = {
  login: async () => {
    const resp = await refreshToken()
    localStorage.setItem(TOKEN, resp.token)
    localStorage.setItem(EXPIRE_AT, resp.expireAt)
    return
  },
  logout: async () => {
    localStorage.removeItem(TOKEN)
    localStorage.removeItem(EXPIRE_AT)
    return
  },
  checkError: async (error) => {
    console.log(error)
  },
  checkAuth: async () => {
    const expireAt = getExpireAt()
    if (expireAt) {
      const expireDate = dayjs(expireAt)
      if (expireDate.isBefore(dayjs())) {
        return Promise.reject()
      }
    }

    const token = getToken()
    return token ? Promise.resolve() : Promise.reject()
  },
  getPermissions: async () => {
    return Promise.resolve()
  },
}
