exports.colors = {
  'bg': '#F7FAFC',
  'nav-bg': '#EDF2F7',
  'tag-bg': '#F0F0F0',
  'clear': '#FFFFFF',

  'text-primary': '#2D3748',
  'text-secondary': '#CBD5E0',
  'text-inverted': '#FFFFFF',

  'primary': '#3182CE',
  'error': '#E53E3E',
  'divider': '#CBD5E0',
  'control': '#CBD5E0',
  'anchor': '#0366d6',

  'btn-light-bg': '#EDF2F7',
  'btn-light': '#4A5568',
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
