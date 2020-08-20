import React from 'react'
import { Box, Grid } from '@material-ui/core'
import TagChip from './TagChip'
import get from 'lodash/get'

const TagArrayField = ({ record = {}, source, ...props }) => {
  const tags = get(record, source, [])

  return (
    <Box display="flex" m={-0.5} flexWrap="wrap">
      {tags.map(tag => (
        <Grid item key={tag}>
          <TagChip label={tag} />
        </Grid>
      ))}
    </Box>
  )
}

export default TagArrayField
