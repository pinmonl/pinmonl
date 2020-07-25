<script>
import pageMixin from '@/mixins/page'
import ShareListItem from '@/components/share/ShareListItem'
import AppContainer from '@/components/app/AppContainer'
import { share as api } from '@/api'
import ShareEditor from '@/components/share/ShareEditor'

export default {
  name: 'share-view',
  mixins: [pageMixin],
  data: () => ({
    items: [],
    loading: false,
    editorValue: null,
  }),
  created () {
    this.fetchData()
  },
  computed: {
    showEditor: {
      get () { return !!this.editorValue },
      set (v) { this.editorValue = v ? v : null },
    },
  },
  methods: {
    async fetchData () {
      try {
        this.loading = true

        const query = { page: this.page }
        const res = await api.list(query)
        this.items = res.data
        this.updateTotalPage(res)
      } catch (e) {
        //
      } finally {
        this.loading = false
      }
    },
    renderList () {
      return this.$createElement(ShareListItem, null, [
      ])
    },
    renderEditor () {
      return this.$createElement(ShareEditor, {
        attrs: { show: true, autofocus: true },
        on: {
          'update:show': (v) => this.showEditor = v,
        },
      })
    },
  },
  render (h) {
    const children = [
      this.renderList(),
      this.renderEditor(),
    ]

    return h(AppContainer, null, children)
  },
  watch: {
  },
}
</script>
