import { prettyNumber } from '@/utils/pretty'

const useNumberFormat = () => {
  return (value) => prettyNumber(Number(value))
}

export default useNumberFormat
