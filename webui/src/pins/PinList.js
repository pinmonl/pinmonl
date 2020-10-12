import React from 'react'
import {
  Title,
  useListController,
  ListContextProvider,
  useListContext,
  useTranslate,
} from 'react-admin'
import {
  Card,
  List,
  ListItem,
  Drawer,
  Typography,
  makeStyles,
} from '@material-ui/core'
import { useSelector } from 'react-redux'
import { Route, useHistory, useParams } from 'react-router-dom'
import PinEdit from './PinEdit'
import PinCreate from './PinCreate'
import PinListItem from './PinListItem'
import PinFilter from './PinFilter'

const useStyles = makeStyles(theme => ({
  root: {
    flex: '1 1 auto',
    display: 'flex',
    alignItems: 'flex-start',
    position: 'relative',
    overflowX: 'hidden',
  },
  list: {
    flex: '0 0 auto',
    alignSelf: 'stretch',
    width: '100%',
    position: 'relative',
    [theme.breakpoints.up('sm')]: {
      width: 400,
    },
  },
  listContent: {
    position: 'absolute',
    top: 0,
    left: 0,
    width: '100%',
    height: '100%',
    overflow: 'auto',
    padding: theme.spacing(1),
    boxSizing: 'border-box',
    '-webkit-overflow-scrolling': 'touch',
  },
  paper: {
    flex: '1 1 auto',
    height: '100%',
    [theme.breakpoints.up('sm')]: {
      position: 'relative',
    },
  },
  drawerPaper: {
    zIndex: theme.zIndex.detailDrawer,
    width: '100%',
    position: 'absolute',
  },
}), { name: 'PinList' })

const PinListView = (props) => {
  const classes = useStyles()
  const translate = useTranslate()
  const history = useHistory()
  const { id: activeId } = useParams()
  const listContext = useListContext(props)
  const {
    ids,
    data,
    total,
    loading,
    resource,
    defaultTitle,
  } = listContext
  const drawerOpen = (!!activeId && activeId !== 'create')

  return (
    <div className={classes.root}>
      <div className={classes.list}>
        <div className={classes.listContent}>
          <PinFilter />
          <Card>
            <Title defaultTitle={defaultTitle} />
            <List
              component="div"
            >
              {ids.map(id => (
                <PinListItem
                  key={id}
                  record={data[id]}
                  active={activeId === id}
                />
              ))}
              {!loading && total === 0 && (
                <ListItem>
                  <Typography>
                    {translate('ra.page.empty', {
                      name: translate(`resources.${resource}.name`, {
                        _: resource,
                      }),
                    })}
                  </Typography>
                </ListItem>
              )}
            </List>
          </Card>
        </div>
      </div>
      <Drawer
        variant="persistent"
        anchor="right"
        open={drawerOpen}
        classes={{
          root: classes.paper,
          paper: classes.drawerPaper,
        }}
      >
        {drawerOpen && (
          <PinEdit
            {...props}
            id={activeId}
            onClose={() => history.push('/pin')}
          />
        )}
      </Drawer>
    </div>
  )
}

const PinList = (props) => {
  const selectedTags = useSelector(state => state.app.tag.selected)
  const controllerProps = useListController({
    ...props,
    perPage: 0,
    filter: { tagId: selectedTags },
  })

  return (
    <ListContextProvider value={controllerProps}>
        <Route
          path={[
            '/pin/t/:tagIds/:id',
            '/pin/t/:tagIds',
            '/pin/:id',
            '/pin',
          ]}
        >
          <PinListView {...props} />
        </Route>
    </ListContextProvider>
  )
}

export default PinList
