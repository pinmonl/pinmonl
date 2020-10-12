import React from 'react'
import PropTypes from 'prop-types'
import { makeStyles } from '@material-ui/core'
import clsx from 'clsx'

const useStyles = makeStyles(theme => ({
  anchor: {
    textDecoration: 'none',
    '&:hover': {
      textDecoration: 'underline',
    },
  },
}))

const ExternalLink = ({ className, children, ...props }) => {
  const classes = useStyles(props)
  return (
    <a {...props} className={clsx(classes.anchor, className)}>
      {children}
    </a>
  )
}

ExternalLink.defaultProps = {
  target: '_blank',
}

ExternalLink.propTypes = {
  children: PropTypes.any.isRequired,
}

export default ExternalLink
