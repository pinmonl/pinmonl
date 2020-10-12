import React from 'react'
import md from '../../utils/markdown'
import get from 'lodash/get'
import { Typography } from '@material-ui/core'
import clsx from 'clsx'
import './MarkdownField.css'

const MarkdownField = ({ 
  record, 
  source, 
  value,
  className,
  ...props
}) => {
  const safeValue = get(record, source, value) || ''
  return (
    <Typography
      component="div"
      dangerouslySetInnerHTML={{__html: md.render(safeValue)}}
      className={clsx(className, 'markdown-body')}
      {...props}
    />
  )
}

MarkdownField.defaultProps = {
  variant: 'body2',
}

export default MarkdownField
