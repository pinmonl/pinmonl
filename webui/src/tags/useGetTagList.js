import { useSelector } from 'react-redux'
import get from 'lodash/get'

const useGetTagList = ({ parentId }) => {
  const loading = useSelector(state => state.app.tag.loading.includes(parentId))
  const ids = useSelector(state => state.app.tag.children[parentId] || [])
  const data = useSelector(state =>
    ids.reduce((acc, id) => {
      const record = get(state, ['admin', 'resources', 'tag', 'data', id])
      if (!record) return acc
      return { ...acc, [id]: record }
    }, {})
  )

  return {
    ids,
    data,
    total: ids.length,
    loading,
    loaded: !loading,
  }
}

export default useGetTagList
