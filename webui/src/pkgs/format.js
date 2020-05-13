export function numberFormat (n) {
  return n.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",")
}

export function humanFileSize (bytes, si) {
  const thresh = si ? 1000 : 1024
  if(Math.abs(bytes) < thresh) {
    return bytes + ' B'
  }
  const units = si
    ? ['kB','MB','GB','TB','PB','EB','ZB','YB']
    : ['KiB','MiB','GiB','TiB','PiB','EiB','ZiB','YiB']
  let u = -1
  do {
    bytes /= thresh
    ++u
  } while (Math.abs(bytes) >= thresh && u < units.length - 1)
  return bytes.toFixed(1)+' '+units[u]
}

export function humanizeNumber (value) {
  const thresh = 1000
  if (Math.abs(value) < thresh) {
    return value
  }
  const units = ['k', 'M', 'B']
  let u = -1
  do {
    value /= thresh
    ++u
  } while (Math.abs(value) >= thresh && u < units.length - 1)
  return value.toFixed(0) + units[u]
}