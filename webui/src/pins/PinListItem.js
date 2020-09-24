import React, { useMemo } from 'react'
import {
  ListItem,
  ListItemText,
  ListItemAvatar,
  Avatar,
  makeStyles,
} from '@material-ui/core'
import BookmarkIcon from '@material-ui/icons/Bookmark'
import { Link } from 'react-router-dom'

const useStyles = makeStyles(theme => ({
  avatarContainer: {
    minWidth: theme.spacing(5),
  },
  avatar: {
    width: theme.spacing(3),
    height: theme.spacing(3),
  },
  avatarIcon: {
    fontSize: theme.typography.pxToRem(18),
  },
}), { name: 'PinListItem' })

const PinListItem = ({ record }) => {
  const classes = useStyles()
  const linkTo = useMemo(() => record ? `/pin/${record.id}` : null, [record])
  return (
    <ListItem
      button
      component={Link}
      to={linkTo}
    >
      <ListItemAvatar className={classes.avatarContainer}>
        <Avatar className={classes.avatar} variant="rounded">
          <BookmarkIcon className={classes.avatarIcon} />
        </Avatar>
      </ListItemAvatar>
      <ListItemText
        primary={record.title || '[no title]'}
      />
    </ListItem>
  )
}

export default PinListItem
