import React, { useCallback } from 'react'
import {
  makeStyles,
  List,
  ListItem,
  ListItemText,
  ListItemSecondaryAction,
  Collapse,
  IconButton,
  CircularProgress,
  Typography,
  useTheme,
} from '@material-ui/core'
import clsx from 'clsx'
import UpIcon from '@material-ui/icons/KeyboardArrowUp'
import DownIcon from '@material-ui/icons/KeyboardArrowDown'
import { baseTagName } from '../utils/tag'
import useGetTagList from '../tags/useGetTagList'
import { useDispatch, useSelector } from 'react-redux'
import TagIcon from '@material-ui/icons/LabelOutlined'
import CloseIcon from '@material-ui/icons/CancelOutlined'
import {
  TAG_OPENED,
  TAG_CLOSED,
  TAG_SELECTED,
  TAG_UNSELECTED,
} from '../actions'

const TagMenu = ({ parentId = '', open = true }) => {
  const {
    ids,
    data,
    loaded,
  } = useGetTagList({ parentId })

  return (
    <Collapse
      in={loaded && open}
      unmountOnExit
    >
      <List
        component="div"
        disablePadding
      >
        {ids.map(id => (
          <TagItem
            key={id}
            tag={data[id]}
          />
        ))}
      </List>
    </Collapse>
  )
}

const useItemStyles = makeStyles(theme => ({
  root: {
    position: 'relative',
  },
  item: {
    paddingRight: theme.spacing(4),
  },
  action: {
    position: 'absolute',
    zIndex: 10,
    top: '50%',
    transform: 'translateY(-50%)',
    display: 'flex',
  },
  startAction: {},
  endAction: {
    right: 16,
  },
  itemText: {
    margin: 0,
    display: 'flex',
    alignItems: 'center',
  },
  expandButton: {
    position: 'absolute',
    left: theme.spacing(1),
    bottom: 0,
    zIndex: 1,
    margin: 'auto',
  },
}), { name: 'TagItem' })

const TagItem = ({ tag }) => {
  const { id, hasChildren, level } = tag
  const classes = useItemStyles()
  const theme = useTheme()
  const dispatch = useDispatch()
  const open = useSelector(state => state.app.tag.opened.includes(id))
  const selected = useSelector(state => state.app.tag.selected.includes(id))
  const loading = useSelector(state => state.app.tag.loading.includes(id))

  const handleOpen = useCallback(() => {
    dispatch({
      type: open ? TAG_CLOSED : TAG_OPENED,
      payload: { id },
    })
  }, [open, dispatch, id])

  const handleSelect = useCallback(() => {
    dispatch({
      type: selected ? TAG_UNSELECTED : TAG_SELECTED,
      payload: { id },
    })
  }, [selected, dispatch, id])

  return (
    <React.Fragment>
      <div className={classes.root}>
        <div
          className={clsx(classes.action, classes.startAction)}
          style={{left: theme.spacing(2 + level)}}
        >
          {selected ? (
            <IconButton size="small" onClick={handleSelect}>
              <CloseIcon fontSize="inherit" />
            </IconButton>
          ) : (
            <IconButton size="small" onClick={handleSelect}>
              <TagIcon fontSize="inherit" />
            </IconButton>
          )}
        </div>
        <ListItem
          ContainerComponent="div"
          button
          onClick={handleSelect}
          className={classes.item}
          style={{paddingLeft: theme.spacing(6 + level)}}
          selected={selected}
        >
          <ListItemText
            className={classes.itemText}
            disableTypography
          >
            <Typography variant="body2" component="span">
              {baseTagName(tag.name)}
            </Typography>
          </ListItemText>
        </ListItem>
        {hasChildren && (
          <div className={clsx(classes.action, classes.endAction)}>
            {loading ? (
              <CircularProgress size={20} />
            ) : (
              <IconButton size="small" onClick={handleOpen} edge="end">
                {open ? <UpIcon fontSize="inherit" /> : <DownIcon fontSize="inherit" />}
              </IconButton>
            )}
          </div>
        )}
      </div>
      {hasChildren && (
        <TagMenu parentId={id} open={open} />
      )}
    </React.Fragment>
  )
}

export default TagMenu
