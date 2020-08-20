import React from 'react'
import { Chip, Box } from '@material-ui/core'
import { absTagName } from './utils'

const TagChip = ({ label, p, ...props }) => {
  return (
    <Box p={p}>
      <Chip
        size="small"
        label={absTagName(label)}
        {...props}
      />
    </Box>
  )
}

TagChip.defaultProps = {
  p: 0.5,
}

export default TagChip
