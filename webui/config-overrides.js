const path = require('path')

module.exports = {
  paths: (paths, env) => {
    paths.appBuild = path.resolve(paths.appPath, 'dist')
    return paths
  },
}
