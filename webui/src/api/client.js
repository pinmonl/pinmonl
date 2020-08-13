import { baseURL } from '@/utils/constants'
import { encodeQuery } from './utils'

class Client {
  constructor () {
    this.token = ''
    this.errorHandler = () => {}
  }

  setToken (token) {
    this.token = token
  }

  async do (url, opts) {
    let headers = (opts || {}).headers || {}
    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`
    }

    const res = await fetch(`${baseURL}${url}`, { ...opts, headers })
    return res
  }

  async doJson (url, opts) {
    const res = await this.do(url, opts)

    // No Content.
    if (res.status == 204) {
      return res

    // Success status.
    } else if (res.status >= 200 && res.status < 300) {
      return res.json()

    // Errors.
    } else {
      const error = { status: res.status, response: res }

      if (this.errorHandler(error) === false) {
        return
      }
      throw error
    }
  }

  // Login.
  async login (login, password) {
    const data = { login, password }

    return this.doJson('/api/login', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  // Signup.
  async signup (login, password, username) {
    const data = { login, password, username }

    return this.doJson('/api/login', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  // Refresh user token.
  async refresh (token) {
    return this.doJson('/api/refresh', {
      method: 'POST',
    }, token)
  }

  // List pins.
  async listPins (query) {
    return this.doJson(`/api/pinl${encodeQuery(query)}`)
  }

  // Find pin by id.
  async findPin (id) {
    return this.doJson(`/api/pinl/${id}`)
  }

  // Create pin.
  async createPin (data) {
    return this.doJson('/api/pinl', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  // Update pin by id.
  async updatePin (id, data) {
    return this.doJson(`/api/pinl/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  }

  // Delete pin by id.
  async deletePin (id) {
    return this.doJson(`/api/pinl/${id}`, {
      method: 'DELETE',
    })
  }

  // List tags.
  async listTags (query) {
    return this.doJson(`/api/tag${encodeQuery(query)}`)
  }

  // Find tag by id.
  async findTag (id) {
    return this.doJson(`/api/tag/${id}`)
  }

  // Create tag.
  async createTag (data) {
    return this.doJson('/api/tag', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  // Update tag by id.
  async updateTag (id, data) {
    return this.doJson(`/api/tag/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  }

  // Delete tag by id.
  async deleteTag (id) {
    return this.doJson(`/api/tag/${id}`, {
      method: 'DELETE',
    })
  }

  // Crawl card information of url.
  async card(url) {
    return this.doJson(`/api/card${encodeQuery({ url })}`)
  }

  // Upload image to the target.
  async uploadImage (targetPath, file) {
    const data = new FormData()
    data.append('file', file)

    return this.doJson(`/api/${targetPath}/image`, {
      method: 'POST',
      body: data,
    })
  }

  // Upload image to pin.
  async uploadPinImage (pinId, file) {
    return this.uploadImage(`pinl/${pinId}`, file)
  }
}

export default Client
