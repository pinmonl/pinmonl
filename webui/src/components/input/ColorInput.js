import React from 'react'
import { useInput } from 'react-admin'
import { TwitterPicker } from 'react-color'

const ColorInput = (props) => {
  const {
    input: { value, onChange },
  } = useInput(props)

  const handleChange = (color, event) => {
    onChange(color.hex)
  }

  return (
    <TwitterPicker
      color={value}
      onChangeComplete={handleChange}
    />
  )
}

export default ColorInput
