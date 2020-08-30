import React from 'react'
import { Layout as RaLayout } from 'react-admin'
import { makeStyles } from '@material-ui/core'

const useStyles = makeStyles((theme) => ({
  root: { minWidth: 0 },
  content: { minWidth: 0 },
}), { name: 'Layout' })

const Layout = (props) => {
  const classes = useStyles()
  return (
    <RaLayout classes={classes} {...props} />
  )
}

export default Layout
