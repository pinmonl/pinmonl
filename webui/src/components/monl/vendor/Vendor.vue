<script>
export default {
  props: {
    name: {
      type: String,
    },
    url: {
      type: String,
    },
  },
  render (h) {
    let iconEl = h('Icon', {
      class: this.$style.icon, 
      props: { name: this.iconName }
    })
    if (this.url) {
      iconEl = h('Anchor', {
        props: { external: true },
        attrs: { to: this.url },
      }, [iconEl])
    }

    return h('div', {
      class: this.$style.vendor,
    }, [
      iconEl,
      h('div', {
        class: this.$style.stats,
      }, this.$slots.default),
    ])
  },
  computed: {
    iconName () {
      const map = {
        'github': 'github',
      }
      return map[this.name]
    },
  },
}
</script>

<style lang="scss" module>
.vendor {
  @apply flex;
}

.icon {
}

.stats {
  @apply flex;
  @apply items-center;
  @apply flex-1;
  @apply min-w-0;
  @apply overflow-x-auto;
  @apply ml-2;
  @apply flex-wrap;
}
</style>
