import React from 'react'
import ImageAvatar from '../images/ImageAvatar'
import BookmarkIcon from '@material-ui/icons/Bookmark'

const PinAvatar = ({ fontSize, ...props }) => {
  return (
    <ImageAvatar {...props}>
      <BookmarkIcon fontSize={fontSize} />
    </ImageAvatar>
  )
}

PinAvatar.defaultProps = {
  variant: 'rounded',
  source: 'imageId',
}

export default PinAvatar
