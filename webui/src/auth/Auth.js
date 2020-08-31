import React, {
  useState,
  useEffect,
  useMemo,
  useCallback,
} from 'react'
import AuthContext from './AuthContext'
import { baseURL } from '@/utils/constants'
import { doRequest } from '@/dataProvider'
import { useLocation } from 'react-router-dom'
import dayjs from 'dayjs'

const TOKEN = 'token'
const EXPIRE_AT = 'expire_at'

const Auth = ({ children }) => {
  const [token, setToken] = useState(localStorage.getItem(TOKEN))
  const [expireAt, setExpireAt] = useState(localStorage.getItem(EXPIRE_AT))
  const location = useLocation()

  const save = useCallback(({ token, expireAt }) => {
    localStorage.setItem(TOKEN, token)
    localStorage.setItem(EXPIRE_AT, expireAt)
    setToken(token)
    setExpireAt(expireAt)
  }, [setToken, setExpireAt])

  const shouldRefresh = useMemo(() => {
    if (!token || !expireAt) {
      return true
    }
    const expireDate = dayjs(expireAt)
    return expireDate.isBefore(dayjs().add(1, 'day'))
  }, [expireAt, token])

  useEffect(() => {
    let cancelled = false
    const refreshToken = async () => {
      try {
        const resp = await doRequest(`${baseURL}/api/refresh`, {
          method: 'POST',
        })
        if (cancelled) return
        save(resp)
      } catch (e) {
        //
      }
    }

    if (shouldRefresh) {
      refreshToken()
    }
    return () => { cancelled = true }
  }, [save, location, shouldRefresh])

  return (
    <AuthContext.Provider value={{ token }}>
      {children}
    </AuthContext.Provider>
  )
}

export default Auth
