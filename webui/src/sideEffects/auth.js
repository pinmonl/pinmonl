import { createAuthChannel } from '../auth'
import { takeEvery, put } from 'redux-saga/effects'
import { TOKEN_UPDATED } from '../actions'

const auth = (authProvider) => {
  const channel = createAuthChannel(authProvider)

  return function* () {
    yield takeEvery(channel, function* ({ token }) {
      yield put({ type: TOKEN_UPDATED, payload: token })
    })
  }
}

export default auth
