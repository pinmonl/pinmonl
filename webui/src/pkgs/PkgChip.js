import React, {
  useMemo,
} from 'react'
import { Chip } from '@material-ui/core'
import MonlerIcon from '../monlers/MonlerIcon'
import { prettyTime } from '@/utils/pretty'

const PkgChip = ({ pkg, ...props }) => {
  const provider = useMemo(() => pkg.provider, [pkg])
  const stats = useMemo(() => pkg.stats || [], [pkg])

  const label = useMemo(() => {
    switch (provider) {
      case 'youtube': {
        const filtered = filterStats(stats, 'video')
        if (filtered.length > 0) {
          return prettyTime(filtered[0].recordedAt)
        }
        return
      }
      case 'npm':
      case 'docker':
      case 'git':
      case 'github': {
        const filtered = filterStats(stats, 'tag')
        if (filtered.length > 0) {
          return filtered[0].value
        }
        return 'latest'
      }
      default:
    }
  }, [provider, stats])

  return (
    <Chip
      {...props}
      icon={<MonlerIcon name={provider} size="small" />}
      variant="outlined"
      label={label}
      size="small"
    />
  )
}

const filterStats = (stats, kind, latest) => 
  stats.filter(stat => stat.kind === kind && (typeof latest === 'boolean' ? stat.isLatest === latest : true))

export default PkgChip
