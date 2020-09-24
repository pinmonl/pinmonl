import React, { cloneElement, Children } from 'react'
import { Drawer, Toolbar, useMediaQuery, makeStyles } from '@material-ui/core'
import { useDispatch, useSelector } from 'react-redux'
import { setSidebarVisibility } from 'react-admin'

const useStyles = makeStyles(theme => ({
  paper: {
    zIndex: theme.zIndex.sidebar,
    width: theme.sidebar.width,
    backgroundColor: 'transparent',
    border: 'none',
    [theme.breakpoints.down('xs')]: {
      backgroundColor: theme.palette.background.default,
    },
  },
}))

const Sidebar = ({
  children,
}) => {
  const classes = useStyles()
  const dispatch = useDispatch()
  const isXs = useMediaQuery(theme => theme.breakpoints.down('xs'))
  const open = useSelector(state => state.admin.ui.sidebarOpen)
  const handleClose = () => dispatch(setSidebarVisibility(false))

  return isXs ? (
    <Drawer
      variant="temporary"
      open={open}
      classes={classes}
      onClose={handleClose}
    >
      {cloneElement(Children.only(children), {
        onMenuClick: handleClose,
      })}
    </Drawer>
  ) : (
    <Drawer
      variant="persistent"
      open={open}
      classes={classes}
      onClose={handleClose}
    >
      <Toolbar variant="dense" />
      {children}
    </Drawer>
  )
}

export default Sidebar
