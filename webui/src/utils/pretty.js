import dayjs from 'dayjs'
import prettyBytes from 'pretty-bytes'
import relativeTime from 'dayjs/plugin/relativeTime'

dayjs.extend(relativeTime)

const prettyTime = (target, from, withSuffix) => {
  return dayjs(target).from(dayjs(from), withSuffix)
}

const prettySize = (number, options) => {
  return prettyBytes(number, options)
}

const prettyNumber = (number) => {
  var s = ['', 'k', 'M', 'B', 'T']
  var e = Math.floor(Math.log(number) / Math.log(1000))
  if (e < 1) return number
  return (number / Math.pow(1000, e)).toFixed(2) + s[e]
}

const timeFormat = (target, format = 'DD MMM YYYY, HH:mm') => {
  return dayjs(target).format(format)
}

export {
  prettyTime,
  prettySize,
  prettyNumber,
  timeFormat,
}
