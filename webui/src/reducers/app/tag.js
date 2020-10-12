import {
  TAG_OPENED,
  TAG_CLOSED,
  TAG_SELECTED,
  TAG_UNSELECTED,
  TAG_LOADING,
  TAG_LOADED,
} from '../../actions'
import { 
  UNREGISTER_RESOURCE, 
  GET_LIST,
  GET_MANY,
  GET_MANY_REFERENCE,
  UPDATE,
  CREATE,
  GET_ONE,
} from 'react-admin'

const appendId = (arr, id) => {
  if (arr.includes(id)) {
    return arr
  }
  return [ ...arr, id ]
}
const removeId = (arr, id) => arr.filter(id2 => id2 !== id)

const mergeChildren = (previousChildren, parentId, data) => {
  if (!data) {
    const { [parentId]: removed, ...rest } = previousChildren
    return rest
  }
  return {
    ...previousChildren,
    [parentId]: data.map(record => record.id),
  }
}

const appendOpened = (opened, id) => {
  const newOpened = appendId(opened, id)
  saveOpened(newOpened)
  return newOpened
}

const removeOpened = (opened, id) => {
  const newOpened = removeId(opened, id)
  saveOpened(newOpened)
  return newOpened
}

const saveOpened = (opened) => {
  localStorage.setItem('tag_opened', JSON.stringify(opened))
}

const loadOpened = () => {
  try {
    return JSON.parse(localStorage.getItem('tag_opened')) || []
  } catch (e) {
    return []
  }
}

const initialState = {
  children: {},
  selected: [],
  opened: loadOpened(),
  loading: [],
  data: {},
}

const tagReducer = (previousState = initialState, { type, payload, meta }) => {
  if (meta && meta.resource === 'tag') {
    if ([GET_LIST, GET_MANY, GET_MANY_REFERENCE].includes(meta.fetchResponse)) {
      const newData = {
        ...previousState.data,
        ...Object.values(payload.data).reduce((acc, record) => {
          acc[record.id] = record
          return acc
        }, {})
      }
      return {
        ...previousState,
        data: newData,
      }
    }
    if ([UPDATE, CREATE, GET_ONE].includes(meta.fetchResponse)) {
      const newData = {
        ...previousState.data,
        [payload.data.id]: payload.data,
      }
      return {
        ...previousState,
        data: newData,
      }
    }
  }

  switch (type) {
    case TAG_LOADING:
      return { ...previousState, loading: appendId(previousState.loading, payload.parentId) }
    case TAG_LOADED:
      return {
        ...previousState,
        loading: removeId(previousState.loading, payload.parentId),
        children: mergeChildren(previousState.children, payload.parentId, payload.data),
      }
    case TAG_OPENED:
      return { ...previousState, opened: appendOpened(previousState.opened, payload.id) }
    case TAG_CLOSED:
      return { ...previousState, opened: removeOpened(previousState.opened, payload.id) }
    case TAG_SELECTED:
      return { ...previousState, selected: appendId(previousState.selected, payload.id) }
    case TAG_UNSELECTED:
      return { ...previousState, selected: removeId(previousState.selected, payload.id) }
    case UNREGISTER_RESOURCE:
      if (payload === 'tag') {
        return initialState
      }
      return previousState
    default:
      return previousState
  }
}

export default tagReducer
