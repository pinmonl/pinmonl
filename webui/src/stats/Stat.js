import React from 'react'
import {
  Box,
  Link,
} from '@material-ui/core'
import Icon from '@mdi/react'

const Stat = ({ value, prefix, suffix, iconPath, href }) => {
  const inner = (
    <Box display="flex" alignItems="center" py={0.8} px={1.2} mx={-0.25}>
      <Box px={0.25} display="inline-flex">
        <Icon size={0.65} path={iconPath} />
      </Box>
      {!!prefix && <Box px={0.25}>{prefix}</Box>}
      <Box px={0.25}>{value}</Box>
      {!!suffix && <Box px={0.25}>{suffix}</Box>}
    </Box>
  )

  if (href) {
    return (
      <Link href={href} color="textPrimary">{inner}</Link>
    )
  }

  return inner
}

export default Stat
