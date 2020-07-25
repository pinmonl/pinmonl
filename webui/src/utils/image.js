import { baseURL } from './constants'

export function getURL (imageId) {
  return `${baseURL}/image/${imageId}`
}
