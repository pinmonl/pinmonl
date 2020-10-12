import React from 'react'
import {
  useListContext,
} from 'react-admin'
import {
  Toolbar,
  IconButton,
  TextField,
  InputAdornment,
  makeStyles,
} from '@material-ui/core'
import { Form, useField } from 'react-final-form'
import FilterIcon from '@material-ui/icons/FilterList'
import SearchIcon from '@material-ui/icons/Search'

const useStyles = makeStyles(theme => ({
  toolbar: {
    paddingLeft: theme.spacing(1.5),
    paddingRight: theme.spacing(1.5),
  },
  search: {
    flex: '1 1 auto',
  },
}))

const PinFilterView = (props) => {
  const classes = useStyles()
  const {
    input,
  } = useField('q')

  return (
    <Toolbar variant="dense" className={classes.toolbar}>
      <TextField
        {...input}
        className={classes.search}
        InputProps={{
          startAdornment: (
            <InputAdornment position="start">
              <SearchIcon />
            </InputAdornment>
          ),
        }}
      />
      <IconButton edge="end">
        <FilterIcon />
      </IconButton>
    </Toolbar>
  )
}

const PinFilter = (props) => {
  const listContext = useListContext()

  return (
    <Form
      onSubmit={() => {}}
      render={(formProps) => (
        <PinFilterView {...formProps} />
      )}
    />
  )
}

export default PinFilter
