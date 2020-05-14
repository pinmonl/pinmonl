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
  mounted () {
    document.addEventListener('keyup', this.handleKeyPress)
  },
  beforeDestroy () {
    document.removeEventListener('keyup', this.handleKeyPress)
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
    disableKeys () {
      return this.$store.getters.globalSearch
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
          input: (q) => {
            this.search = q
            this.$refs.search.blur()
          },
          focus: () => this.$store.commit('SET_GLOBAL_SEARCH', true),
          blur: () => this.$store.commit('SET_GLOBAL_SEARCH', false),
        },
        ref: 'search',
      })
    },
    handleKeyPress (e) {
      if (this.disableKeys) {
        return
      }

      if (e.key == '1' && this.$route.name != 'bookmark.list') {
        this.$router.push({ name: 'bookmark.list' })
        return
      }
      if (e.key == '2' && this.$route.name != 'tag.list') {
        this.$router.push({ name: 'tag.list' })
        return
      }
    },
  },
}
</script>
