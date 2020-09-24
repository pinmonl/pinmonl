const loadToken = () => {
  try {
    return {
      token: JSON.parse(localStorage.getItem('token')),
      expireAt: JSON.parse(localStorage.getItem('expire_at')),
    }
  } catch (e) {
    return { token: null, expireAt: null }
  }
}

export default loadToken
