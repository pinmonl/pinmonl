<template>
  <div :class="$style.tagView">
    <Box v-if="tags.length > 0">
      <template v-for="tag in tags">
        <TagNode :tag="tag" :previsouParentName="safeParentName" :key="tag.id" :active="isActive(tag)">
          <template #before>
            <Anchor :to="{ name: 'tag', params: {id: tag.id} }" :replace="showPanel" inset />
          </template>
        </TagNode>
      </template>
    </Box>
    <Box v-else :class="$style.emptyResult">Empty.</Box>

    <Modal
      v-if="showPanel"
      @backdrop="handlePanelClose"
      @close="handlePanelClose"
    >
      <template #header="{ headerClass }" v-if="editing">
        <div :class="headerClass" v-text="`${isNew ? 'New' : 'Update'} Tag`" />
      </template>
      <TagDetail
        :tag="model"
        :loading="loading"
        :editable.sync="editing"
        @input="handlePanelSave"
        @cancel="isNew ? handlePanelClose() : null"
      >
        <template #controls="slotProps">
          <div :class="$style.panelButtons">
            <template v-if="!editing">
              <Button @click="handlePanelEdit">Edit</Button>
            </template>
            <template v-else>
              <Button :disabled="slotProps.error" @click="slotProps.submit">Save</Button>
              <Button @click="slotProps.cancel" light>Cancel</Button>
            </template>
          </div>
        </template>
      </TagDetail>
    </Modal>
  </div>
</template>

<script>
import { formatRepeatParam } from '@/pkgs/utils'
import Box from '@/components/app/Box.vue'
import Button from '@/components/form/Button.vue'
import Modal from '@/components/modal/Modal.vue'
import TagDetail from '@/components/tag/TagDetail.vue'
import TagNode from '@/components/tag/TagNode.vue'

