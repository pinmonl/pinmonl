import React from 'react'
import { Admin, Resource } from 'react-admin'
import './App.css'
import dataProvider from './dataProvider'

import pinls from './pinls'

function App() {
  return (
    <Admin dataProvider={dataProvider}>
      <Resource name="pinl" {...pinls} />
    </Admin>
  )
}

export default App
