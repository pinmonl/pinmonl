import { takeEvery, put, select, call } from 'redux-saga/effects'
import { createMatchSelector } from 'connected-react-router'
import {
  GET_TAG_LIST,
  GET_TAG_LIST_LOADING,
  GET_TAG_LIST_SUCCESS,
  GET_TAG_LIST_FAILURE,
  TAG_LOADING,
  TAG_LOADED,
  TAG_SELECTED,
} from '../actions'
import { FETCH_END } from 'react-admin'

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

function* restoreSelected() {
  const match = yield select(createMatchSelector('/pin/t/:tagIds'))
  if (!match) {
    return
  }

  const tagIds = match.params.tagIds.split(',')
  for (const id of tagIds) {
    yield put({
      type: TAG_SELECTED,
      payload: { id },
    })
  }
}

const takeGetTagList = (action) => action.type.startsWith(GET_TAG_LIST)

const tag = () => {
  return function* () {
    yield call(restoreSelected)
    yield takeEvery(takeGetTagList, handleTagListRequest)
  }
}

export default tag
