import React from 'react'
import { Box } from '@material-ui/core'

const StatList = ({ children }) => {
  return (
    <Box
      display="flex"
      flexWrap="wrap"
      mx={-1.2}
      fontSize={14}
    >
      {children}
    </Box>
  )
}

export default StatList
