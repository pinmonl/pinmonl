import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'

const md = MarkdownIt('commonmark', {
  highlight: (str, lang) => {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return hljs.highlight(lang, str).value
      } catch (_) {}
    }
    return ''
  }
})

export default md
