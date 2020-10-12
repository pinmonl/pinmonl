import { useSelector } from 'react-redux'
import get from 'lodash/get'
import { useQueryWithStore } from 'react-admin'
import { GET_TAG_LIST } from '../actions'

const useGetTagList = ({ parentId }) => {
  const query = {
    type: 'getList',
    resource: 'tag',
    payload: {
      pagination: { page: 1, perPage: 0 },
      sort: { field: 'id', order: 'ASC' },
      filter: { parentId },
    },
  }

  const {
    data: ids,
    total,
    loading,
    loaded,
  } = useQueryWithStore(
    query,
    { action: GET_TAG_LIST },
    state => get(state, ['app', 'tag', 'children', parentId], [])
  )

  const data = useSelector(state =>
    ids.reduce((acc, id) => {
      const record = get(state, ['app', 'tag', 'data', id])
      if (!record) return acc
      return { ...acc, [id]: record }
    }, {})
  )

  return {
    ids,
    data,
    total,
    loading,
    loaded,
  }
}

export default useGetTagList
