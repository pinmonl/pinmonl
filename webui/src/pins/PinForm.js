import React, { useCallback, useState } from 'react'
import { Form } from 'react-final-form'
import {
  IconButton,
  makeStyles,
} from '@material-ui/core'
import {
  TextInput,
} from 'react-admin'
import MarkdownInput from '../components/input/MarkdownInput'
import CloseIcon from '@material-ui/icons/Close'
import EditIcon from '@material-ui/icons/Edit'
import SaveIcon from '@material-ui/icons/Save'
import PinShow from './PinShow'

const useStyles = makeStyles(theme => ({
  container: {
    display: 'block',
    overflow: 'hidden auto',
    position: 'relative',
    width: '100%',
    height: '100%',
    boxSizing: 'border-box',
    [theme.breakpoints.up('sm')]: {
      padding: theme.spacing(4, 2),
    },
    [theme.breakpoints.down('xs')]: {
      paddingTop: theme.spacing(3),
    },
  },
  form: {
    display: 'flex',
    flexDirection: 'column',
    position: 'relative',
    width: '100%',
    minHeight: '100%',
  },
  markdownInput: {
    flex: '1 1 0',
    minHeight: 500,
  },
  actions: {
    position: 'absolute',
    top: 4,
    right: 4,
    display: 'flex',
    flexDirection: 'row-reverse',
  },
}), { name: 'PinForm' })

const PinFormView = (props) => {
  const classes = useStyles()
  const [edit, setEdit] = useState(false)
  const {
    handleSubmit,
    resource,
    basePath,
    record,
    onClose,
  } = props

  const inputProps = {
    record,
    resource,
    basePath,
    margin: 'dense',
    variant: 'standard',
  }

  const handleSave = useCallback(async () => {
    await handleSubmit()
    setEdit(false)
  }, [handleSubmit])

  return (
    <div className={classes.container}>
      {edit ? (
        <form onSubmit={handleSubmit} className={classes.form}>
          <TextInput {...inputProps} source="title" />
          <TextInput {...inputProps} source="url" />
          <MarkdownInput
            {...inputProps}
            source="description"
            className={classes.markdownInput}
          />
        </form>
      ) : (
        <PinShow {...props} />
      )}

      <div className={classes.actions}>
        {edit ? (
          <React.Fragment>
            <IconButton onClick={() => setEdit(false)}>
              <CloseIcon />
            </IconButton>
            <IconButton onClick={handleSave}>
              <SaveIcon />
            </IconButton>
          </React.Fragment>
        ) : (
          <React.Fragment>
            <IconButton onClick={(e) => onClose(e)}>
              <CloseIcon />
            </IconButton>
            <IconButton onClick={() => setEdit(true)}>
              <EditIcon />
            </IconButton>
          </React.Fragment>
        )}
      </div>
    </div>
  )
}

const PinForm = (props) => {
  const { record, save } = props

  const handleSubmit = useCallback(async (values) => {
    return new Promise((resolve, reject) => {
      save(values, '', {
        onSuccess: () => {
          resolve()
        },
        onFailure: () => {
          reject()
        },
      })
    })
  }, [save])

  const validate = useCallback((values) => {
  }, [])

  return (
    <Form
      {...props}
      initialValues={record}
      onSubmit={handleSubmit}
      validate={validate}
      render={(formProps) => (
        <PinFormView {...formProps} />
      )}
    />
  )
}

export default PinForm
