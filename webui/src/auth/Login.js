import React, {
  useEffect,
} from 'react'
import { useLogin, useAuthState } from 'react-admin'

const Login = (props) => {
  const { authenticated } = useAuthState()
  const login = useLogin()

  useEffect(() => {
    if (authenticated) return
    login()
  }, [authenticated])

  return <div />
}

export default Login
