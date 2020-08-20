import React, {
  useState,
  useRef,
  useMemo,
  useEffect,
  useCallback,
} from 'react'
import {
  useInput,
  useMutation,
} from 'react-admin'
import {
  TextField,
  Popper,
  Paper,
  List,
  ListItem,
  makeStyles,
  ClickAwayListener,
} from '@material-ui/core'
import { sanitizeTagName, absTagName } from './utils'
import TagChip from './TagChip'

const useStyles = makeStyles((theme) => ({
  root: {
    position: 'relative',
  },
  popper: {
    width: '100%',
    zIndex: theme.zIndex.modal,
  },
  paper: {
    margin: theme.spacing(0.5, 0),
  },
  list: {
    maxHeight: '200px',
    overflow: 'auto',
  },
  inputRoot: {
    paddingTop: '20px',
    paddingLeft: '8px',
    display: 'flex',
    flexWrap: 'wrap',
  },
  inputInput: {
    paddingTop: '5px',
    paddingBottom: '4px',
    paddingLeft: '4px',
    flex: '1 1 0',
    minWidth: '60px',
  },
}), { name: 'TagArrayInput' })

const TagArrayInput = React.forwardRef(({
  label,
  variant,
  margin,
  ...props
}, ref) => {
  const [mutate] = useMutation({ type: 'getList', resource: 'tag' })
  const [inputValue, setInputValue] = useState('')
  const [options, setOptions] = useState([])
  const classes = useStyles()
  const [highlighted, setHighlighted] = useState(-1)
  const [focus, setFocus] = useState(false)
  const rootRef = useRef()
  const fieldRef = useRef()
  const inputRef = useRef()
  const listRef = useRef()
  const filteredOptions = options

  const { input } = useInput(props)
  const value = input.value || []
  const onChange = input.onChange

  const handleFetch = ({ data }) => {
    const tagNames = data.map((item) => item.name)
    setOptions(tagNames)
    setHighlighted(-1)
  }

  const fetch = useCallback(async (query) => new Promise((resolve) => mutate({
    payload: { filter: {q: query} },
  }, {
    onSuccess: resolve,
  })), [mutate])

  const getSelectedIndex = useCallback((tag) => {
    const sanitizedTag = sanitizeTagName(tag)
    return value.findIndex((item) => item === sanitizedTag)
  }, [value])

  const push = useCallback((tag) => {
    onChange([ ...value, sanitizeTagName(tag) ])
  }, [value, onChange])

  const removeAt = useCallback((index) => {
    const newValue = [...value]
    const removed = newValue.splice(index, 1)
    onChange(newValue)
    return removed
  }, [value, onChange])

  const toggle = useCallback((tag) => {
    const index = getSelectedIndex(tag)
    if (index >= 0) {
      removeAt(index)
    } else {
      push(tag)
    }
  }, [push, removeAt, getSelectedIndex])

  const changeHighlightIndex = useCallback(({ dir, step = 1, cycle = true }) => {
    let index = highlighted
    let len = filteredOptions.length
    switch (dir) {
      case 'up':
        index -= step
        if (highlighted < 0) {
          index++
        }
        break
      case 'down':
        index += step
        break
      case 'top':
        index = 0
        break
      case 'bottom':
        index = len - 1
        break
      default:
    }

    if (cycle) {
      if (index < 0) {
        index += len
      } else if (index >= len) {
        index -= len
      }
    }
    setHighlighted(index)
  }, [highlighted, filteredOptions])

  const handleKeyDown = useCallback((e) => {
    switch (e.key) {
      case 'Tab':
        setFocus(false)
        break
      case 'ArrowUp':
        e.preventDefault()
        changeHighlightIndex({ dir: 'up' })
        break
      case 'ArrowDown':
        e.preventDefault()
        changeHighlightIndex({ dir: 'down' })
        break
      case 'Backspace':
        if (value.length && !inputValue) {
          e.preventDefault()
          removeAt(value.length - 1)
        }
        break
      case 'Enter':
        if (highlighted >= 0) {
          e.preventDefault()
          toggle(filteredOptions[highlighted])
          setInputValue('')
        } else if (inputValue) {
          e.preventDefault()
          toggle(inputValue)
          setInputValue('')
        }
        break
      case 'Escape':
        e.preventDefault()
        setFocus(false)
        inputRef.current.blur()
        break
      default:
    }
  }, [changeHighlightIndex, filteredOptions, inputValue, toggle])

  const handleItemClick = useCallback((e, option) => {
    inputRef.current.focus()
    toggle(option)
    setInputValue('')
  }, [toggle])

  useEffect(() => {
    if (!inputValue.length) return

    let cancelled = false
    fetch(inputValue).then((resp) => {
      if (cancelled) return
      handleFetch(resp)
    })
    return () => cancelled = true
  }, [inputValue, fetch])

  useEffect(() => {
    if (highlighted < 0) return
    if (!listRef.current) return

    const selectedEl = listRef.current.children[highlighted]
    if (!selectedEl) return

    const scrollTo = (top) => listRef.current.scroll(0, top)

    if (selectedEl.offsetTop < listRef.current.scrollTop) {
      scrollTo(selectedEl.offsetTop)
    } else if (selectedEl.offsetTop + selectedEl.offsetHeight > listRef.current.scrollTop + listRef.current.clientHeight) {
      scrollTo(selectedEl.offsetTop + selectedEl.offsetHeight - listRef.current.clientHeight)
    }
  }, [highlighted])

  useEffect(() => {
    if (ref && ref.current) {
      ref.current = rootRef.current
    }
  }, [rootRef])

  const open = useMemo(() => focus && filteredOptions.length > 0, [focus, filteredOptions])

  return (
    <ClickAwayListener onClickAway={() => setFocus(false)}>
      <div ref={rootRef} className={classes.root}>
        <TextField
          ref={fieldRef}
          inputRef={inputRef}
          variant={variant}
          margin={margin}
          label={label}
          value={inputValue}
          onChange={(e) => setInputValue(e.target.value)}
          fullWidth={props.fullWidth}
          onKeyDown={handleKeyDown}
          onFocus={() => setFocus(true)}
          InputProps={{
            startAdornment: (value.length > 0 &&
              <SelectedValues
                value={value}
                onDelete={(tag) => toggle(tag)}
                tabIndex="-1"
                variant={variant === 'filled' ? 'outlined' : 'default'}
              />
            ),
            classes: {
              input: classes.inputInput,
              root: classes.inputRoot,
            },
          }}
        />
        {fieldRef.current && (
          <Popper
            anchorEl={fieldRef.current}
            open={open}
            disablePortal
            placement="bottom-start"
            className={classes.popper}
          >
            <Paper className={classes.paper}>
              <List className={classes.list} ref={listRef}>
                {filteredOptions.map((option, n) => (
                  <ListItem
                    onMouseEnter={() => setHighlighted(n)}
                    key={option}
                    selected={n === highlighted}
                    onClick={(e) => handleItemClick(e, option)}
                    button
                    tabIndex="-1"
                    data-selected={n === highlighted ? true : undefined}
                  >{absTagName(option)}</ListItem>
                ))}
              </List>
            </Paper>
          </Popper>
        )}
      </div>
    </ClickAwayListener>
  )
})

TagArrayInput.defaultProps = {
  variant: 'filled',
  margin: 'dense',
}

const SelectedValues = ({ value, onDelete, ...props }) => {
  return value.map(tag => (
    <TagChip
      key={tag}
      label={tag}
      onDelete={() => onDelete(tag)}
      {...props}
    />
  ))
}

export default TagArrayInput
