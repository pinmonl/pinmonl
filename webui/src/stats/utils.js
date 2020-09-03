export const getStats = (stats, where) => {
  let matches = []
  let { kind, value } = where

  if (!Array.isArray(kind)) {
    kind = [kind]
  }
  matches.push((stat) => kind.includes(stat.kind))

  if (typeof value !== 'undefined') {
    value = [value]
    matches.push((stat) => value.includes(stat.value))
  }

  return stats.filter((stat) => {
    return matches.every((match) => match(stat))
  })
}

export const findStat = (stats, where) => {
  const filtered = getStats(stats, where)
  return filtered[0]
}
