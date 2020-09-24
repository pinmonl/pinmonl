import React from 'react'
import { makeStyles, ThemeProvider } from '@material-ui/core'
import { Notification } from 'react-admin'
import clsx from 'clsx'
import AppBar from './AppBar'
import Menu from './Menu'
import Sidebar from './Sidebar'
import { useSelector } from 'react-redux'

const useStyles = makeStyles((theme) => ({
  root: {
    display: 'flex',
    flexDirection: 'column',
    zIndex: 1,
    minHeight: '100vh',
    minWidth: 0,
    backgroundColor: theme.palette.background.default,
    position: 'relative',
    width: '100%',
    color: theme.palette.getContrastText(
      theme.palette.background.default
    ),
  },
  appFrame: {
    display: 'flex',
    flexDirection: 'column',
    flexGrow: 1,
    [theme.breakpoints.up('xs')]: {
      marginTop: theme.spacing(6),
    },
    [theme.breakpoints.down('xs')]: {
      marginTop: theme.spacing(6),
    },
  },
  contentWithSidebar: {
    display: 'flex',
    flexGrow: 1,
  },
  content: {
    display: 'flex',
    flexDirection: 'column',
    flexGrow: 1,
    flexBasis: 0,
    minWidth: 0,
    padding: theme.spacing(1),
    transition: theme.transitions.create(['margin'], {
      easing: theme.transitions.easing.easeOut,
      duration: theme.transitions.duration.enteringScreen,
    }),
    [theme.breakpoints.down('sm')]: {
      padding: 0,
    },
    [theme.breakpoints.up('sm')]: {
      marginLeft: theme.sidebar.width,
    },
  },
  contentFullWidth: {
    transition: theme.transitions.create(['margin'], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    }),
    [theme.breakpoints.up('sm')]: {
      marginLeft: 0,
    },
  },
}), { name: 'Layout' })

const Layout = ({
  title,
  open,
  notification,
  logout,
  dashboard,
  children,
}) => {
  const classes = useStyles()
  const sidebarOpen = useSelector(state => state.admin.ui.sidebarOpen)

  return (
    <React.Fragment>
      <div className={classes.root}>
        <div className={classes.appFrame}>
          <AppBar
            title={title}
            open={sidebarOpen}
            logout={logout}
          />
          <main className={classes.contentWithSidebar}>
            <Sidebar>
              <Menu logout={logout} hasDashboard={!!dashboard} />
            </Sidebar>
            <div
              className={clsx(classes.content, {
                [classes.contentFullWidth]: !sidebarOpen,
              })}
            >
              {children}
            </div>
          </main>
        </div>
      </div>
      <Notification />
    </React.Fragment>
  )
}

const LayoutWithTheme = ({ theme, ...props }) => {
  return (
    <ThemeProvider theme={theme}>
      <Layout {...props} />
    </ThemeProvider>
  )
}

export default LayoutWithTheme
