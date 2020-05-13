<script>
export default {
  inheritAttrs: false,
  props: {
    tags: {
      type: Array,
      default: () => ([]),
    },
  },
  render (h) {
    return h('div', {
      class: this.$style.list,
    }, this.renderTags(h, this.tags))
  },
  methods: {
    renderTags (h, tags) {
      let elems = []
      for (const i in tags) {
        if (i > 0) {
          elems.push(this.renderChevron(h))
        }
        elems.push(this.renderTag(h, tags[i]))
      }
      return elems
    },
    renderTag (h, tag) {
      return h('div', {
        class: this.$style.tag,
        props: { tag },
      }, tag.name)
    },
    renderFolder (h) {
      return h('div', {
        class: this.$style.icon,
      }, [
        h('Icon', {
          props: { name: 'folder' },
        }),
      ])
    },
    renderChevron (h) {
      return h('div', {
        class: this.$style.icon,
      }, [
        h('Icon', {
          props: { name: 'chevronRight' },
        }),
      ])
    },
  },
}
</script>

<style lang="scss" module>
.list {
  @apply flex;
  @apply flex-wrap;
  @apply items-center;
  @apply text-sm;
}

.tag {
  @apply m-1;
}

.icon {
  @apply -mx-1;
}
</style>
