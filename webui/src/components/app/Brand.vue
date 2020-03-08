<script>
export default {
  props: {
    tag: {
      type: [Object, String],
      default: 'h3',
    },
    to: {
      type: [Object, String],
      default: null,
    },
  },
  render (h) {
    return this.renderRoot(h, [
      this.renderLogo(h),
      this.renderName(h),
    ])
  },
  methods: {
    renderRoot (h, children) {
      let elem = h(this.tag, {
        class: this.$style.container,
      }, children)
      if (this.to !== null) {
        elem = h('Anchor', {
          attrs: { to: this.to },
        }, [elem])
      }
      return elem
    },
    renderLogo (h) {
      return h('div', {
        class: this.$style.logo,
      }, [
        h('img', {
          attrs: { src: require('@/assets/logo.svg') },
        }),
      ])
    },
    renderName (h) {
      return h('div', 'Pinmonl')
    },
  },
}
</script>

<style lang="scss" module>
$logo-size: 2rem;

.container {
  @apply flex;
  @apply items-center;
  @apply font-bold;
}

.logo {
  width: $logo-size;
  height: $logo-size;
  @apply mr-2;
}
</style>
