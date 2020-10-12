import React, { useMemo } from 'react'
import {
  AdminContext,
  AdminUI,
  Resource,
} from 'react-admin'
import { createMuiTheme } from '@material-ui/core'
import { unstable_createMuiStrictModeTheme } from '@material-ui/core/styles'
import { blue, deepOrange } from '@material-ui/core/colors'
import './App.css'

import Layout from './layout/Layout'
import { baseURL, hasDefaultUser } from './utils/constants'
import { createDataProvider } from './data'
import { createAuthProvider } from './auth'
import { createI18nProvider } from './i18n'
import { createSagas } from './sideEffects'
import reducers from './reducers'
import pins from './pins'

const dataProvider = createDataProvider(baseURL)
const authProvider = createAuthProvider(baseURL, hasDefaultUser)
const i18nProvider = createI18nProvider()
const sagas = createSagas(authProvider)
const createTheme = process.env.NODE_ENV === 'production'
  ? createMuiTheme
  : unstable_createMuiStrictModeTheme

function App() {
  const theme = useMemo(() => createTheme({
    sidebar: {
      width: 240,
    },
    palette: {
      type: 'light',
      primary: {
        main: blue[800],
      },
      secondary: {
        main: deepOrange[500],
      },
      background: {
        main: '#fafafa',
      },
    },
    typography: {
      fontFamily: [
        '-apple-system',
        'BlinkMacSystemFont',
        '"Segoe UI"',
        'Roboto',
        '"Helvetica Neue"',
        'Arial',
        'sans-serif',
        '"Apple Color Emoji"',
        '"Segoe UI Emoji"',
        '"Segoe UI Symbol"',
      ].join(','),
      body1: {
        fontSize: 14,
      },
      body2: {
        fontSize: 12,
      },
    },
    zIndex: {
      detailDrawer: 200,
      sidebar: 100,
    },
  }), [])

  return (
    <AdminContext
      dataProvider={dataProvider}
      authProvider={authProvider}
      i18nProvider={i18nProvider}
      customSagas={sagas}
      customReducers={reducers}
    >
      <AdminUI
        theme={theme}
        layout={Layout}
      >
        <Resource name="pin" {...pins} />
        <Resource name="tag" />
        <Resource name="pkg" />
        <Resource name="stat" />
      </AdminUI>
    </AdminContext>
  )
}

export default App
