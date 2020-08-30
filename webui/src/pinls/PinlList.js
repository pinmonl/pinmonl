import React from 'react'
import {
  useListContext,
  List,
  Filter,
  TextInput,
  ListActions,
} from 'react-admin'
import {
  List as MuiList,
  Divider,
} from '@material-ui/core'
import TagArrayInput from '../tags/TagArrayInput'
import PinlListItem from './PinlListItem'

const PinlList = (props) => {
  return (
    <List {...props} filters={<PinlFilter />} actions={<ListActions exporter={false} />}>
      <PinlGrid />
    </List>
  )
}

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

const PinlFilter = (props) => {
  return (
    <Filter {...props}>
      <TextInput label="Search" source="q" />
      <TagArrayInput label="Tag" source="tag" />
    </Filter>
  )
}

export default PinlList
