import React from 'react'
import { Form } from 'react-final-form'

const PinFormView = (props) => {
  return (
    <div>
    </div>
  )
}

const PinForm = (props) => {
  const handleSubmit = () => {
  }

  return (
    <Form
      {...props}
      onSubmit={handleSubmit}
      render={(formProps) => (
        <PinFormView {...formProps} />
      )}
    />
  )
}

export default PinForm
