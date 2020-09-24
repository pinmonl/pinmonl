import React from 'react'
import { useSelector } from 'react-redux'
import { useMediaQuery, makeStyles } from '@material-ui/core'
import { MenuItemLink, getResources } from 'react-admin'
import clsx from 'clsx'
import TagMenu from './TagMenu'

const useStyles = makeStyles(theme => ({
  root: {
    display: 'flex',
    flexDirection: 'column',
    flex: '1 1 auto',
    padding: theme.spacing(1, 0),
    overflow: 'auto',
  },
  tags: {
    flex: '1 1 auto',
  },
  menu: {
    flex: '0 0 auto',
    fontSize: 14,
  },
  divider: {
    padding: theme.spacing(0, 1),
  },
}), { name: 'Menu' })

const Menu = ({
  onMenuClick,
  logout,
  dense,
  hasDashboard,
  className,
  ...props
}) => {
  const isXs = useMediaQuery(theme => theme.breakpoints.down('xs'))
  const open = useSelector(state => state.admin.ui.sidebarOpen)
  const resources = useSelector(getResources)
  const classes = useStyles()

  return (
    <div className={clsx(classes.root, className)} {...props}>
      <div className={classes.tags}>
        <TagMenu />
      </div>
      <div className={classes.menu}>
        {resources
          .filter(r => r.hasList)
          .map(resource => (
            <MenuItemLink
              key={resource.name}
              to={`/${resource.name}`}
              primaryText={resource.options && (resource.options.label || resource.name)}
              leftIcon={resource.icon ? <resource.icon /> : null}
              onClick={onMenuClick}
              sidebarIsOpen={open}
              dense={dense}
            />
          ))}
        {isXs && logout}
      </div>
    </div>
  )
}

export default Menu
