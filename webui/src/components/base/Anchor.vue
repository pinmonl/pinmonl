<script>
export default {
  inheritAttrs: false,
  props: {
    external: {
      type: Boolean,
      default: false,
    },
    underline: {
      type: Boolean,
      default: false,
    },
    inset: {
      type: Boolean,
      default: false,
    },
    color: {
      type: Boolean,
      default: false,
    },
  },
  render (h) {
    const c = this.external ? 'a' : 'router-link'
    return h(c, {
      class: [this.$style.anchor, {
        [this.$style.underline]: this.underline,
        [this.$style.inset]: this.inset,
        [this.$style.color]: this.color,
      }],
      attrs: this.attrs,
      on: this.listeners,
    }, this.$slots.default)
  },
  computed: {
    attrs () {
      const attrs = { ...this.$attrs }
      if (this.external) {
        const to = attrs.to
        attrs.href = to || attrs.href
        attrs.target = '_blank'
        delete attrs.to
      }
      return attrs
    },
    listeners () {
      return {
        ...this.$listeners,
      }
    },
  },
}
</script>

<style lang="scss" module>
.anchor {
}

.underline {
  &:hover {
    @apply underline;
  }
}

.color {
  @apply text-anchor;
}

.inset {
  @apply absolute;
  @apply inset-0;
}
</style>
