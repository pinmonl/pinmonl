import React from 'react'
import {
  AdminContext,
  AdminUI,
  Resource,
} from 'react-admin'
import './App.css'

import Layout from './Layout'
import dataProvider from './dataProvider'
import { i18nProvider } from './i18n'
import { Pubsub } from './pubsub'
import { Auth } from './auth'
import customRoutes from './custom/routes'
import pinls from './pinls'

function App() {
  return (
    <AdminContext
      dataProvider={dataProvider}
      i18nProvider={i18nProvider}
    >
      <Auth>
        <Pubsub>
          <AdminUI customRoutes={customRoutes} layout={Layout}>
            <Resource name="pinl" {...pinls} />
            <Resource name="tag" />
            <Resource name="pkg" />
            <Resource name="stat" />
          </AdminUI>
        </Pubsub>
      </Auth>
    </AdminContext>
  )
}

export default App
