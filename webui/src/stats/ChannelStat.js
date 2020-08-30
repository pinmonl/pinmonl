import React from 'react'
import { Box } from '@material-ui/core'

const ChannelStat = ({ channel }) => {
  const { substats } = channel

  return (
    <Box display="flex" my={1} fontSize="14px">
      <Box>{channel.value}</Box>
      <Box
        flexShrink={1}
        flexGrow={1}
        borderBottom={"1px dotted"}
        minWidth={0}
        mx={0.6}
      />
      {!!substats && <Box style={{textOverflow: 'ellipsis', whiteSpace: 'nowrap', overflow: 'hidden'}}>
        {substats[0].value}
      </Box>}
    </Box>
  )
}

export default ChannelStat
