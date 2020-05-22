<script>
import Footer from '@/components/app/Footer.vue'
import Nav from '@/components/app/Nav.vue'
import ThreeColumn from '@/components/layout/ThreeColumn.vue'
import scrollable from '@/provides/scrollable'

export default {
  provide () {
    return {
      [scrollable.name]: () => this.$el,
    }
  },
  render (h) {
    return h(ThreeColumn, {
      scopedSlots: {
        default: () => h('router-view'),
        left: () => this.renderLeft(h),
        right: () => h(Footer),
      },
    })
  },
  beforeRouteEnter (to, from, next) {
    next(vm => {
      vm.$el.scrollTo({ top: 0, left: 0 })
    })
  },
  computed: {
    scrollable () {
      return this.$el
    },
    navShowNewBookmark () {
      return this.$route.matched.filter(m => m.meta.navShowNewBookmark).length > 0
    },
    navShowNewTag () {
      return this.$route.matched.filter(m => m.meta.navShowNewTag).length > 0
    },
  },
  methods: {
    renderLeft (h) {
      return h(Nav, {
        scopedSlots: {
          controls: (props) => ([
            this.navShowNewBookmark && this.renderControl(
              h, props, { name: 'bookmark.new', query: this.$route.query }, '+ Bookmark'
            ),
            this.navShowNewTag && this.renderControl(
              h, props, { name: 'tag.new', query: this.$route.query }, '+ Tag'
            ),
          ]),
        }
      })
    },
    renderControl (h, props, to, label) {
      return h('Anchor', {
        class: props.anchorClass,
        attrs: { to },
      }, [
        h('div', {
          class: props.labelClass,
        }, label),
      ])
    }
  },
}
</script>
