import { loadToken } from '../../auth'
import { TOKEN_UPDATED } from '../../actions'

const initialState = {
  token: null,
}

const getInitialState = () => {
  const { token } = loadToken()
  initialState.token = token
}

getInitialState()

const authReducer = (previousState = initialState, { type, payload }) => {
  if (type !== TOKEN_UPDATED) {
    return previousState
  }
  return { token: payload || null }
}

export default authReducer
