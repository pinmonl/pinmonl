import { useMemo } from 'react'
import { useAuthState } from 'react-admin'
import { TOKEN } from './provider'

const useAuthToken = () => {
  const { loading, authenticated } = useAuthState()

  const token = useMemo(() => {
    if (!authenticated || loading) {
      return ''
    } else {
      return localStorage.getItem(TOKEN)
    }
  }, [authenticated, loading])

  return { token }
}

export default useAuthToken
