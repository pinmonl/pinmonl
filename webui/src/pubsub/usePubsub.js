import { useContext } from 'react'
import PubsubContext from './PubsubContext'

const usePubsub = () => {
  const pubsub = useContext(PubsubContext)
  return pubsub
}

export default usePubsub
