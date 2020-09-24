import { eventChannel } from 'redux-saga'

const createAuthChannel = (authProvider) => {
  return eventChannel(emit => {
    const offTokenUpdate = authProvider.onTokenUpdate(({ token, expireAt }) => {
      emit({ token, expireAt })
    })

    return () => {
      offTokenUpdate()
    }
  })
}

export default createAuthChannel
