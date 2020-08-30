import React, {
  useMemo,
} from 'react'
import GithubStatPanel from './GithubStatPanel'
import DockerStatPanel from './DockerStatPanel'
import NpmStatPanel from './NpmStatPanel'
import YoutubeStatPanel from './YoutubeStatPanel'

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
      default:
        return 'div'
    }
  }, [pkg])

  return (
    <LazyPanel {...props} />
  )
}

export default StatPanel
