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
import { findStat } from './utils'
import StatList from './StatList'
import Stat from './Stat'
import ChannelStat from './ChannelStat'
import ChannelSection from './ChannelSection'
import useNumberFormat from './useNumberFormat'
import useSizeFormat from './useSizeFormat'

const NpmStatPanel = ({ pkg }) => {
  const [latestStats, setLatestStats] = useState(pkg.stats)
  const [channels, setChannels] = useState([])
  const fetchStats = useFetchStats({ pkg })
  const numberFormat = useNumberFormat()
  const sizeFormat = useSizeFormat()

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

  const downloadCount = useMemo(() => findStat(latestStats, { kind: 'download_count' }), [latestStats])
  const fileCount = useMemo(() => findStat(latestStats, { kind: 'file_count' }), [latestStats])
  const size = useMemo(() => findStat(latestStats, { kind: 'size' }), [latestStats])

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
              stat={downloadCount}
              format={numberFormat}
              iconPath={mdiCloudDownloadOutline}
              suffix={"Monthly Downloads"}
            />
            <Stat
              stat={fileCount}
              iconPath={mdiFileMultipleOutline}
              suffix={"Files"}
            />
            <Stat
              stat={size}
              format={sizeFormat}
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
