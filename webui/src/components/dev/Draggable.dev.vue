<script>
export default {
  props: {
    list: {
      type: Array,
      default: () => ([]),
    },
    tag: {
      type: [String, Object],
      default: 'div',
    },
  },
  render (h) {
    return h(this.tag, {
      class: 'draggable',
    }, this.children)
  },
  computed: {
    children () {
      const slot = this.$scopedSlots.default
      if (!slot) {
        return []
      }
      return this.list.map(item => {
        const vn = slot({ item })
        console.log(vn)
        vn[0].context.$el.addEventListener('mouseover', () => console.log(1))
        return vn
      })
    },
  },
}
</script>
