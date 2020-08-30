import React, {
  Children,
  useState,
  useEffect,
  useCallback,
} from 'react'
import {
  SimpleForm,
  TextInput,
  ImageInput,
  ImageField,
  FormInput,
  useMutation,
  useRedirect,
  required,
  useNotify,
} from 'react-admin'
import {
  Grid,
  IconButton,
  InputAdornment,
} from '@material-ui/core'
import RefreshIcon from '@material-ui/icons/Refresh'
import { useForm, useFormState } from 'react-final-form'
import { getImageUrl } from '../images/utils'
import TagArrayInput from '../tags/TagArrayInput'

const imageKey = '_image'

const PinlForm = (props) => {
  const [mutate] = useMutation({ resource: 'pinl' })
  const [mutateImage] = useMutation({ type: 'createImage', resource: 'pinl' })
  const redirect = useRedirect()
  const notify = useNotify()

  const save = useCallback(async (data) => {
    const isNew = !data.id
    const payload = { data }
    if (!isNew) {
      payload.id = data.id
    }

    return new Promise((resolve, reject) => mutate({
      type: isNew ? 'create' : 'update',
      payload,
    }, {
      onSuccess: resolve,
      onFailure: reject,
    }))
  }, [mutate])

  const saveImage = useCallback(async (id, image) => {
    return new Promise((resolve, reject) => mutateImage({
      payload: { id, image },
    }, {
      onSuccess: resolve,
      onFailure: reject,
    }))
  }, [mutateImage])

  const handleSave = useCallback(async (values, redirectTo) => {
    const { [imageKey]: image, ...data } = values

    try {
      const { data: target } = await save(data)
      if (image && image.rawFile) {
        const { data: { id: imageId } } = await saveImage(target.id, image)
        target.imageId = imageId
      }
      redirect(redirectTo, props.basePath, target.id, target)
    } catch (e) {
      notify(
        typeof e === 'string' ? e : e.message,
        'warning'
      )
    }
  }, [save, saveImage, redirect, props.basePath, notify])

  return (
    <SimpleForm {...props} save={handleSave}>
      <FormBody />
    </SimpleForm>
  )
}

const FormBody = (props) => {
  const form = useForm()
  const [mutate] = useMutation({ type: 'getCard' })

  const urlValue = form.getState().values.url

  const fetch = useCallback(async () => {
    return new Promise((resolve) => mutate({
      payload: { url: urlValue },
    }, {
      onSuccess: resolve,
    }))
  }, [urlValue, mutate])

  const handleFetch = ({ data }) => {
    form.batch(() => {
      form.change('title', data.title)
      form.change('description', data.description)
      form.change(imageKey, data.image)
    })
  }

  return (
    <InitializeFormValues>
      <Grid container spacing={4}>
        <Grid item xs={12} md={6}>
          <FormCol {...props}>
            <PinlUrlInput source="url" fullWidth onRefreshClick={() => fetch().then(handleFetch)} validate={[required()]} />
            <TagArrayInput label="Tags" source="tagNames" fullWidth strict />
          </FormCol>
        </Grid>
        <Grid item xs={12} md={6}>
          <FormCol {...props}>
            <TextInput source="title" fullWidth validate={[required()]} />
            <TextInput source="description" multiline fullWidth />
            <ImageInput source={imageKey} fullWidth accept="image/*">
              <ImageField source="src" />
            </ImageInput>
          </FormCol>
        </Grid>
      </Grid>
    </InitializeFormValues>
  )
}

const InitializeFormValues = ({ children }) => {
  const form = useForm()
  const formState = useFormState()

  useEffect(() => {
    if (typeof formState.values[imageKey] !== 'undefined') {
      return
    }
    form.initialize((values) => {
      let image = null
      const imageUrl = getImageUrl(values, 'imageId')
      if (imageUrl) {
        image = { src: imageUrl }
      }
      values[imageKey] = image
      return values
    })
  }, [formState.values])
  
  return children
}

const PinlUrlInput = ({ onRefreshClick, ...props }) => {
  const formState = useFormState()
  const [hasUrl, setHasUrl] = useState()

  useEffect(() => {
    setHasUrl(!!formState.values.url)
  }, [formState.values.url])

  return (
    <TextInput
      {...props}
      InputProps={{
        endAdornment: hasUrl && (
          <InputAdornment position="end">
            <IconButton onClick={onRefreshClick} tabIndex="-1">
              <RefreshIcon />
            </IconButton>
          </InputAdornment>
        )
      }}
    />
  )
}

const FormCol = ({ children, ...props }) => {
  return mapFormInput(children, props)
}

const mapFormInput = (children, props) => {
  const {
    basePath,
    record,
    resource,
    variant,
    margin,
  } = props

  return Children.map(children, input => input && (
    <FormInput
      basePath={basePath}
      input={input}
      record={record}
      resource={resource}
      variant={input.props.variant || variant}
      margin={input.props.margin || margin}
      fullWidth={input.props.fullWidth}
    />
  ))
}

export default PinlForm
