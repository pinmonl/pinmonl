<script>
import { VTextField } from 'vuetify/lib'

export default VTextField.extend({
  name: 'text-field',

  props: {
    labelActive: Boolean,
    inputClasses: [String, Array],
    hideInput: Boolean,
    inputContainerClasses: [String, Array],
  },

  computed: {
    isLabelActive () {
      return VTextField.options.computed.isLabelActive.call(this) || this.labelActive
    },
  },

  methods: {
    genInput () {
      const input = VTextField.options.methods.genInput.call(this)
      if (this.inputClasses) {
        input.data.class = [this.inputClasses]
      }
      if (this.hideInput) {
        input.data.class.push(this.$style.hideInput)
      }

      return this.$createElement('div', {
        class: [this.$style.inputContainer, this.inputContainerClasses],
      }, [
        this.$slots['before'],
        input,
      ])
    },
  },
});
</script>

<style lang="scss" module>
.inputContainer {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  padding-top: 2px;
  min-height: 32px;
  width: 100%;

  input {
    width: auto;
  }
}

input.hideInput {
  width: 0;
  max-width: 0;
}
</style>
