import React from 'react'
import { Avatar } from '@material-ui/core'
import getImageUrl from './getImageUrl'

const ImageAvatar = ({ record, source, value, ...props }) => {
  const imageSrc = getImageUrl({ record, source, value })
  return (
    <Avatar
      src={imageSrc}
      {...props}
    />
  )
}

export default ImageAvatar
