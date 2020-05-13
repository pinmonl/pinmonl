const themeVars = require('./src/theme/variables.js')
const boxGutter = '10px'

module.exports = {
  separator: '_',
  theme: {
    extend: {
      width: {
        nav: '240px',
      },
      height: {
        header: '55px',
        input: '38px',
      },
      colors: () => {
        const colors = {}
        for (const key of Object.keys(themeVars.colors)) {
          colors[`${key}`] = `var(--color-${key})`
        }
        return  colors
      },
      screens: () => {
        return themeVars.mediaQueries.getAll()
      },
      zIndex: () => {
        return {
          '-10': '-10',
          '-1': '-1',
          '100': '100',
          '200': '200',
          '300': '300',
          '400': '400',
          '500': '500',
        }
      },
      boxShadow: () => {
        return {
          'l': '-0.5px 0 0.5px rgba(0, 0, 0, 0.05), -4px 0 4px rgba(0, 0, 0, 0.1)',
          'b': '0 0.5px 0.5px rgba(0, 0, 0, 0.05), 0 4px 4px rgba(0, 0, 0, 0.1)',
          'b-sm': '0 0.3px 0.3px rgba(0, 0, 0, 0.05), 0 2px 2px rgba(0, 0, 0, 0.1)',
          'r-sm': '0.3px 0 0.3px rgba(0, 0, 0, 0.05), 2px 0 2px rgba(0, 0, 0, 0.1)',
          'box': '0 0 5px 0 #00000033',
          'modal': `0 0 9px 0 ${themeVars.colors.backdrop}`,
        }
      },
      margin: {
        'box': boxGutter,
      },
      padding: {
        'box': boxGutter,
      },
      inset: {
        'box': boxGutter,
      },
      lineHeight: {
        '0': '0px',
      },
    },
  },
  variants: {},
  plugins: [],
}
