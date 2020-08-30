import React, {
  useEffect,
  useCallback,
  useMemo,
  useState,
} from 'react'
import {
  Box,
  Paper,
} from '@material-ui/core'
import {
  mdiLink,
  mdiCloudDownloadOutline,
  mdiPackageVariant,
  mdiFileMultipleOutline,
} from '@mdi/js'
import useFetchStats from './useFetchStats'
import { getStats } from './utils'
import StatList from './StatList'
import Stat from './Stat'
import ChannelStat from './ChannelStat'
import ChannelSection from './ChannelSection'
import { prettySize, prettyNumber } from '@/utils/pretty'

const NpmStatPanel = ({ pkg }) => {
  const [latestStats, setLatestStats] = useState(pkg.stats)
  const [channels, setChannels] = useState([])
  const fetchStats = useFetchStats({ pkg })

  useEffect(() => {
    let cancelled = false
    const fetchData = async () => {
      try {
        const { data } = await fetchStats({ latest: 1, perPage: 0 })
        if (cancelled) return
        setLatestStats(data)
      } catch (e) {
        //
      }
    }
    fetchData()
    return () => cancelled = true
  }, [fetchStats])

  useEffect(() => {
    let cancelled = false
    const fetchData = async () => {
      try {
        const { data } = await fetchStats({ perPage: 0, kind: 'channel' })
        if (cancelled) return
        setChannels(data)
      } catch (e) {
        //
      }
    }
    fetchData()
    return () => cancelled = true
  }, [])

  const findStat = useCallback((kind) => {
    const filtered = getStats(latestStats, { kind })
    return filtered.length > 0 ? filtered[0] : null
  }, [latestStats])

  const downloadCount = useMemo(() => findStat('download_count'), [findStat])
  const fileCount = useMemo(() => findStat('file_count'), [findStat])
  const size = useMemo(() => findStat('size'), [findStat])

  return (
    <React.Fragment>
      <Paper>
        <Box p={3}>
          <StatList>
            <Stat
              value={"NPM"}
              iconPath={mdiLink}
              href={`https://www.npmjs.com/package/${pkg.providerUri}`}
            />
            <Stat
              value={prettyNumber(Number(downloadCount.value))}
              iconPath={mdiCloudDownloadOutline}
              suffix={"Monthly Downloads"}
            />
            <Stat
              value={fileCount.value}
              iconPath={mdiFileMultipleOutline}
              suffix={"Files"}
            />
            <Stat
              value={prettySize(Number(size.value))}
              iconPath={mdiPackageVariant}
            />
          </StatList>
        </Box>
      </Paper>
      <ChannelSection title="Tags">
        {channels.map((channel) => (
          <ChannelStat key={channel.id} channel={channel} />
        ))}
      </ChannelSection>
    </React.Fragment>
  )
}

export default NpmStatPanel
