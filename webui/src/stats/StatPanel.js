import React, {
  useMemo,
} from 'react'
import GithubStatPanel from './GithubStatPanel'
import DockerStatPanel from './DockerStatPanel'
import NpmStatPanel from './NpmStatPanel'
import YoutubeStatPanel from './YoutubeStatPanel'
import GitStatPanel from './GitStatPanel'

const StatPanel = (props) => {
  const { pkg } = props

  const LazyPanel = useMemo(() => {
    switch (pkg.provider) {
      case 'github':
        return GithubStatPanel
      case 'docker':
        return DockerStatPanel
      case 'npm':
        return NpmStatPanel
      case 'youtube':
        return YoutubeStatPanel
      case 'git':
        return GitStatPanel
      default:
        return 'div'
    }
  }, [pkg])

  return (
    <LazyPanel {...props} />
  )
}

export default StatPanel
