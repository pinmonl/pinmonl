const themeVars = require('./src/theme/variables.js')

function scssPrependData() {
  let pd = []

  const colors = []
  for (const [key, value] of Object.entries(themeVars.colors)) {
    colors.push(`${key}: ${value}`)
  }
  pd.push(`$colors: (${colors.join(',')});`)

  pd.push(`@import '@/theme/functions.scss';`)
  return pd.join('\n')
}

module.exports = {
  css: {
    requireModuleExtension: true,
    loaderOptions: {
      css: {
        modules: {
          localIdentName: process.env.NODE_ENV == 'production' ? '[hash:6]' : '[name]_[local]_[hash:4]',
        },
      },
      scss: {
        prependData: scssPrependData,
      },
    },
  },
  chainWebpack: (config) => {
    const svgRule = config.module.rule('svg')
 
    svgRule.uses.clear()
 
    svgRule
      .use('babel-loader')
      .loader('babel-loader')
      .end()
      .use('vue-svg-loader')
      .loader('vue-svg-loader')
  },
}
