import React, {
  useCallback,
  useMemo,
  useState,
  useEffect,
} from 'react'
import { 
  Box, 
  Paper,
} from '@material-ui/core'
import StatList from './StatList'
import Stat from './Stat'
import {
  mdiLink,
  mdiDownload,
  mdiStar,
} from '@mdi/js'
import useFetchStats from './useFetchStats'
import { getStats } from './utils'
import ChannelStat from './ChannelStat'
import ChannelSection from './ChannelSection'
import { prettyNumber } from '@/utils/pretty'

const DockerStatPanel = ({ pkg }) => {
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

  const pullCount = useMemo(() => findStat('pull_count'), [findStat])
  const starCount = useMemo(() => findStat('star_count'), [findStat])
  const link = useMemo(() => {
    const path = pkg.providerUri.startsWith('library')
      ? '_/' + pkg.providerUri.replace('library/', '')
      : 'r/' + pkg.providerUri
    return `https://hub.docker.com/${path}`
  }, [pkg])

  return (
    <React.Fragment>
      <Paper>
        <Box p={3}>
          <StatList>
            <Stat
              value={"Docker"}
              iconPath={mdiLink}
              href={link}
            />
            <Stat
              value={prettyNumber(Number(starCount.value))}
              iconPath={mdiStar}
              suffix={"Stars"}
            />
            <Stat
              value={prettyNumber(Number(pullCount.value))}
              iconPath={mdiDownload}
              suffix={"Pulls"}
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

export default DockerStatPanel
