import { base64ToFile } from '../utils/file'
import doRequest from './doRequest'

const isUndefined = (v) => typeof v === 'undefined'

const createDataProvider = (baseURL) => {
  const provider = {
    getList: async (resource, params = {}) => {
      const { page, perPage } = params.pagination || {}
      const { field, order } = params.sort || {}

      const query = new URLSearchParams()

      if (!isUndefined(page)) query.set('page', page)
      if (!isUndefined(perPage)) query.set('pageSize', perPage)

      if (!isUndefined(field)) query.set('sort', field)
      if (!isUndefined(order)) query.set('order', order)

      for (const [key, value] of Object.entries(params.filter || {})) {
        if (!isUndefined(value)) query.set(key, value)
      }

      let dest = `${baseURL}/api/${resource}`
      let queryString = query.toString()
      if (queryString) {
        dest += "?" + queryString
      }

      const { json } = await doRequest(dest)
      return {
        data: json.data,
        total: json.totalCount,
      }
    },
    getOne: async (resource, { id }) => {
      const dest = `${baseURL}/api/${resource}/${id}`
      const { json } = await doRequest(dest)
      return { data: json.data }
    },
    getMany: async (resource, params) => {
      const { ids } = params || {}
      return await provider.getList(resource, {
        pagination: { page: 1, perPage: 0 },
        filter: { id: ids.join(',') },
      })
    },
    getManyReference: async (resource, params) => {
      return { data: [], total: 0 }
    },
    create: async (resource, { data }) => {
      const dest = `${baseURL}/api/${resource}`
      const { json } = await doRequest(dest, {
        method: 'POST',
        body: JSON.stringify(data),
      })
      return { data: json.data }
    },
    update: async (resource, { id, data }) => {
      const dest = `${baseURL}/api/${resource}/${id}`
      const { json } = await doRequest(dest, {
        method: 'PUT',
        body: JSON.stringify(data),
      })
      return { data: json.data }
    },
    updateMany: async (resource, params) => {
      return { data: [] }
    },
    delete: async (resource, { id, previousData }) => {
      const dest = `${baseURL}/api/${resource}/${id}`
      const { json } = await doRequest(dest, {
        method: 'DELETE',
      })
      if (!json) {
        return { data: previousData }
      }
      return { data: json.data }
    },
    deleteMany: async (resource, params) => {
      return { data: [] }
    },
    getCard: async (resource, params) => {
      const { url } = params || {}
      const query = new URLSearchParams()
      query.set('url', url)

      const dest = `${baseURL}/api/card?${query.toString()}`
      const { json } = await doRequest(dest)

      const data = {
        title: json.data.title,
        description: json.data.description,
        image: null,
      }
      if (json.data.imageData) {
        const file = base64ToFile(json.data.imageData, 'card_cover.png', 'image/png')
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
      const { json } = await doRequest(dest, {
        method: 'POST',
        body: data,
      })
      return { data: json.data }
    }
  }
  return provider
}

export default createDataProvider
