import React from 'react'
import { useSelector } from 'react-redux'
import { useMediaQuery, makeStyles } from '@material-ui/core'
import { MenuItemLink, getResources, useTranslate } from 'react-admin'
import clsx from 'clsx'
import TagMenu from './TagMenu'
import TagIcon from '@material-ui/icons/LocalOffer'
import PinIcon from '@material-ui/icons/Bookmark'
import { useRouteMatch } from 'react-router-dom'

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
  const classes = useStyles()
  const translate = useTranslate()
  const pinMatch = useRouteMatch('/pin')

  return (
    <div className={clsx(classes.root, className)} {...props}>
      <div className={classes.menu}>

        <MenuItemLink
          to="/pin"
          primaryText={translate('resources.pin.name')}
          leftIcon={<PinIcon />}
          onClick={onMenuClick}
          sidebarIsOpen={open}
          dense={dense}
        />
        <MenuItemLink
          to="/tag"
          primaryText={translate('resources.tag.name')}
          leftIcon={<TagIcon />}
          onClick={onMenuClick}
          sidebarIsOpen={open}
          dense={dense}
        />
        <div className={classes.tags}>
          <TagMenu open={!!pinMatch} />
        </div>
        {isXs && logout}
      </div>
    </div>
  )
}

export default Menu
