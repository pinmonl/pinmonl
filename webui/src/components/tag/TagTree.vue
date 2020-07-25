<template>
  <v-treeview
    :items="items"
    :load-children="loadChildren"
    item-key="id"
    item-text="name"
    item-children="children"
    activatable
    :active.sync="active"
    :open.sync="open"
    v-bind="$attrs"
    v-on="$listeners"
  >
    <template #prepend="{ item, open }">
      <template v-if="!!item.children">
        <v-icon v-if="open">{{ mdiTagMultipleOutline }}</v-icon>
        <v-icon v-else>{{ mdiTagMultiple }}</v-icon>
      </template>
      <template v-else>
        <v-icon>{{ mdiTagOutline }}</v-icon>
      </template>
    </template>
    <template #append="{ item, active }">
      <v-scroll-x-reverse-transition hide-on-leave v-if="editable">
        <div v-if="active">
          <v-btn icon :to="showPins(item)">
            <v-icon>{{ mdiArrowRightCircle }}</v-icon>
          </v-btn>
          <v-btn @click.stop="edit(item)" icon>
            <v-icon>{{ mdiPencil }}</v-icon>
          </v-btn>
        </div>
      </v-scroll-x-reverse-transition>
    </template>
  </v-treeview>
</template>

<script>
import { mapState } from 'vuex'
import {
  mdiPencil,
  mdiTagMultiple,
  mdiTagMultipleOutline,
  mdiTagOutline,
  mdiArrowRightCircle,
} from '@mdi/js'
import { SearchParams } from '@/utils/search'

export default {
  name: 'tag-tree',
  props: {
    editable: Boolean,
  },
  data () {
    return {
      tags: [],
      open: [],
      active: [],
      loading: true,
      inited: false,

      mdiPencil,
      mdiTagMultiple,
      mdiTagMultipleOutline,
      mdiTagOutline,
      mdiArrowRightCircle,
    }
  },
  mounted () {
    this.loadAll()
  },
  computed: {
    ...mapState(['client']),
    items () {
      return this.getTagTree()
    },
    isEmpty () {
      return this.tags.length == 0
    },
  },
  methods: {
    edit (item) {
      this.$emit('edit', item.data)
    },
    async reload () {
      const open = [ ...this.open ]
      const active = [ ...this.active ]

      this.tags = []
      await this.loadAll()

      const tagIds = this.tags.map(tag => tag.id)
      this.open = open.filter(id => tagIds.includes(id))
      this.active = active.filter(id => tagIds.includes(id))
    },
    openAll () {
      this.open = this.tags.
        filter(tag => tag.hasChildren).
        map(tag => tag.id)
    },
    closeAll () {
      this.open = []
    },
    showPins (item) {
      const params = new SearchParams()
      params.add('tag', `/${item.data.name}`)

      return {
        path: '/pin',
        query: { q: params.encode() },
      }
    },
    async loadAll () {
      try {
        this.loading = true
        await this.fetchData()
      } catch (e) {
        //
      } finally {
        this.loading = false
      }
    },
    async loadRootTags () {
      try {
        this.loading = true
        await this.fetchData('')
      } catch (e) {
        //
      } finally {
        this.loading = false
      }
    },
    async loadChildren (parent) {
      try {
        await this.fetchData(parent.id)
      } catch (e) {
        //
      }
    },
    async fetchData (parentId) {
      const query = { page_size: 0 }
      if (typeof parentId == 'string') {
        query.parent = parentId
      }

      const res = await this.client.listTags(query)
      this.tags.unshift(...res.data)
      this.tags = this.tags.filter((tag, n) => this.findIndex(tag) == n)
    },
    findIndex (item) {
      return this.tags.findIndex(tag => tag.id == item.id)
    },
    getTagTree (parentId = '') {
      return this.getChildrenTags(parentId).map(tag => {
        const item = {
          id: tag.id,
          name: tag.name.replace(/^.*\//g, ''),
          data: tag,
        }

        if (tag.hasChildren) {
          item.children = this.getTagTree(tag.id)
        }

        return item
      })
    },
    getChildrenTags (parentId = '') {
      const tags = this.tags.filter(tag => tag.parentId == parentId)
      return tags || []
    },
  },
  watch: {
    loading (val) {
      if (!val && !this.inited) {
        this.inited = true
      }
    },
  },
}
</script>
