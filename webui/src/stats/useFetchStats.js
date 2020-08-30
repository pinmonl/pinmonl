import { useCallback } from 'react'
import { useDataProvider } from 'react-admin'

const useFetchStats = ({ pkg }) => {
  const dataProvider = useDataProvider()

  const fetch = useCallback(async ({ latest, parent, kind, perPage = 25, page = 1 }) => {
    return await dataProvider.getList('stat', {
      pagination: { perPage, page },
      filter: { latest, parent, kind, pkg: pkg.id },
    })
  }, [dataProvider, pkg])

  const fetchTree = useCallback(async ({ depth = 1, kind, latest, ...params }) => {
    const { data: rootStats } = await fetch({ ...params, latest, kind })
    const parentIds = rootStats.filter(stat => stat.hasChildren).map(stat => stat.id)
    if (parentIds.length < 1 || depth < 1) {
      return { data: rootStats }
    }

    const { data: children } = await fetchTree({
      ...params,
      parent: parentIds,
      depth: depth - 1,
    })
    rootStats.forEach(stat => {
      if (!stat.hasChildren) return
      stat.substats = children.filter((child) => child.parentId === stat.id)
    })
    return { data: rootStats }
  }, [fetch])

  return fetchTree
}

export default useFetchStats
