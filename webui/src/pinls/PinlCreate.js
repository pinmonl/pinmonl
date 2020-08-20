import React from 'react'
import { Create } from 'react-admin'
import PinlForm from './PinlForm'

const PinlCreate = (props) => (
  <Create {...props}>
    <PinlForm />
  </Create>
)

export default PinlCreate
