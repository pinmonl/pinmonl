import React from 'react'
import { Title, useListController, ListContextProvider, useListContext } from 'react-admin'
import { Card, List, Drawer, makeStyles } from '@material-ui/core'
import { useSelector } from 'react-redux'
import { Route } from 'react-router-dom'
import clsx from 'clsx'
import PinEdit from './PinEdit'
import PinListItem from './PinListItem'

const PinListView = (props) => {
  const listContext = useListContext(props)
  const {
    ids,
    data,
    defaultTitle,
  } = listContext

  return (
    <React.Fragment>
      <Title defaultTitle={defaultTitle} />
      <List
        component="div"
      >
        {ids.map(id => (
          <PinListItem
            key={id}
            record={data[id]}
          />
        ))}
      </List>
    </React.Fragment>
  )
}

const useStyles = makeStyles(theme => ({
  root: {
    flex: '1 1 auto',
    display: 'flex',
    alignItems: 'flex-start',
  },
  list: {
    flexGrow: 1,
  },
  drawerPaper: {
    zIndex: theme.zIndex.detailDrawer,
    width: 600,
    paddingTop: 48,
  },
}))

const PinList = (props) => {
  const classes = useStyles()
  const selectedTags = useSelector(state => state.app.tag.selected)
  const controllerProps = useListController({
    ...props,
    perPage: 0,
    filter: { tagId: selectedTags },
  })

  return (
    <ListContextProvider value={controllerProps}>
      <div className={classes.root}>
        <Route path="/pin/:id">
          {({ match }) => {
            const drawerOpen = !!(match && match.params.id !== 'create')
            return (
              <React.Fragment>
                <Card
                  className={clsx(classes.list, {
                    [classes.listWithDrawer]: drawerOpen,
                  })}
                >
                  <PinListView {...props} />
                </Card>
                <Drawer
                  variant="persistent"
                  anchor="right"
                  open={drawerOpen}
                  classes={{
                    paper: classes.drawerPaper,
                  }}
                >
                  {drawerOpen && (
                    <PinEdit {...props} id={match.params.id} />
                  )}
                </Drawer>
              </React.Fragment>
            )
          }}
        </Route>
      </div>
    </ListContextProvider>
  )
}

export default PinList
