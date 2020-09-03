import React, {
  useEffect,
  useCallback,
  useMemo,
  useState,
} from 'react'
import {
  Box,
  Link,
  Paper,
} from '@material-ui/core'
import {
  mdiLink,
  mdiAccount,
  mdiTelevisionPlay,
  mdiVideo,
} from '@mdi/js'
import useFetchStats from './useFetchStats'
import { findStat } from './utils'
import StatList from './StatList'
import Stat from './Stat'
import { timeFormat } from '@/utils/pretty'
import useNumberFormat from './useNumberFormat'

const YoutubeStatPanel = ({ pkg }) => {
  const [latestStats, setLatestStats] = useState(pkg.stats)
  const [videos, setVideos] = useState([])
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
        const { data } = await fetchStats({ perPage: 20, kind: 'video' })
        if (cancelled) return
        setVideos(data)
      } catch (e) {
        //
      }
    }
    fetchData()
    return () => cancelled = true
  }, [])

  const subscriberCount = useMemo(() => findStat(latestStats, { kind: 'subscriber_count' }), [latestStats])
  const viewCount = useMemo(() => findStat(latestStats, { kind: 'view_count' }), [latestStats])
  const videoCount = useMemo(() => findStat(latestStats, { kind: 'video_count' }), [latestStats])

  return (
    <React.Fragment>
      <Paper>
        <Box p={3}>
          <StatList>
            <Stat
              value={"YouTube"}
              iconPath={mdiLink}
              href={`https://youtube.com/channel/${pkg.providerUri}`}
            />
            <Stat
              stat={subscriberCount}
              format={numberFormat}
              iconPath={mdiAccount}
              suffix={"Subscribers"}
            />
            <Stat
              stat={viewCount}
              format={numberFormat}
              iconPath={mdiTelevisionPlay}
              suffix={"Views"}
            />
            <Stat
              stat={videoCount}
              format={numberFormat}
              iconPath={mdiVideo}
              suffix={"Videos"}
            />
          </StatList>
        </Box>
      </Paper>
      <Box my={2} width={1}>
        <Paper>
          <Box p={3}>
            <Box mb={3} fontWeight={600}>Uploads</Box>
            <Box display="column">
              {videos.map((video) => (
                <Box key={video.id} my={1.2} fontSize="14px" display="flex">
                  <Box style={{overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap'}}>
                    {video.name}
                  </Box>
                  <Box mx={1} flexGrow={1} flexShrink={1} borderBottom="1px dotted" />
                  <Box flexShrink={0}>
                    <Link href={`https://www.youtube.com/watch/${video.value}`}>
                      {timeFormat(video.recordedAt)}
                    </Link>
                  </Box>
                </Box>
              ))}
            </Box>
          </Box>
        </Paper>
      </Box>
    </React.Fragment>
  )
}

export default YoutubeStatPanel
