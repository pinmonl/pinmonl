import React, { cloneElement } from 'react'
import {
  AppBar as MuiAppBar,
  Toolbar,
  Tooltip,
  IconButton,
  Typography,
  LinearProgress,
  makeStyles,
} from '@material-ui/core'
import { UserMenu, toggleSidebar, useTranslate } from 'react-admin'
import { useSelector, useDispatch } from 'react-redux'
import MenuIcon from '@material-ui/icons/Menu'
import clsx from 'clsx'

const useStyles = makeStyles(theme => ({
  menuButton: {
    marginLeft: theme.spacing(-2),
  },
  spacer: {
    flex: '1 1 0',
  },
}))

const AppBar = ({
  logout,
  userMenu,
  open,
}) => {
  const classes = useStyles()
  const translate = useTranslate()
  const dispatch = useDispatch()
  const loading = useSelector(state => state.admin.loading > 0)

  return (
    <MuiAppBar color="primary">
      <Toolbar
        variant="dense"
      >
        <Tooltip title={translate(open ? 'ra.action.close_menu' : 'ra.action.open_menu')}>
          <IconButton
            color="inherit"
            onClick={() => dispatch(toggleSidebar())}
            className={clsx(classes.menuButton)}
          >
            <MenuIcon />
          </IconButton>
        </Tooltip>
        <Typography
          variant="h6"
          color="inherit"
          id="react-admin-title"
        />
        <div className={classes.spacer} />
        {cloneElement(userMenu, { logout })}
      </Toolbar>
      {loading && <LinearProgress />}
    </MuiAppBar>
  )
}

AppBar.defaultProps = {
  userMenu: <UserMenu />,
}

export default AppBar
