export default (opts) => {
  const { prop = 'value', event = 'input' } = opts || {}

  return {
    model: opts,
    computed: {
      $m () { 
        return opts
      },
    },
    methods: {
      getter (key) {
        const model = this[prop] || {}
        return model[key]
      },
      setter (key, value) {
        const model = this[prop]
        this.$emit(event, { ...model, [key]: value })
      },
    },
  }
}
