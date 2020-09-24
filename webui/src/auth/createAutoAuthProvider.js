import dayjs from '../utils/day'

const createAutoAuthProvider = (emitter, { refreshToken, loadToken, saveToken }) => {
  const provider = {
    login: async (params) => {
      const data = await refreshToken()
      saveToken(data)
    },
    logout: async (params) => {
      saveToken()
    },
    checkAuth: async (params) => {
      const { expireAt } = loadToken()
      const expireDate = dayjs(expireAt)
      if (expireDate.isValid() && expireDate.isAfter(dayjs().add(1, 'day'))) {
        return
      }
      const data = await refreshToken()
      saveToken(data)
    },
    checkError: async (params) => {
    },
    getPermissions: async (params) => {
    },
  }
  return provider
}

export default createAutoAuthProvider
