import React from 'react'
import { Box } from '@material-ui/core'

const ChannelStat = ({ channel }) => {
  const { substats } = channel
  const aliases = substats.filter((s) => s.kind === 'alias')

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
      {aliases.length > 0 && <Box style={{textOverflow: 'ellipsis', whiteSpace: 'nowrap', overflow: 'hidden'}}>
        {aliases[0].value}
      </Box>}
    </Box>
  )
}

export default ChannelStat
