import { takeEvery, put, select } from 'redux-saga/effects'
import {
  GET_TAG_LIST,
  GET_TAG_LIST_LOADING,
  GET_TAG_LIST_SUCCESS,
  GET_TAG_LIST_FAILURE,
  TAG_LOADING,
  TAG_LOADED,
  TAG_OPENED,
} from '../actions'
import { FETCH_END, GET_LIST, REGISTER_RESOURCE } from 'react-admin'

function* handleTagListRequest(action) {
  const { type, payload, requestPayload, meta } = action
  const parentId = meta.fetchStatus === FETCH_END
    ? requestPayload.filter.parentId
    : payload.filter.parentId

  switch (type) {
    case GET_TAG_LIST_LOADING:
      yield put({ type: TAG_LOADING, payload: { parentId } })
      break
    case GET_TAG_LIST_SUCCESS:
      yield put({ type: TAG_LOADED, payload: { parentId, data: payload.data } })
      break
    case GET_TAG_LIST_FAILURE:
      yield put({ type: TAG_LOADED, payload: { parentId, data: null } })
      break
    default:
      return
  }
}

function* fetchTagChildren(action) {
  const { payload: { id: tagId } } = action
  const isFetched = yield select(state => typeof state.app.tag.children[tagId] !== 'undefined')
  if (isFetched) {
    return
  }
  yield put(getTagList(tagId))
}

function* restoreOpenedTags() {
  yield put(getTagList(''))

  const opened = yield select(state => state.app.tag.opened)
  for (const tagId of opened) {
    yield put(getTagList(tagId))
  }
}

const getTagList = (parentId) => ({
  type: GET_TAG_LIST,
  payload: {
    pagination: { page: 1, perPage: 0 },
    sort: { field: 'id', order: 'ASC' },
    filter: { parentId },
  },
  meta: {
    resource: 'tag',
    fetch: GET_LIST,
  },
})

const takeGetTagList = (action) => action.type.startsWith(GET_TAG_LIST)
const takeTagOpened = (action) => action.type === TAG_OPENED
const takeRegisterTagResource = (action) => action.type === REGISTER_RESOURCE && action.payload.name === 'tag'

const tag = () => {
  return function* () {
    yield takeEvery(takeGetTagList, handleTagListRequest)
    yield takeEvery(takeTagOpened, fetchTagChildren)
    yield takeEvery(takeRegisterTagResource, restoreOpenedTags)
  }
}

export default tag
