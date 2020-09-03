import React from 'react'
import PropTypes from 'prop-types'
import {
  Box,
  Link,
} from '@material-ui/core'
import Icon from '@mdi/react'

const Stat = ({ stat, format, value, prefix, suffix, iconPath, href }) => {
  if (!stat && !value) {
    return null
  }

  let finalValue
  if (stat && stat.value) {
    finalValue = stat.value
  } else {
    finalValue = value
  }
  if (format) {
    finalValue = format(finalValue)
  }

  const inner = (
    <Box display="flex" alignItems="center" py={0.8} px={1.2} mx={-0.25}>
      <Box px={0.25} display="inline-flex">
        <Icon size={0.65} path={iconPath} />
      </Box>
      {!!prefix && <Box px={0.25}>{prefix}</Box>}
      <Box px={0.25}>{finalValue}</Box>
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

Stat.propTypes = {
  stat: PropTypes.object,
  format: PropTypes.func,
}

export default Stat
