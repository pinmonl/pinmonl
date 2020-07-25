export function encodeQuery (query) {
  const params = new URLSearchParams(query)
  const encoded = params.toString()
  if (encoded) {
    return `?${encoded}`
  }
  return ''
}
