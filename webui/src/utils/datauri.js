export function dataURIToBlob (dataURI) {
  const mime = dataURI.split(',')[0].split(':')[1].split(';')[0]
  const b64str = dataURI.split(',')[1]
  return base64ToBlob(b64str, mime)
}

export function base64ToBlob (src, mime) {
  const binary = atob(src)
  const arr = []
  for (let i = 0; i < binary.length; i++) {
    arr.push(binary.charCodeAt(i))
  }

  return new Blob([new Uint8Array(arr)], { type: mime })
}
