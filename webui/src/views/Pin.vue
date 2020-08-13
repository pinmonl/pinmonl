<template>
  <div>
    <v-list two-line>
      <template v-if="loading">
        <div v-for="n in 5" :key="n">
          <v-divider v-if="n == 1"></v-divider>
          <v-skeleton-loader
            type="list-item-avatar-two-line"
          ></v-skeleton-loader>
          <v-divider></v-divider>
        </div>
      </template>

      <template v-else>
        <template v-for="(item, n) in items">
          <div :key="item.id">
            <v-divider v-if="n == 0"></v-divider>
            <pin-list-item
              :pin="item"
              @edit="edit(item)"
            ></pin-list-item>
            <v-divider></v-divider>
          </div>
        </template>
        <template v-if="items.length == 0">
          <v-divider></v-divider>
          <v-list-item>
            <v-list-item-title>
              No result.
            </v-list-item-title>
          </v-list-item>
          <v-divider></v-divider>
        </template>
      </template>
    </v-list>

    <pin-editor
      :show.sync="showEditor"
      :value="editorValue"
      autofocus
      @change="showEditor = false"
      @remove="showEditor = false"
    ></pin-editor>

    <v-speed-dial
      v-model="fab"
      bottom
      right
      fixed
      direction="top"
    >
      <template #activator>
        <v-btn color="primary" fab>
          <v-icon v-if="fab">{{ mdiClose }}</v-icon>
          <v-icon v-else>{{ mdiBookmarkMultiple }}</v-icon>
        </v-btn>
      </template>

      <v-btn
        fab
        dark
        small
        color="green"
        @click="create"
        title="Create"
      >
        <v-icon>{{ mdiPlus }}</v-icon>
      </v-btn>
    </v-speed-dial>
  </div>
</template>

<script>
import PinListItem from '@/components/pin/PinListItem'
import PinEditor from '@/components/pin/PinEditor'
import cloneDeep from 'lodash.clonedeep'
import { pin as pinDefault } from '@/utils/model'
import { mapState } from 'vuex'
import { Topics } from '@/utils/constants'
import {
  mdiPlus,
  mdiBookmarkMultiple,
  mdiClose,
} from '@mdi/js'

export default {
  name: 'pin-view',
  components: {
    PinListItem,
    PinEditor,
  },
  data () {
    return {
      items: [],
      loading: true,
      editorValue: null,
      fab: false,

      mdiBookmarkMultiple,
      mdiPlus,
      mdiClose,
    }
  },
  mounted () {
    this.fetchData()
    this.$store.state.socket.on(Topics.PINL_UPDATED, this.onPinlUpdated)
    this.$store.state.socket.on(Topics.PINL_DELETED, this.onPinlDeleted)
  },
  beforeDestroy () {
    this.$store.state.socket.off(Topics.PINL_UPDATED, this.onPinlUpdated)
    this.$store.state.socket.off(Topics.PINL_DELETED, this.onPinlDeleted)
  },
  computed: {
    ...mapState(['client', 'search']),
    showEditor: {
      get () {
        return !!this.editorValue
      },
      set (val) {
        this.editorValue = val ? val : null
      },
    },
  },
  methods: {
    async fetchData () {
      try {
        this.loading = true

        const query = { page_size: 0 }
        if (this.search.isEmpty()) {
          query.notag = 1
        } else {
          query.q = this.search.getText()
          query.tag = this.search.getValues('tag').join(',')
        }

        const res = await this.client.listPins(query)
        this.items = res.data
      } catch (e) {
        //
      } finally {
        this.loading = false
      }
    },
    create () {
      this.edit(pinDefault)
    },
    edit (item) {
      this.editorValue = cloneDeep(item)
    },
    findIndex (item) {
      return this.items.findIndex(p => p.id == item.id)
    },
    update (item) {
      const idx = this.findIndex(item)

      // Replace if found.
      if (idx >= 0) {
        this.items = [
          ...this.items.slice(0, idx),
          { ...this.items[idx], ...item },
          ...this.items.slice(idx + 1),
        ]
        return
      }

      // Compare tags with query.
      if (this.search.isEmpty() && item.tags.length > 0) {
        return
      }

      const searchTags = this.search.getValues('tag')
      const remainings = searchTags.filter(tag => !item.tags.includes(tag.replace(/^\/+|\/+$/, '')))
      if (remainings.length > 0) {
        return
      }

      this.items = [ item, ...this.items ]
    },
    remove (item) {
      const idx = this.findIndex(item)
      if (idx < 0) {
        return
      }

      // Remove if found.
      this.items = [
        ...this.items.slice(0, idx),
        ...this.items.slice(idx + 1),
      ]
    },
    onPinlUpdated ({ data }) {
      this.update(data)
    },
    onPinlDeleted ({ data }) {
      this.remove(data)
    },
  },
  watch: {
    page () {
      this.fetchData()
    },
    search () {
      this.fetchData()
    },
  },
}
</script>
