import React, {
  useCallback,
  useMemo,
} from 'react'
import { 
  ReferenceArrayField,
} from 'react-admin'
import {
  ListItem,
  ListItemAvatar,
  ListItemText,
  Avatar,
  Link,
  Box,
  makeStyles,
} from '@material-ui/core'
import clsx from 'clsx'
import { useHistory } from 'react-router-dom'
import BookmarkIcon from '@material-ui/icons/Bookmark'
import { getImageUrl } from '../images/utils'
import TagArrayField from '../tags/TagArrayField'
import PkgChip from '../pkgs/PkgChip'

const useStyles = makeStyles((theme) => ({
  container: {
    position: 'relative',
  },
  listItem: {
    [theme.breakpoints.up('md')]: {
      paddingRight: '240px',
    },
    [theme.breakpoints.up('lg')]: {
      paddingRight: '300px',
    },
  },
  avatar: {
    marginBottom: theme.spacing(1),
  },
  pkgs: {
    position: 'absolute',
    [theme.breakpoints.down('sm')]: {
      left: theme.spacing(9),
      bottom: theme.spacing(1),
    },
    [theme.breakpoints.up('md')]: {
      left: '100%',
      marginLeft: '-240px',
      top: theme.spacing(2.5),
    },
    [theme.breakpoints.up('lg')]: {
      marginLeft: '-300px',
    },
  },
}), { name: 'PinlListItem' })

const PinlListItem = (props) => {
  const { record } = props
  const imageUrl = getImageUrl(record, 'imageId')
  const classes = useStyles()
  const history = useHistory()

  const handleClick = useCallback(() => {
    history.push(`/pinl/${record.id}`)
  }, [record, history])

  const handlePkgClick = useCallback(() => {
    history.push(`/pkg/of-pinl/${record.id}`)
  }, [record, history])

  const hasPkgs = useMemo(() => {
    return (record.pkgIds || []).length > 0
  }, [record])

  return (
    <li className={classes.container}>
      <ListItem
        button
        onClick={handleClick}
        alignItems="flex-start"
        className={classes.listItem}
      >
        <ListItemAvatar className={classes.avatar}>
          <Avatar variant="rounded" src={imageUrl}>
            <BookmarkIcon />
          </Avatar>
        </ListItemAvatar>
        <Box display="flex" flexDirection={{xs: 'column', md: 'row'}} flexGrow={1} pb={{xs: hasPkgs ? 4 : 0, md: 0}}>
          <Box display="flex" flexDirection="column" flexGrow={1} flexShrink={1} width={{xs: 1, md: 'auto'}}>
            <Box display="flex">
              <ListItemText
                primary={
                  <Link href={record.url} target="_blank" onClick={e => e.stopPropagation()}>
                    {record.title}
                  </Link>
                }
              />
            </Box>
            <Box display="flex" flexWrap="wrap">
              <TagArrayField record={record} source="tagNames" />
            </Box>
          </Box>
        </Box>
      </ListItem>
      <Box className={classes.pkgs}>
        <ReferenceArrayField record={record} basePath={props.basePath} source="pkgIds" reference="pkg">
          <PinlPkgChip onClick={handlePkgClick} />
        </ReferenceArrayField>
      </Box>
    </li>
  )
}

const PinlPkgChip = ({ ids, data, onClick }) => {
  return ids.map(id => (
    !!data[id] && <PkgChip key={id} pkg={data[id]} onClick={onClick} />
  ))
}

export default PinlListItem
