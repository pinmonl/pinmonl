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
import { findStat } from './utils'
import StatList from './StatList'
import Stat from './Stat'
import ChannelStat from './ChannelStat'
import ChannelSection from './ChannelSection'
import { prettyNumber } from '@/utils/pretty'
import useNumberFormat from './useNumberFormat'

const GithubStatPanel = ({ pkg }) => {
  const [latestStats, setLatestStats] = useState(pkg.stats)
  const [channels, setChannels] = useState([])
  const fetchStats = useFetchStats({ pkg })
  const numberFormat = useNumberFormat()

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

  const starCount = useMemo(() => findStat(latestStats, { kind: 'star_count' }), [latestStats])
  const openIssueCount = useMemo(() => findStat(latestStats, { kind: 'open_issue_count' }), [latestStats])
  const watcherCount = useMemo(() => findStat(latestStats, { kind: 'watcher_count' }), [latestStats])
  const forkCount = useMemo(() => findStat(latestStats, { kind: 'fork_count' }), [latestStats])
  const license = useMemo(() => findStat(latestStats, { kind: 'license' }), [latestStats])
  const lang = useMemo(() => findStat(latestStats, { kind: 'lang' }), [latestStats])

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
              stat={starCount}
              format={numberFormat}
              iconPath={mdiStar}
              suffix={"Stars"}
            />
            <Stat
              stat={openIssueCount}
              iconPath={mdiAlertCircleOutline}
              suffix={"Open Issues"}
            />
            <Stat
              stat={watcherCount}
              value={numberFormat}
              iconPath={mdiEye}
              suffix={"Watchers"}
            />
            <Stat
              stat={forkCount}
              format={numberFormat}
              iconPath={mdiSourceFork}
              suffix={"Forks"}
            />
            <Stat
              stat={license}
              iconPath={mdiLicense}
              suffix={"License"}
            />
            <Stat
              stat={lang}
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
