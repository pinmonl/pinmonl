import React, {
  Children,
  useCallback,
} from 'react'
import {
  Toolbar,
  SaveButton,
  DeleteButton,
  useMutation,
} from 'react-admin'
import { Box } from '@material-ui/core'
import { useForm } from 'react-final-form'

const PinlToolbar = (props) => {
  return (
    <Toolbar {...props}>
      <ToolbarBody>
        <PinlSaveButton />
        <DeleteButton />
      </ToolbarBody>
    </Toolbar>
  )
}

const ToolbarBody = ({ children, ...props }) => (
  <Box display="flex" justifyContent="space-between" width={1}>
    {Children.map(children, action => React.cloneElement(action, props))}
  </Box>
)

const PinlSaveButton = (props) => {
  const [create] = useMutation({
    type: 'create',
    resource: 'pinl',
  })
  const [createImage] = useMutation({
    type: 'createImage',
    resource: 'pinl',
  })

  const handleSubmit = useCallback((redirectTo) => {
    // console.log(args)
    //
  }, [create, createImage])

  return (
    <SaveButton {...props} />
  )
}

export default PinlToolbar
