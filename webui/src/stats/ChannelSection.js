import React from 'react'
import {
  Paper,
  Box,
} from '@material-ui/core'

const ChannelSection = ({ children, title }) => {
  return (
    <Box my={2} width={{xs: 1, md: '600px'}}>
      <Paper>
        <Box p={3}>
          <Box mb={3} fontWeight={600}>{title}</Box>
          <Box>{children}</Box>
        </Box>
      </Paper>
    </Box>
  )
}

export default ChannelSection
