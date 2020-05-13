exports.colors = {
  'background': '#f8f9fa',
  'container': '#ffffff',
  'anchor': '#0366d6',
  'text': '#2d3748',
  'tag-bg': '#f0f0f0',
  'tag-border': '#dcdcdc',

  'primary': '#3182ce',
  'error': '#e53e3e',
  'disabled': '#a1a9b3',

  'text-primary': '#2d3748',
  'text-secondary': '#5b626b',
  'text-inverted': '#ffffff',

  'btn-light-bg': '#edf2f7',
  'btn-light': '#4a5568',

  'divider': '#cbd5e0',
  'control': '#cbd5e0',

  'hover-bg': '#f4f8ff',
  'hover-border': '#62acff',

  'backdrop': '#00000070',
}

const screens = {
  'sm': '640px',
  'md': '768px',
  'lg': '1024px',
  'xl': '1280px',
  'xxl': '1440px',
}
exports.mediaQueries = {
  get: (name) => {
    const size = screens[name]
    return `(min-width: ${size})`
  },
  getAll: () => {
    return {
      ...screens,
      'sm-down': { max: '767px' },
      'md-down': { max: '1023px' },
      'lg-down': { max: '1279px' },
    }
  },
}
