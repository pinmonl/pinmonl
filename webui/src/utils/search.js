export class SearchParams {
  constructor () {
    this.params = []
    this.text = ''
  }

  parse (text) {
    text = text || ''

    this.clear()
    for (const part of text.split(' ')) {
      let [ type, value ] = part.split(':')
      // This implied plain search when value is undefined.
      if (typeof value == 'undefined') {
        const decVal = decodeURIComponent(type)
        if (this.text) {
          this.text += ' '
        }
        this.text += decVal
      // Normal search param.
      } else {
        this.add(type, decodeURIComponent(value))
      }
    }
  }

  add (type, value) {
    this.params.push([type, value])
  }

  pop () {
    return this.params.pop()
  }

  clear () {
    this.params = []
    this.text = ''
  }

  get (type) {
    return this.params.filter(([ paramType ]) => type == paramType)
  }

  getValues (type) {
    return this.get(type).map(([ , paramValue ]) => paramValue)
  }

  del (type) {
    this.params = this.params.filter(([ paramType ]) => type != paramType)
  }

  getText () {
    return this.text
  }

  setText (value) {
    this.text = value
  }

  isEmpty () {
    return this.isEmptyParam() && !this.getText()
  }

  isEmptyParam () {
    return this.params.length == 0
  }

  clone () {
    const cloned = new SearchParams()
    cloned.parse(this.encode())
    return cloned
  }

  encode () {
    let str = ''
    for (const [type, value] of this.params) {
      if (str != '') {
        str += ' '
      }
      str += `${type}:${encodeURIComponent(value)}`
    }

    if (this.getText()) {
      str += ` ${encodeURIComponent(this.getText())}`
    }

    return str
  }

  [Symbol.iterator] () {
    return this.params.values()
  }
}
