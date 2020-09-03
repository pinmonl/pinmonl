import React, {
  useState,
  useCallback,
  useEffect,
  useMemo,
} from 'react'
import {
  useGetMany,
  useInput,
  useMutation,
  FieldTitle,
  ArrayInput,
  SimpleFormIterator,
  TextInput,
  SelectInput,
  useDataProvider,
  required,
  AutocompleteInput,
} from 'react-admin'
import {
  FormControl,
  InputLabel,
  Paper,
  makeStyles,
  List,
  ListItem,
  ListItemText,
  ListItemAvatar,
  ListItemSecondaryAction,
  Avatar,
  Divider,
  Select,
  MenuItem,
  Box,
  IconButton,
  TextField,
  Dialog,
  DialogTitle,
  DialogActions,
  DialogContent,
} from '@material-ui/core'
import {
  Form,
  useField,
  useForm,
  useFormState,
} from 'react-final-form'
import get from 'lodash/get'
import MonlerIcon from '../monlers/MonlerIcon'
import { icons } from '../monlers/icons'
import EditIcon from '@material-ui/icons/Edit'
import AddIcon from '@material-ui/icons/Add'
import DoneIcon from '@material-ui/icons/Done'
import CloseIcon from '@material-ui/icons/Close'
import DeleteIcon from '@material-ui/icons/Delete'

const useStyles = makeStyles((theme) => ({
  label: {
    position: 'relative',
  },
  itemAvatar: {
    backgroundColor: 'white',
  },
}))

const PkgArrayInput = (props) => {
  const {
    record,
    source,
    variant,
    margin,
    fullWidth,
    className,
    label,
  } = props
  const pkgIds = get(record, source, [])
  const { data: pkgData, loading, error } = useGetMany('pkg', pkgIds)
  const classes = useStyles()
  const dataProvider = useDataProvider()
  const form = useForm()
  const [pkgs, setPkgs] = useState([])
  const [formValues, setFormValues] = useState(null)

  const { input, meta, isRequired } = useInput(props)

  const edit = useCallback((pkg) => {
    setFormValues({ ...pkg, uri: pkg.providerUri })
  }, [setFormValues])

  const isNew = useMemo(() => !(formValues || {}).id, [formValues])
  const formIndex = useMemo(() => {
    return isNew ? pkgs.length : pkgs.findIndex((pkg) => pkg.id === formValues.id)
  }, [isNew, formValues, pkgs])

  const handleChange = useCallback((newPkgs) => {
    setPkgs(newPkgs)
    input.onChange(newPkgs.map((p) => p.id))
    form.change('hasPinpkgs', true)
  }, [setPkgs, input])

  const handleSubmit = useCallback(async (values) => {
    try {
      const { data } = await dataProvider.create('pkg', { data: values })
      const newPkgs = [ ...pkgs ]
      newPkgs.splice(formIndex, isNew ? 0 : 1, data)
      handleChange(newPkgs)
    } finally {
      setFormValues(null)
    }
  }, [formValues, dataProvider, pkgs, isNew, formIndex])

  const handleDelete = useCallback(() => {
    const newPkgs = [ ...pkgs ]
    newPkgs.splice(formIndex, 1)
    handleChange(newPkgs)
    setFormValues(null)
  }, [formValues, pkgs, formIndex])

  useEffect(() => {
    if (loading) return
    setPkgs(pkgData)
  }, [pkgData, loading])

  return (
    <FormControl
      variant={variant}
      margin={margin}
      fullWidth={fullWidth}
      className={className}
    >
      <InputLabel shrink maring={margin} className={classes.label}>
        <FieldTitle record={record} source={source} label={label} isRequired={isRequired} />
      </InputLabel>
      <Paper variant="outlined">
        <List>
          {!loading && pkgs.map((pkg) => (
            <ListItem key={pkg.id}>
              <ListItemAvatar>
                <Avatar className={classes.itemAvatar}>
                  <MonlerIcon name={pkg.provider} />
                </Avatar>
              </ListItemAvatar>
              <ListItemText>
                {pkg.provider == 'git' && `${pkg.providerHost}/`}
                {pkg.providerUri}
              </ListItemText>
              <ListItemSecondaryAction>
                <IconButton edge="end" onClick={(e) => edit(pkg)}>
                  <EditIcon />
                </IconButton>
              </ListItemSecondaryAction>
            </ListItem>
          ))}
          {pkgs.length > 0 && <Divider />}
          <Box display="flex" justifyContent="center">
            <IconButton onClick={() => setFormValues({})}>
              <AddIcon />
            </IconButton>
          </Box>
          {!!formValues && (
            <PkgForm
              initialValues={formValues}
              onSubmit={handleSubmit}
              isNew={isNew}
              onClose={() => setFormValues(null)}
              onDelete={handleDelete}
            />
          )}
        </List>
      </Paper>
    </FormControl>
  )
}

const usePkgFormStyle = makeStyles((theme) => ({

}))

const PkgForm = ({ onClose, onDelete, isNew, ...props }) => {
  return (
    <Form {...props}>
      {({ handleSubmit }) => (
        <Dialog open={true} onClose={onClose}>
          <DialogTitle>{isNew ? 'Add Pkg' : 'Edit Pkg'}</DialogTitle>
          <DialogContent>
            <Box display="flex" flexWrap="wrap" alignItems="center" width={1}>
              <Box width={1}>
                <ProviderSelect validate={[required()]} />
              </Box>
              <Box width={1}>
                <URIInput fullWidth validate={[required()]} />
              </Box>
            </Box>
          </DialogContent>
          <DialogActions>
            {!isNew && (
              <IconButton onClick={onDelete}>
                <DeleteIcon />
              </IconButton>
            )}
            <Box flexGrow={1} />
            <IconButton onClick={onClose}>
              <CloseIcon />
            </IconButton>
            <IconButton onClick={handleSubmit} color="primary">
              <DoneIcon />
            </IconButton>
          </DialogActions>
        </Dialog>
      )}
    </Form>
  )
}

const ProviderSelect = (props) => {
  const choices = useMemo(() => {
    return Object.entries(icons).map(([name, icon]) => ({ ...icon, name }))
  })

  return (
    <SelectInput
      {...props}
      source="provider"
      choices={choices}
      optionValue="name"
      optionText={<ProviderChoice />}
    />
  )
}

const ProviderChoice = ({ record }) => {
  const { name, title } = record
  return (
    <Box display="flex" alignItems="center">
      <MonlerIcon name={name} size="small" />
      <Box ml={1} fontSize="14px">{title}</Box>
    </Box>
  )
}

const URIInput = (props) => {
  const formState = useFormState()

  const label = useMemo(() => {
    const provider = formState.values.provider
    switch (provider) {
      case 'youtube':
        return 'Channel ID'
      case 'git':
        return 'Url'
      default:
        return 'Uri'
    }
  }, [formState.values.provider])

  return (
    <TextInput
      {...props}
      source="uri"
      label={label}
    />
  )
}

export default PkgArrayInput
