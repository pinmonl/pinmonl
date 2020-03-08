<template>
  <div :class="$style.container">
    <Header :class="$style.header">
      <h2 :class="$style.title">Bookmark</h2>
      <div :class="$style.search">
        <Search v-model="search" />
      </div>
      <Anchor :class="$style.addBtn" :to="{ name: 'bookmark.new' }">
        <IconButton name="add" block />
      </Anchor>
    </Header>
    <div :class="$style.list">
      <template v-for="pinl in pinls">
        <div :key="pinl.id">
          <Pinl :pinl="pinl">
            <template #before>
              <Anchor :to="{ name: 'bookmark', params: {id: pinl.id} }" :replace="showPanel" inset />
            </template>
          </Pinl>
          <Divider />
        </div>
      </template>
    </div>
    <RightPanel
      v-if="showPanel"
      @close="handlePanelClose"
      @edit="handlePanelEdit"
      @save="handlePanelSave"
      @cancel="handlePanelCancel"
      :noSave="!editing"
      :noCancel="!editing || $props.new"
      :noEdit="editing"
    >
      <PinlDetail
        v-model="model"
        :class="$style.detail"
        :loading="loading"
        :editable="editing"
      />
    </RightPanel>
  </div>
</template>

<script>
import Header from '@/components/app/Header.vue'
import IconButton from '@/components/form/IconButton.vue'
import Pinl from '@/components/pinl/Pinl.vue'
import PinlDetail from '@/components/pinl/PinlDetail.vue'
import RightPanel from '@/components/modal/RightPanel.vue'
import Search from '@/components/app/Search.vue'

export default {
  components: {
    Header,
    IconButton,
    Pinl,
    PinlDetail,
    RightPanel,
    Search,
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
  },
  data () {
    return {
      loading: false,
      editing: false,
      storedSearch: '',

      model: null,
      original: null,
    }
  },
  created () {
    this.initModel()
    this.storedSearch = this.search
  },
  computed: {
    hasId () {
      return !!this.id
    },
    showPanel () {
      return this.new || this.hasId
    },
    pinls () {
      const { input, tags } = this.$store.getters['pinl/parseSearch'](this.search)
      const tagNames = this.$store.getters['tag/mapName'](tags)
      let pinls = this.$store.getters['pinl/pinls']
      pinls = this.$store.getters['pinl/getByTag'](pinls, tagNames)
      pinls = this.$store.getters['pinl/searchByTitle'](pinls, input)
      return pinls
    },
    search: {
      get () {
        return this.$route.query.q || this.storedSearch
      },
      set (q) {
        if (q == this.search) {
          return
        }
        this.storedSearch = q
        this.$router.replace({ query: {q} })
      },
    },
  },
  methods: {
    handlePanelClose () {
      this.$router.push({ name: 'bookmark' })
    },
    async initModel () {
      this.editing = this.new
      this.loading = true
      if (this.hasId) {
        this.original = await this.find(this.id)
      } else if (this.new) {
        this.original = this.$store.getters['pinl/new']()
      }
      this.model = this.original
      this.loading = false
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
    async find (id) {
      return await this.$store.dispatch('pinl/find', { id })
    },
    async save (model) {
      if (model.id) {
        return await this.$store.dispatch('pinl/update', model)
      } else {
        const pinl = await this.$store.dispatch('pinl/create', model)
        this.$router.replace({ name: 'bookmark', params: {id: pinl.id} })
        return pinl
      }
    },
  },
  watch: {
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
  @apply h-full;
  @apply w-full;
  @apply overflow-y-auto;
  @apply scrolling-touch;
}

.header {
}

.title {
  @apply flex-grow;
}

.search {
  @apply relative;
}

.addBtn {
}

.list {
}

.pinl {
  @apply border-b;
  @apply border-solid;
  border-color: theme('colors.divider');
}

.detail {
  @apply h-full;
  @apply p-4;
  @apply overflow-y-auto;
  @apply scrolling-touch;
}
</style>
