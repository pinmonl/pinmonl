import React from 'react'
import { Chip } from '@material-ui/core'
import MonlerIcon from '../monlers/MonlerIcon'
import { prettyTime } from '@/utils/pretty'

const PkgChip = ({ pkg, ...props }) => {
  const provider = pkg.provider

  let label
  switch (provider) {
    case 'youtube': {
      const stats = filterStats(pkg.stats, 'video')
      if (stats.length > 0) {
        label = prettyTime(stats[0].recordedAt)
      }
      break
    }
    case 'npm':
    case 'docker':
    case 'git':
    case 'github': {
      const stats = filterStats(pkg.stats, 'tag')
      if (stats.length > 0) {
        label = stats[0].value
      }
      break
    }
    default:
  }

  return (
    <Chip
      {...props}
      icon={<MonlerIcon name={pkg.provider} size="small" />}
      variant="outlined"
      label={label}
      size="small"
    />
  )
}

const filterStats = (stats, kind, latest) => 
  stats.filter(stat => stat.kind === kind && (typeof latest === 'boolean' ? stat.isLatest === latest : true))

export default PkgChip
