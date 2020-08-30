import React from 'react'
import { Route } from 'react-router'
import PinlPkgList from '../pkgs/PinlPkgList'

export default [
  <Route exact path="/pkg/of-pinl/:pinlId" component={PinlPkgList} />
]
