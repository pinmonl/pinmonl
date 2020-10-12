const path = require('path')
const MonacoWebpackPlugin = require('monaco-editor-webpack-plugin')

module.exports = {
  webpack: (config) => {
    config.plugins.push(new MonacoWebpackPlugin())
    return config
  },
  paths: (paths, env) => {
    paths.appBuild = path.resolve(paths.appPath, 'dist')
    return paths
  },
}
