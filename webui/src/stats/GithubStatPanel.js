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
  mdiStar,
  mdiEye,
  mdiAlertCircleOutline,
  mdiSourceFork,
  mdiLicense,
  mdiLink,
  mdiCodeBracesBox,
} from '@mdi/js'
import useFetchStats from './useFetchStats'
import { getStats } from './utils'
import StatList from './StatList'
import Stat from './Stat'
import ChannelStat from './ChannelStat'
import ChannelSection from './ChannelSection'
import { prettyNumber } from '@/utils/pretty'

const GithubStatPanel = ({ pkg }) => {
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

  const starCount = useMemo(() => findStat('star_count'), [findStat])
  const openIssueCount = useMemo(() => findStat('open_issue_count'), [findStat])
  const watcherCount = useMemo(() => findStat('watcher_count'), [findStat])
  const forkCount = useMemo(() => findStat('fork_count'), [findStat])
  const license = useMemo(() => findStat('license'), [findStat])
  const lang = useMemo(() => findStat('lang'), [findStat])

  return (
    <React.Fragment>
      <Paper>
        <Box p={3}>
          <StatList>
            <Stat
              value={"GitHub"}
              iconPath={mdiLink}
              href={`https://github.com/${pkg.providerUri}`}
            />
            <Stat
              value={prettyNumber(Number(starCount.value))}
              iconPath={mdiStar}
              suffix={"Stars"}
            />
            <Stat
              value={openIssueCount.value}
              iconPath={mdiAlertCircleOutline}
              suffix={"Open Issues"}
            />
            <Stat
              value={prettyNumber(Number(watcherCount.value))}
              iconPath={mdiEye}
              suffix={"Watchers"}
            />
            <Stat
              value={prettyNumber(Number(forkCount.value))}
              iconPath={mdiSourceFork}
              suffix={"Forks"}
            />
            <Stat
              value={license.value}
              iconPath={mdiLicense}
              suffix={"License"}
            />
            <Stat
              value={lang.value}
              iconPath={mdiCodeBracesBox}
              prefix={"Written in"}
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

export default GithubStatPanel
