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
import { getStats } from './utils'
import StatList from './StatList'
import Stat from './Stat'
import { prettyNumber, timeFormat } from '@/utils/pretty'

const YoutubeStatPanel = ({ pkg }) => {
  const [latestStats, setLatestStats] = useState(pkg.stats)
  const [videos, setVideos] = useState([])
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

  const findStat = useCallback((kind) => {
    const filtered = getStats(latestStats, { kind })
    return filtered.length > 0 ? filtered[0] : null
  }, [latestStats])

  const subscriberCount = useMemo(() => findStat('subscriber_count'), [findStat])
  const viewCount = useMemo(() => findStat('view_count'), [findStat])
  const videoCount = useMemo(() => findStat('video_count'), [findStat])

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
              value={prettyNumber(Number(subscriberCount.value))}
              iconPath={mdiAccount}
              suffix={"Subscribers"}
            />
            <Stat
              value={prettyNumber(Number(viewCount.value))}
              iconPath={mdiTelevisionPlay}
              suffix={"Views"}
            />
            <Stat
              value={prettyNumber(Number(videoCount.value))}
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
