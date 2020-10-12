import React, { useEffect, useCallback, useRef } from 'react'
import {
  FormControl,
  InputLabel,
  makeStyles,
} from '@material-ui/core'
import clsx from 'clsx'
import { useInput, useTranslate } from 'react-admin'
import * as monaco from 'monaco-editor'

const useStyles = makeStyles(theme => ({
  root: {
    position: 'relative',
    display: 'flex',
    flexDirection: 'column',
  },
  label: {
  },
  content: {
    marginTop: theme.spacing(2),
    flex: '1 1 0',
  },
  button: {
    position: 'absolute',
    bottom: 4,
    right: 4,
  },
}), { name: 'MarkdownInput' })

const MarkdownInput = (props) => {
  const classes = useStyles(props)
  const { className, resource, source } = props
  const editorRef = useRef(null)
  const nodeRef = useRef()
  const translate = useTranslate()

  const {
    input: { value, onChange },
  } = useInput(props)

  const updateLayout = useCallback(() => {
    if (!editorRef.current) return
    console.log('hihi')
    editorRef.current.layout()
  }, [])

  useEffect(() => {
    if (editorRef.current) return
    editorRef.current = monaco.editor.create(nodeRef.current, {
      language: 'markdown',
      fontSize: 12,
      wordWrap: 'on',
    })

    return () => {
      editorRef.current.dispose()
      editorRef.current = null
    }
  }, [])

  useEffect(() => {
    if (!editorRef.current) return
    if (editorRef.current.getValue() !== value) {
      editorRef.current.setValue(value)
    }
  }, [value])

  useEffect(() => {
    if (!editorRef.current) return
    const disposableChangeHandler = editorRef.current.onDidChangeModelContent(() => {
      onChange(editorRef.current.getValue())
    })
    return () => {
      disposableChangeHandler.dispose()
    }
  }, [onChange])

  useEffect(() => {
    window.addEventListener('resize', updateLayout)
    return () => window.removeEventListener('resize', updateLayout)
  }, [updateLayout])

  return (
    <FormControl className={clsx(className, classes.root)}>
      <InputLabel
        shrink
        className={classes.label}
      >
        {translate(`resources.${resource}.fields.${source}`, {
          _: source,
        })}
      </InputLabel>
      <div
        ref={nodeRef}
        className={classes.content}
      />
    </FormControl>
  )
}

export default MarkdownInput
