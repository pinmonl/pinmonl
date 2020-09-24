import React from 'react'
import { useEditController } from 'react-admin'
import PinForm from './PinForm'

const PinEdit = (props) => {
  return (
    <PinForm {...props} {...useEditController(props)} />
  )
}

export default PinEdit
