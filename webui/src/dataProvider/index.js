import { baseURL } from '@/utils/constants'
import { HttpError } from 'react-admin'
import { base64ToFile } from '@/utils/file'

const doRequest = async (...args) => {
  try {
    const resp = await fetch(...args)
    if (resp.status === 204) {
      return true
    }
    const json = await resp.json()
    if (resp.status >= 400) {
      throw new HttpError(json.error)
    }
    return json
  } catch (e) {
    throw e
  }
}

const provider = {
  getList: async (resource, params = {}) => {
    const { page, perPage } = params.pagination || {}
    const { field, order } = params.sort || {}

    const query = new URLSearchParams()

    if (page) query.set('page', page)
    if (perPage) query.set('page_size', perPage)

    if (field) query.set('sort', field)
    if (order) query.set('order', order)

    for (const [key, value] of Object.entries(params.filter || {})) {
      if (value) query.set(key, value)
    }

    let dest = `${baseURL}/api/${resource}`
    if (query.toString()) {
      dest += "?" + query.toString()
    }

    const resp = await doRequest(dest)
    return {
      data: resp.data,
      total: resp.totalCount,
    }
  },
  getOne: async (resource, { id }) => {
    const dest = `${baseURL}/api/${resource}/${id}`
    const resp = await doRequest(dest)
    return { data: resp }
  },
  getMany: async (resource, params) => {
    return { data: [] }
  },
  getManyReference: async (resource, params) => {
    return { data: [], total: 0 }
  },
  create: async (resource, { data }) => {
    const dest = `${baseURL}/api/${resource}`
    const resp = await doRequest(dest, {
      method: 'POST',
      body: JSON.stringify(data),
    })
    return { data: resp }
  },
  update: async (resource, { id, data }) => {
    const dest = `${baseURL}/api/${resource}/${id}`
    const resp = await doRequest(dest, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
    return { data: resp }
  },
  updateMany: async (resource, params) => {
    return { data: [] }
  },
  delete: async (resource, { id, previousData }) => {
    const dest = `${baseURL}/api/${resource}/${id}`
    const resp = await doRequest(dest, {
      method: 'DELETE',
    })
    if (resp === true) {
      return { data: previousData }
    }
    return { data: resp }
  },
  deleteMany: async (resource, params) => {
    return { data: [] }
  },
  getCard: async (resource, params) => {
    const { url } = params || {}
    const query = new URLSearchParams()
    query.set('url', url)

    const dest = `${baseURL}/api/card?${query.toString()}`
    const resp = await doRequest(dest)

    const data = {
      title: resp.title,
      description: resp.description,
      image: null,
    }
    if (resp.imageData) {
      const file = base64ToFile(resp.imageData, 'card_cover.png', 'image/png')
      data.image = {
        rawFile: file,
        src: URL.createObjectURL(file),
      }
    }

    return { data }
  },
  createImage: async (targetResource, params) => {
    const { id, image } = params
    if (typeof id === 'undefined') {
      throw new Error('id is missing in payload')
    }
    if (!image || !(image.rawFile instanceof File)) {
      return { data: null }
    }

    const data = new FormData()
    data.append('file', image.rawFile)

    const dest = `${baseURL}/api/${targetResource}/${id}/image`
    const resp = await doRequest(dest, {
      method: 'POST',
      body: data,
    })
    return { data: resp }
  }
}

export default provider
