<template>
  <div>
    <div v-if="loading" :class="$style.skeleton">
      <v-skeleton-loader
        v-for="n in 5"
        :key="n"
        type="list-item"
      ></v-skeleton-loader>
    </div>

    <tag-tree
      v-show="!loading"
      :search="searchInput"
      editable
      ref="tree"
      @edit="edit"
      autofocus
    ></tag-tree>

    <v-list v-if="isEmpty">
      <v-divider></v-divider>
      <v-list-item two-line>
        <v-list-item-title>
          It's empty! <a href="#" @click.prevent="create">Click here to create one</a>.
        </v-list-item-title>
      </v-list-item>
      <v-divider></v-divider>
    </v-list>

    <tag-editor
      :show.sync="showEditor"
      :value="editorValue"
      @change="reload"
      @delete="reload"
    ></tag-editor>

    <v-speed-dial
      v-model="fab"
      bottom
      right
      fixed
      direction="top"
    >
      <template #activator>
        <v-btn
          color="primary"
          fab
        >
          <v-icon v-if="!fab">{{ mdiTagMultiple }}</v-icon>
          <v-icon v-else>{{ mdiClose }}</v-icon>
        </v-btn>
      </template>

      <v-btn
        color="green"
        dark
        small
        fab
        @click="create"
        title="Create"
      >
        <v-icon>{{ mdiPlus }}</v-icon>
      </v-btn>

      <v-btn
        fab
        @click="expandAll"
        small
        dark
        color="indigo"
        v-if="!isExpanded"
        title="Expand all"
      >
        <v-icon>{{ mdiExpandAll }}</v-icon>
      </v-btn>
      <v-btn
        fab
        @click="collapseAll"
        small
        dark
        color="pink"
        v-if="isExpanded"
        title="Collapse all"
      >
        <v-icon>{{ mdiCollapseAll }}</v-icon>
      </v-btn>
    </v-speed-dial>
  </div>
</template>

<script>
import cloneDeep from 'lodash.clonedeep'
import TagTree from '@/components/tag/TagTree'
import TagEditor from '@/components/tag/TagEditor'
import {
  mdiMagnify,
  mdiTagMultiple,
  mdiExpandAll,
  mdiCollapseAll,
  mdiClose,
  mdiPlus,
} from '@mdi/js'
import { tag as tagDefault } from '@/utils/model'

export default {
  name: 'tag-view',
  components: {
    TagTree,
    TagEditor,
  },
  data () {
    return {
      searchInput: '',
      editorValue: null,
      isExpanded: false,
      fab: false,
      bootstrapped: false,
      isMounted: false,

      mdiMagnify,
      mdiTagMultiple,
      mdiExpandAll,
      mdiCollapseAll,
      mdiClose,
      mdiPlus,
    }
  },
  mounted () {
    this.isMounted = true
  },
  computed: {
    showEditor: {
      get () {
        return !!this.editorValue
      },
      set () {
        this.editorValue = null
      },
    },
    loading () {
      if (!this.isMounted) {
        return true
      }
      if (!this.$refs.tree.inited) {
        return this.$refs.tree.loading
      }
      return false
    },
    isEmpty () {
      if (!this.isMounted) {
        return false
      }
      return this.$refs.tree.inited && this.$refs.tree.isEmpty
    },
  },
  methods: {
    edit (item) {
      this.editorValue = cloneDeep(item)
    },
    create () {
      this.edit(tagDefault)
    },
    async reload () {
      await this.$refs.tree.reload()
      this.showEditor = false
    },
    expandAll () {
      this.isExpanded = true
      this.$refs.tree.openAll()
    },
    collapseAll () {
      this.isExpanded = false
      this.$refs.tree.closeAll()
    },
  },
}
</script>

<style lang="scss" module>
.skeleton {
  :global(.v-skeleton-loader) {
    &:nth-child(1) {
      :global(.v-skeleton-loader__list-item) {
        max-width: 200px;
      }
    }

    &:nth-child(2) {
      :global(.v-skeleton-loader__list-item) {
        max-width: 240px;
      }
    }

    &:nth-child(3) {
      :global(.v-skeleton-loader__list-item) {
        max-width: 160px;
      }
    }

    &:nth-child(4) {
      :global(.v-skeleton-loader__list-item) {
        max-width: 200px;
      }
    }

    &:nth-child(5) {
      :global(.v-skeleton-loader__list-item) {
        max-width: 180px;
      }
    }
  }
}
</style>
