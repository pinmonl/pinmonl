import isEqual from 'lodash/isEqual'
import cloneDeep from 'lodash/cloneDeep'

export default (opts) => {
  const { prop = 'value', event = 'input' } = opts || {}

  return {
    model: opts,
    data () {
      return {
        localModel: null,
      }
    },
    created () {
      this.cloneModel()
    },
    computed: {
      $m () {
        return this[prop]
      },
    },
    methods: {
      getter (key) {
        return (this.localModel || {})[key]
      },
      setter (key, value) {
        this.localModel[key] = value
      },
      syncModel () {
        this.$emit(event, this.localModel)
      },
      revertModel () {
        this.cloneModel()
      },
      cloneModel (model = this.$m) {
        this.localModel = cloneDeep(model)
      },
    },
    watch: {
      '$m' (newValue) {
        if (!isEqual(newValue, this.localModel)) {
          this.cloneModel(newValue)
        }
      },
    },
  }
}
