import React, { useMemo, useState } from 'react'
import {
  ListItem,
  ListItemText,
  ListItemAvatar,
  ListItemSecondaryAction,
  IconButton,
  makeStyles,
} from '@material-ui/core'
import { Link } from 'react-router-dom'
import PropTypes from 'prop-types'
import PinAvatar from './PinAvatar'
import LinkIcon from '@material-ui/icons/Link'
import { useSelector } from 'react-redux'

const useStyles = makeStyles(theme => ({
  avatarContainer: {
    minWidth: theme.spacing(5),
  },
  avatar: {
    width: theme.spacing(3),
    height: theme.spacing(3),
    fontSize: 18,
  },
}), { name: 'PinListItem' })

const PinListItem = ({ record = {}, active }) => {
  const [showAction, setShowAction] = useState(false)
  const classes = useStyles()
  const hasUrl = useMemo(() => !!record.url, [record])
  const selectedTagIds = useSelector(state => state.app.tag.selected)

  let linkTo = '/pin'
  if (selectedTagIds.length > 0) {
    linkTo += `/t/${selectedTagIds.join(',')}`
  }
  linkTo += `/${record.id}`

  return (
    <ListItem
      button
      component={Link}
      to={linkTo}
      selected={active}
      ContainerProps={{
        onMouseEnter: () => setShowAction(true),
        onMouseLeave: () => setShowAction(false),
      }}
    >
      <ListItemAvatar className={classes.avatarContainer}>
        <PinAvatar record={record} className={classes.avatar} fontSize="inherit" />
      </ListItemAvatar>
      <ListItemText
        primary={record.title || '[no title]'}
      />
      <ListItemSecondaryAction>
        {showAction && hasUrl && (
          <IconButton component="a" href={record.url} target="_blank" edge="end">
            <LinkIcon />
          </IconButton>
        )}
      </ListItemSecondaryAction>
    </ListItem>
  )
}

PinListItem.propTypes = {
  active: PropTypes.bool,
}

export default PinListItem
