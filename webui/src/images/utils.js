import { baseURL } from '@/utils/constants'
import get from 'lodash/get'

export const getImageUrl = (record, source) => {
  const imageId = get(record, source)
  if (!imageId) return
  return `${baseURL}/image/${imageId}`
}
