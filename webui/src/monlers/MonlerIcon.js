import React from 'react'
import { getIcon } from './icons'
import { makeStyles } from '@material-ui/core'
import clsx from 'clsx'

const useStyles = makeStyles((theme) => ({
  root: {
    width: '20px',
    height: '20px',
    padding: '2px',
    fill: 'currentColor',
  },
  sizeSmall: {
    width: '16px',
    height: '16px',
  },
}), { name: 'MonlerIcon' })

const MonlerIcon = ({ name, className, size, ...props }) => {
  const icon = getIcon(name)
  const classes = useStyles()

  return (
    <svg
      {...props}
      className={clsx(className, classes.root, {
        [classes.sizeSmall]: size === 'small',
      })}
      viewBox="0 0 24 24"
      style={{ color: `#${icon.hex}` }}
    >
      <path d={icon.path} />
    </svg>
  )
}

export default MonlerIcon
