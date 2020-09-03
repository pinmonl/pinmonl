import { prettySize } from '@/utils/pretty'

const useSizeFormat = () => {
  return (value) => prettySize(Number(value))
}

export default useSizeFormat
