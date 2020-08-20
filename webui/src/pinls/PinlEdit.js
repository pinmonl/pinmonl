import React from 'react'
import { Edit } from 'react-admin'
import PinlForm from './PinlForm'

const PinlEdit = (props) => (
  <Edit {...props} undoable={false}>
    <PinlForm />
  </Edit>
)

export default PinlEdit
