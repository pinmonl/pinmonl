import React, {
  useState,
  useEffect,
  useCallback,
  useMemo,
} from 'react'
import {
  useDataProvider,
} from 'react-admin'
import {
  Card,
  Tabs,
  Box,
  ButtonBase,
} from '@material-ui/core'
import { useParams } from 'react-router-dom'
import MonlerIcon from '../monlers/MonlerIcon'
import StatPanel from '../stats/StatPanel'

const PinlPkgList = (props) => {
  const { pinlId, tab: initialTab } = useParams()
  const dataProvider = useDataProvider()
  const [pkgs, setPkgs] = useState([])
  const [tab, setTab] = useState(initialTab ? Number(initialTab) : 0)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    let cancelled = false
    const fetch = async () => {
      try {
        const { data: pinl } = await dataProvider.getOne('pinl', { id: pinlId })
        if (cancelled) return
        const { data } = await dataProvider.getMany('pkg', { ids: pinl.pkgIds })
        const ordered = pinl.pkgIds.map((pkgId) => data.find((pkg) => pkg.id === pkgId))
        if (cancelled) return
        setPkgs(ordered)
      } catch (e) {
        //
      } finally {
        setLoading(false)
      }
    }
    fetch()
    return () => cancelled = true
  }, [])

  return !loading && (
    <React.Fragment>
      <Card>
        <Box>
          <Tabs
            variant="scrollable"
            value={tab}
            onChange={(e, val) => setTab(val)}
          >
            {pkgs.map((pkg, n) => (
              <Tab
                key={pkg.id}
                iconName={pkg.provider}
                pkg={pkg}
                value={n}
              />
            ))}
          </Tabs>
        </Box>
      </Card>
      {pkgs.length > 0 && <StatPanel key={tab} pkg={pkgs[tab]} />}
    </React.Fragment>
  )
}

const Tab = ({ pkg, iconName, children, value, onChange, classes }) => {
  const handleClick = useCallback((e) => {
    onChange(e, value)
  }, [onChange, value])

  const label = useMemo(() => {
    switch (pkg.provider) {
      case 'youtube':
        return pkg.title
      default:
        return pkg.providerUri
    }
  }, [pkg])

  return (
    <ButtonBase onClick={handleClick}>
      <Box px={2.5} py={1.5} display="flex" alignItems="center">
        {!!iconName && (
          <Box mr={0.5} display="flex">
            <MonlerIcon name={iconName} />
          </Box>
        )}
        <Box fontSize={14}>
          {label}
        </Box>
      </Box>
    </ButtonBase>
  )
}

export default PinlPkgList
