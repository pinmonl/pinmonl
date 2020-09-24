import React from 'react'
import { useCreateController } from 'react-admin'
import PinForm from './PinForm'

const PinCreate = (props) => {
  return (
    <PinForm {...props} {...useCreateController(props)} />
  )
}

export default PinCreate
