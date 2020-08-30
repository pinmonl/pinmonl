import polyglotI18nProvider from 'ra-i18n-polyglot'
import englishMessages from 'ra-language-english'

const messages = {
  ...englishMessages,
  resources: {
    pinl: {
      name: 'Pin |||| Pins',
    },
  },
}

const provider = polyglotI18nProvider(() => messages, 'en')

export default provider