export default {
  components: {
    Box,
    Button,
    Modal,
    TagDetail,
    TagNode,
  },
  props: {
    id: {
      type: String,
      default: null,
    },
    isNew: {
      type: Boolean,
      default: false,
    },
    parentName: {
      type: [String, Array],
      default: null,
    },
  },
  data () {
    return {
      loading: false,
      editing: false,
      storedParentName: [],

      highlight: [],
      cursor: -1,

      model: null,
      original: null,
    }
  },
  created () {
    this.initModel()
  },
  mounted () {
    document.addEventListener('keyup', this.handleKeyPress)
  },
  beforeDestroy () {
    document.removeEventListener('keyup', this.handleKeyPress)
  },
  computed: {
    showPanel () {
      return this.isNew || this.hasId
    },
    hasId () {
      return !!this.id
    },
    tags () {
      const tags = this.$store.getters['tag/tags']
      return this.$store.getters['tag/getByParent'](tags, this.parentId)
    },
    safeParentName () {
      return this.parentName
        ? formatRepeatParam(this.parentName)
        : this.storedParentName
    },
    parentId () {
      return this.parent ? this.parent.id : ''
    },
    parent () {
      const len = this.parents.length
      if (!len) {
        return null
      }
      return this.parents[len - 1]
    },
    parents () {
      const tags = this.$store.getters['tag/tags']
      return this.$store.getters['tag/getByName'](tags, this.safeParentName)
    },
    disableKeys () {
      return this.$store.getters.globalSearch
    },
  },
  methods: {
    async initModel () {
      this.editing = this.isNew
      this.loading = true
      if (this.hasId) {
        this.original = await this.find(this.id)
      } else if (this.isNew) {
        this.original = this.$store.getters['tag/new']()
      }
      this.model = this.original
      this.loading = false
    },
    async find(id) {
      return await this.$store.dispatch('tag/find', { id })
    },
    handlePanelClose () {
      let to = { name: 'tag.list' }
      if (this.showPanel && this.storedParentName.length) {
        to = { name: 'tag.children', params: {parentName: this.storedParentName} }
      }
      this.$router.push(to)
    },
    handlePanelEdit () {
      this.editing = true
    },
    handlePanelCancel () {
      this.editing = false
      this.model = this.original
    },
    async handlePanelSave (newModel) {
      const model = await this.save(newModel)
      this.model = this.original = model
      this.editing = false
    },
    async save (model) {
      if (model.id) {
        return await this.$store.dispatch('tag/update', model)
      } else {
        const tag = await this.$store.dispatch('tag/create', model)
        this.$router.replace({ name: 'tag', params: {id: tag.id} })
        return tag
      }
    },
    isActive ({ id }) {
      if (this.hasId) {
        return this.id == id
      }
      return this.highlight.includes(id)
    },
    highlightAt (n) {
      if (n >= this.tags.length) {
        return
      }
      this.cursor = n
      const { id } = this.tags[n]
      this.highlight = [id]
    },
    highlightDown () {
      if (this.cursor + 1 >= this.tags.length) {
        return
      }
      return this.highlightAt(this.cursor + 1)
    },
    highlightUp () {
      if (this.cursor - 1 < 0) {
        return
      }
      return this.highlightAt(this.cursor - 1)
    },
    gotoBookmark () {
      if (typeof this.tags[this.cursor] == 'undefined') {
        return
      }
      const tag = this.tags[this.cursor]
      this.$router.push({
        name: 'bookmark.list', 
        query: {
          q: this.$store.getters['pinl/composeSearch']({
            input: '',
            tags: [tag],
          }),
        },
      })
    },
    gotoDetail () {
      if (typeof this.tags[this.cursor] == 'undefined') {
        return
      }
      const { id } = this.tags[this.cursor]
      this.$router.push({
        name: 'tag',
        params: { id },
        query: this.$route.query,
      })
    },
    gotoChildren () {
      if (typeof this.tags[this.cursor] == 'undefined') {
        return
      }
      const { name } = this.tags[this.cursor]
      this.$router.push({
        name: 'tag.children',
        params: {
          parentName: [ ...this.safeParentName, name ],
        },
        query: this.$route.query,
      })
      this.resetHighlightAndCursor()
    },
    gotoParent () {
      const len = this.safeParentName.length
      if (len == 0) {
        return
      }
      if (len == 1) {
        this.$router.push({
          name: 'tag.list',
          query: this.$route.query,
        })
      } else {
        this.$router.push({
          name: 'tag.children',
          params: {
            parentName: [ ...this.safeParentName.slice(0, len - 1) ],
          },
          query: this.$route.query,
        })
      }
      this.resetHighlightAndCursor()
    },
    resetHighlightAndCursor () {
      this.highlight = []
      this.cursor = -1
    },
    handleKeyPress (e) {
      if (this.disableKeys) {
        return
      }
      if (this.hasId || this.isNew) {
        return
      }

      if (e.key == 'a') {
        this.$router.push({ name: 'tag.new' })
        return
      }
      if (e.key == 'j') {
        this.highlightDown()
        return
      }
      if (e.key == 'k') {
        this.highlightUp()
        return
      }
      if (e.key == 'e') {
        this.gotoDetail()
        return
      }
      if (e.key == 'o') {
        this.gotoBookmark()
        return
      }
      if (e.key == 'l') {
        this.gotoChildren()
        return
      }
      if (e.key == 'h') {
        this.gotoParent()
        return
      }
    },
  },
  watch: {
    parentName (newValue, oldValue) {
      if (this.id || this.isNew) {
        this.storedParentName = formatRepeatParam(oldValue)
      } else {
        this.storedParentName = []
      }
    },
    id () {
      this.initModel()
    },
    isNew () {
      this.initModel()
    },
  },
  metaInfo () {
    let title = 'Tag'
    if (this.hasId && this.model) {
      title = `Tag #${this.model.name}`
    }
    if (this.isNew) {
      title = `New ${title}`
    }
    return { title }
  },
}
</script>

<style lang="scss" module>
.tagView {
  @apply leading-normal;
}

.emptyResult {
  @apply p-4;
  @apply text-xs;
}

.panelButtons {
  @apply flex;
  @apply justify-end;
  @apply text-xs;

  > button {
    @apply m-1;
  }
}
</style>
