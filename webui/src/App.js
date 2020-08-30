import React from 'react'
import {
  AdminContext,
  AdminUI,
  Resource,
} from 'react-admin'
import './App.css'
import { Pubsub } from './pubsub'

import Layout from './Layout'
import dataProvider from './dataProvider'
import { Login, authProvider } from './auth'
import { i18nProvider } from './i18n'
import customRoutes from './custom/routes'
import pinls from './pinls'

function App() {
  return (
    <AdminContext
      dataProvider={dataProvider}
      i18nProvider={i18nProvider}
      authProvider={authProvider}
    >
      <Pubsub>
        <AdminUI customRoutes={customRoutes} layout={Layout} loginPage={Login}>
          <Resource name="pinl" {...pinls} />
          <Resource name="tag" />
          <Resource name="pkg" />
          <Resource name="stat" />
        </AdminUI>
      </Pubsub>
    </AdminContext>
  )
}

export default App
