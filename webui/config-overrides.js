const path = require('path')
const rsPaths = require('react-scripts/config/paths')

module.exports = {
  webpack: (config, env) => {
    config.resolve.alias['@'] = rsPaths.appSrc
    return config
  },

  paths: (paths, env) => {
    paths.appBuild = path.resolve(paths.appPath, 'dist')
    return paths
  },
}
