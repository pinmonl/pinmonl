import React from 'react'
import {
  useListContext,
  List,
  Filter,
  TextInput,
  useRedirect,
} from 'react-admin'
import {
  List as MuiList,
  ListItem,
  ListItemAvatar,
  ListItemText,
  Avatar,
  Link,
  Divider,
  Box,
} from '@material-ui/core'
import BookmarkIcon from '@material-ui/icons/Bookmark'
import TagArrayField from '../tags/TagArrayField'
import { getImageUrl } from '../images/utils'
import PkgChip from '../pkgs/PkgChip'
import TagArrayInput from '../tags/TagArrayInput'

const PinlList = (props) => (
  <List {...props} filters={<PinlFilter />}>
    <PinlGrid />
  </List>
)

const PinlGrid = (props) => {
  const { ids, data, ...ctxProps } = useListContext()

  return (
    <MuiList>
      {ids.map((id, n) =>
        <React.Fragment key={id}>
          {n > 0 && <Divider component="li" />}
          <PinlListItem {...props} {...ctxProps} record={data[id]} />
        </React.Fragment>
      )}
    </MuiList>
  )
}

const PinlListItem = ({
  record = {},
  basePath,
}) => {
  const imageUrl = getImageUrl(record, 'imageId')
  const redirect = useRedirect()

  const handleClick = () => {
    redirect(`/pinl/${record.id}`)
  }

  return (
    <ListItem button onClick={handleClick} alignItems="flex-start">
      <ListItemAvatar>
        <Avatar variant="rounded" src={imageUrl}>
          <BookmarkIcon />
        </Avatar>
      </ListItemAvatar>
      <Box display="flex" flexDirection={{xs: 'column', md: 'row'}} flexGrow={1}>
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
        <Box display="flex" width={{md: '300px'}} flexShrink={0} mt={1.5} flexWrap="wrap">
          {record.pkgs.map(pkg => (
            <PkgChip key={pkg.id} pkg={pkg} />
          ))}
        </Box>
      </Box>
    </ListItem>
  )
}

const PinlFilter = (props) => {
  return (
    <Filter {...props}>
      <TextInput label="Search" source="q" />
      <TagArrayInput label="Tag" source="tag" />
    </Filter>
  )
}

export default PinlList
