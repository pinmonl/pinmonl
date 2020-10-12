import { baseURL } from '../utils/constants'
import get from 'lodash/get'

const getImageUrl = ({ record, source, value }) => {
  const imageId = get(record, source, value)
  if (!imageId) return
  return `${baseURL}/image/${imageId}`
}

export default getImageUrl
