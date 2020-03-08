<template>
  <div :class="$style.container">
    <Header :class="$style.header">
      <h2 :class="$style.title">Tag</h2>
      <Anchor :to="{ name: 'tag.new' }" :class="$style.addBtn">
        <IconButton name="add" block />
      </Anchor>
    </Header>
    <div :class="$style.list">
      <template v-for="tag in tags">
        <div :key="tag.id">
          <TagNode :tag="tag" :previsouParentName="safeParentName">
            <template #before>
              <Anchor :to="{ name: 'tag', params: {id: tag.id} }" :replace="showPanel" inset />
            </template>
          </TagNode>
        </div>
      </template>
    </div>
    <RightPanel
      v-if="showPanel"
      @close="handlePanelClose"
      @save="handlePanelSave"
      noCancel
      noEdit
    >
      <TagDetail
        v-model="model"
        :loading="loading"
        :editable="!loading"
      />
    </RightPanel>
  </div>
</template>

<script>
import { formatRepeatParam } from '@/pkgs/utils'
import Header from '@/components/app/Header.vue'
import IconButton from '@/components/form/IconButton.vue'
import RightPanel from '@/components/modal/RightPanel.vue'
import TagDetail from '@/components/tag/TagDetail.vue'
import TagNode from '@/components/tag/TagNode.vue'

export default {
  components: {
    Header,
    IconButton,
    RightPanel,
    TagDetail,
    TagNode,
  },
  props: {
    id: {
      type: String,
      default: null,
    },
    new: {
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

      model: null,
      original: null,
    }
  },
  created () {
    this.initModel()
  },
  computed: {
    showPanel () {
      return this.new || this.hasId
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
  },
  methods: {
    async initModel () {
      this.editing = this.new
      this.loading = true
      if (this.hasId) {
        this.original = await this.find(this.id)
      } else if (this.new) {
        this.original = this.$store.getters['tag/new']()
      }
      this.model = this.original
      this.loading = false
    },
    async find(id) {
      const tags = this.$store.getters['tag/tags']
      return this.$store.getters['tag/find'](tags, id)
    },
    handlePanelClose () {
      let to = { name: 'tag' }
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
    async handlePanelSave () {
      const model = await this.save(this.model)
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
  },
  watch: {
    parentName (newValue, oldValue) {
      if (this.id || this.new) {
        this.storedParentName = formatRepeatParam(oldValue)
      } else {
        this.storedParentName = []
      }
    },
    id () {
      this.initModel()
    },
    new () {
      this.initModel()
    },
  },
}
</script>

<style lang="scss" module>
.container {
  @apply relative;
  @apply w-full;
  @apply h-full;
  @apply overflow-y-auto;
  @apply scrolling-touch;
}

.header {
  @apply flex;
}

.title {
  @apply flex-grow;
}

.addBtn {
}
</style>
