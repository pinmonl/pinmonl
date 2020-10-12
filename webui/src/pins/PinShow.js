import React from 'react'
import { TextField } from 'react-admin'
import { 
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  makeStyles,
} from '@material-ui/core'
import MarkdownField from '../components/field/MarkdownField'
import NoteIcon from '@material-ui/icons/Subject'
import PinAvatar from './PinAvatar'
import get from 'lodash/get'
import ExternalLink from '../components/field/ExternalLink'

const useStyles = makeStyles(theme => ({
  root: {
    display: 'flex',
    flexDirection: 'column',
  },
  descriptionItem: {
    alignItems: 'flex-start',
  },
}))

const PinShow = (props) => {
  const classes = useStyles()
  const { record } = props
  const fieldProps = {
    record,
  }

  return (
    <List component="div" className={classes.root}>
      <ListItem component="div">
        <ListItemIcon>
          <PinAvatar {...fieldProps} />
        </ListItemIcon>
        <ListItemText
          disableTypography
          primary={<TextField {...fieldProps} source="title" component="div" variant="body1" />}
          secondary={
            <ExternalLink href={get(record, 'url')} target="_blank">
              <TextField {...fieldProps} source="url" component="div" color="textSecondary" />
            </ExternalLink>
          }
        />
      </ListItem>
      <ListItem className={classes.descriptionItem} component="div">
        <ListItemIcon>
          <NoteIcon />
        </ListItemIcon>
        <ListItemText
          disableTypography
          primary={
            <MarkdownField
              {...fieldProps}
              source="description"
            />
          }
        />
      </ListItem>
    </List>
  )
}

export default PinShow
