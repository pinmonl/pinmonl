<script>
import Container from '@/components/app/Container.vue'
import Search from '@/components/app/Search.vue'

export default {
  render (h) {
    return h(Container, {
      scopedSlots: {
        default: () => h('router-view'),
        header: () => this.renderHeader(h),
      },
    })
  },
  data () {
    return {
      storedSearch: '',
    }
  },
  computed: {
    search: {
      get () {
        return this.$route.query.q || ''
      },
      set (q) {
        if (q == this.search) {
          return
        }
        this.$router.push({ name: 'bookmark.list', query: {q} })
      },
    },
  },
  methods: {
    renderHeader (h) {
      return h(Search, {
        props: {
          search: this.search,
          disableKeys: this.$route.meta.noSearch,
        },
        on: {
          input: q => this.search = q,
        },
      })
    },
  },
}
</script>
