<script>
export default {
  props: {
    current: {
      type: String,
      required: true,
    },
    previous: {
      type: Array,
      default: null,
    },
    noHome: {
      type: Boolean,
      default: false,
    },
  },
  render (h) {
    return h('div', {
      class: this.$style.breadcrumb,
    }, this.items.map(item => this.renderItem(h, item)))
  },
  computed: {
    items () {
      let items = []
      if (!this.noHome) {
        items.push({ to: '/', label: 'Home' })
      }
      if (this.previous !== null) {
        items = [ ...this.items, ...this.previous ]
      }
      items.push({ current: true, label: this.current })
      return items
    },
  },
  methods: {
    renderItem (h, item) {
      let el
      if (item.to) {
        el = h('Anchor', {
          attrs: { to: item.to },
          props: { color: true },
        }, item.label)
      } else {
        el = h('span', null, item.label)
      }

      return h('div', {
        class: this.$style.item,
      }, [el])
    },
  },
}
</script>

<style lang="scss" module>
.breadcrumb {
  @apply m-box;
  @apply my-6;
  @apply text-sm;
  @apply flex;
  @apply flex-wrap;
}

.item {
  @apply relative;

  &:not(:last-child) {
    &::after {
      content: '/';
      @apply absolute;
      @apply inset-y-0;
      right: -2px;
    }
  }

  > a, span {
    @apply p-2;
  }
}
</style>
