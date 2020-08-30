import React from 'react'
import { Route } from 'react-router'
import PinlPkgList from './PinlPkgList'

const PkgList = (props) => (
  <React.Fragment>
    <Route path="/pkg/of-pinl/:pinlId" render={PinlPkgList} />
  </React.Fragment>
)

export default PkgList
