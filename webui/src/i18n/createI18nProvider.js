import polyglotI18nProvider from 'ra-i18n-polyglot'
import enMessage from './en'

const createI18nProvider = () => {
  const provider = polyglotI18nProvider(locale => {
    if (locale === 'en') {
      return enMessage
    }
  }, 'en')
  return provider
}

export default createI18nProvider
