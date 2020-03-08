<template>
  <div
    :class="$style.editable"
    contenteditable
    v-on="listeners"
    ref="editable"
  />
</template>

<script>
export default {
  model: {
    prop: 'content',
    event: 'input',
  },
  props: {
    content: {
      type: String,
      default: '',
    },
  },
  data () {
    return {
      value: '',
      caret: null,
    }
  },
  mounted () {
    this.copyContent(this.content)
  },
  computed: {
    listeners () {
      return {
        ...this.$listeners,
        input: this.handleInput,
        focus: this.handleFocus,
        blur: this.handleBlur,
      }
    },
  },
  methods: {
    copyContent (content) {
      this.value = content
      this.setContent(content)
    },
    handleInput (event) {
      const content = event.target.innerText
      this.value = content
      this.$emit('input', content)
    },
    handleFocus (event) {
      this.restoreCaret()
      this.$emit('focus', event)
    },
    handleBlur (event) {
      this.saveCaret()
      this.$emit('blur', event)
    },
    setContent (content) {
      this.$refs.editable.innerText = content
    },
    focus () {
      this.$refs.editable.focus()
    },
    blur () {
      this.$refs.editable.blur()
    },
    saveCaret () {
      const sel = window.getSelection()
      if (!this.$el.contains(sel.focusNode)) {
        return
      }
      this.caret = [ sel.focusNode, sel.focusOffset ]
    },
    restoreCaret () {
      if (this.caret == null) {
        return
      }
      const sel = window.getSelection()
      const [ node ] = this.caret
      console.log([node, sel.focusNode])
      // sel.collapse(node, offset)
      this.clearCaret()
    },
    clearCaret () {
      this.caret = null
    },
  },
  watch: {
    content (content) {
      if (this.value != content) {
        this.saveCaret()
        this.copyContent(content)
        this.restoreCaret()
      }
    },
  },
}
</script>

<style lang="scss" module>
.editable {
  @apply text-sm;
  @apply whitespace-pre;
}
</style>