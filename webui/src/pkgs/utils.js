export function formatRepeatParam (param) {
  param = param || ''
  if (param == '') {
    return []
  }
  if (typeof param == 'string') {
    return param.split('/')
  }
  return param
}
