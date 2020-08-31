import { useContext } from 'react'
import AuthContext from './AuthContext'

const useAuth = () => {
  const ctx = useContext(AuthContext)
  return ctx
}

export default useAuth
