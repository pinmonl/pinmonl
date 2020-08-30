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
  const { pinlId } = useParams()
  const dataProvider = useDataProvider()
  const [pkgs, setPkgs] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    let cancelled = false
    const fetch = async () => {
      try {
        const { data: pinl } = await dataProvider.getOne('pinl', { id: pinlId })
        if (cancelled) return
        const { data } = await dataProvider.getMany('pkg', { ids: pinl.pkgIds })
        if (cancelled) return
        setPkgs(data)
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
            value={0}
            onChange={(e, val) => console.log(val)}
          >
            {pkgs.map(pkg => (
              <Tab
                key={pkg.id}
                iconName={pkg.provider}
                pkg={pkg}
              />
            ))}
          </Tabs>
        </Box>
      </Card>
      {pkgs.length > 0 && <StatPanel pkg={pkgs[0]} />}
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
        <Box>
          {label}
        </Box>
      </Box>
    </ButtonBase>
  )
}

export default PinlPkgList
