import { HttpError } from 'react-admin'

const doRequest = async (...args) => {
  try {
    const response = await fetch(...args)
    if (response.status === 204) {
      return { response }
    }
    const json = await response.json()
    if (response.status >= 400) {
      throw new HttpError(json.error)
    }
    return { json, response }
  } catch (e) {
    throw e
  }
}

export default doRequest
