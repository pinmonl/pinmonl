import auth from './auth'
import tag from './tag'

const createSagas = (authProvider) => ([
  auth(authProvider),
  tag(),
])

export default createSagas
