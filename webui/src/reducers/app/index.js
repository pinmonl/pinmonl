import auth from './auth'
import tag from './tag'
import { CLEAR_STATE } from 'react-admin'

const initialState = {}

const reducer = (previousState = initialState, action) => {
  let state = previousState
  if (action.type === CLEAR_STATE) {
    state = initialState
  }

  return {
    auth: auth(state.auth, action),
    tag: tag(state.tag, action),
  }
}

export default reducer
