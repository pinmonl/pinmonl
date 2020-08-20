export const base64ToBlob = (src, mime) => {
  const binary = atob(src)
  const arr = []
  for (let i = 0; i < binary.length; i++) {
    arr.push(binary.charCodeAt(i))
  }

  return new Blob([new Uint8Array(arr)], { type: mime })
}

export const base64ToFile = (src, name, mime) => {
  const binary = atob(src)
  const arr = []
  for (let i = 0; i < binary.length; i++) {
    arr.push(binary.charCodeAt(i))
  }

  return new File([new Uint8Array(arr)], name, { type: mime })
}
