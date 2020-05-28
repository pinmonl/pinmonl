<template>
  <div :class="$style.bookmarkView">
    <Box v-if="pinls.length > 0" :key="listKey">
      <template v-for="pinl in pinls">
        <Pinl :pinl="pinl" :key="pinl.id" :class="$style.pinl" :active="isActive(pinl)" ref="pinls">
          <template #before>
            <Anchor :to="{ name: 'bookmark', params: {id: pinl.id}, query: $route.query }" :replace="showPanel" inset />
          </template>
        </Pinl>
      </template>
    </Box>
    <Box v-else :class="$style.emptyResult">Empty.</Box>

    <Modal
      v-if="showPanel"
      @backdrop="handlePanelClose"
      @close="handlePanelClose"
    >
      <template #header="{ headerClass }" v-if="editing">
        <div :class="headerClass" v-text="`${isNew ? 'New' : 'Update'} Bookmark`" />
      </template>
      <PinlDetail
        :pinl="model"
        :loading="loading"
        :editable.sync="editing"
        @input="handlePanelSave"
        @cancel="isNew ? handlePanelClose() : null"
      >
        <template #controls="slotProps">
          <div :class="$style.panelButtons">
            <template v-if="!editing">
              <Button @click="handlePanelEdit" v-if="!deleting">Edit</Button>
              <ConfirmButton @click="handlePanelDelete" :asking.sync="askingDelete" :loading="deleting" light>Delete</ConfirmButton>
            </template>
            <template v-else>
              <Button :disabled="slotProps.error" @click="slotProps.submit" :loading="loading">Save</Button>
              <Button @click="slotProps.cancel" light>Cancel</Button>
            </template>
          </div>
        </template>
      </PinlDetail>
    </Modal>
  </div>
</template>

<script>
import Box from '@/components/app/Box.vue'
import Button from '@/components/form/Button.vue'
import ConfirmButton from '@/components/form/ConfirmButton.vue'
import Modal from '@/components/modal/Modal.vue'
import Pinl from '@/components/pinl/Pinl.vue'
import PinlDetail from '@/components/pinl/PinlDetail.vue'
import isEqual from 'lodash/isEqual'
import keybindingMixin from '@/mixins/keybinding'
import scrollable from '@/provides/scrollable'

export default {
  mixins: [keybindingMixin()],
  components: {
    Box,
    Button,
    ConfirmButton,
    Modal,
    Pinl,
    PinlDetail,
  },
  inject: {
    scrollable: {
      from: scrollable.name,
    },
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
  },
  data () {
    return {
      loading: false,
      editing: false,
      storedSearch: '',
      askingDelete: false,
      deleting: false,

      highlight: [],
      cursor: -1,

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
      return this.isNew || this.hasId
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
    listKey () {
      return JSON.stringify(this.$route.query)
    },
  },
  methods: {
    handlePanelClose () {
      if (this.deleting) {
        return
      }
      this.$router.push({ name: 'bookmark.list', query: this.$route.query })
    },
    async initModel () {
      this.editing = this.isNew
      this.loading = true
      if (this.hasId) {
        this.original = await this.find(this.id)
      } else if (this.isNew) {
        this.original = this.$store.getters['pinl/new']()
      }
      this.model = this.original
      this.loading = false
    },
    handlePanelEdit () {
      this.editing = true
    },
    handlePanelCancel () {
    },
    async handlePanelSave (newModel) {
      this.loading = true
      const model = await this.save(newModel)
      this.loading = false
      this.model = this.original = model
      this.editing = false
    },
    async handlePanelDelete () {
      this.deleting = true
      await this.delete(this.model)
      this.deleting = false
      // Back to listing
      this.$router.replace({ name: 'bookmark.list', query: this.$route.query })
      // Move cursor up
      this.highlightAt(Math.min(this.cursor, this.pinls.length - 1))
    },
    async find (id) {
      return await this.$store.dispatch('pinl/find', { id })
    },
    async save (model) {
      if (model.id) {
        return await this.$store.dispatch('pinl/update', model)
      } else {
        const pinl = await this.$store.dispatch('pinl/create', model)
        this.$router.replace({ name: 'bookmark', params: {id: pinl.id}, query: this.$route.query })
        return pinl
      }
    },
    async delete (model) {
      return await this.$store.dispatch('pinl/delete', model)
    },
    isActive ({ id }) {
      if (this.hasId) {
        return this.id == id
      }
      return this.highlight.includes(id)
    },
    highlightAt (n) {
      if (n >= this.pinls.length) {
        return
      }
      this.cursor = n
      const { id } = this.pinls[n]
      this.highlight = [id]
      this.scrollToPinl(n)
    },
    highlightDown () {
      if (this.cursor + 1 >= this.pinls.length) {
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
    scrollToPinl (n) {
      const pinlRef = this.$refs.pinls[n]
      if (!pinlRef) {
        return
      }

      const $pinl = pinlRef.$el
      const $parent = this.scrollable()
      const pinlRect = $pinl.getBoundingClientRect()
      const parentRect = $parent.getBoundingClientRect()

      const botDiff = pinlRect.bottom - parentRect.bottom
      const topDiff = pinlRect.top - parentRect.top - 1
      if (botDiff > 0) {
        $parent.scrollTo({ top: $parent.scrollTop + botDiff })
      } else if (topDiff < 0) {
        $parent.scrollTo({ top: $parent.scrollTop + topDiff })
      }
    },
    openHighlightLink () {
      if (typeof this.pinls[this.cursor] == 'undefined') {
        return
      }
      const pinl = this.pinls[this.cursor]
      this.$store.dispatch('pinl/openLink', { pinl })
    },
    gotoDetail () {
      if (typeof this.pinls[this.cursor] == 'undefined') {
        return
      }
      const { id } = this.pinls[this.cursor]
      this.$router.push({
        name: 'bookmark',
        params: { id },
        query: this.$route.query,
      })
    },
    handleKeyPress (e) {
      if (this.shouldDisableKeys) {
        return
      }
      if (this.editing) {
        return
      }
      if (this.hasId || this.isNew) {
        if (e.key == 'd') {
          if (!this.askingDelete) {
            this.askingDelete = true
          } else {
            this.askingDelete = false
            this.handlePanelDelete()
          }
          return
        }

        // Exit if not yet
        return
      }

      if (e.key == 'a') {
        this.$router.push({ name: 'bookmark.new', query: this.$route.query })
        e.preventDefault()
        return false
      }
      if (e.key == 'j') {
        this.highlightDown()
        return
      }
      if (e.key == 'k') {
        this.highlightUp()
        return
      }
      if (e.key == 'o') {
        this.openHighlightLink()
        return
      }
      if (e.key == 'e') {
        this.gotoDetail()
        return
      }
    },
  },
  watch: {
    id () {
      this.initModel()
    },
    isNew () {
      this.initModel()
    },
    '$route.query' (newValue, oldValue) {
      if (isEqual(newValue, oldValue)) {
        return
      }
      this.highlight = []
      this.cursor = -1
    },
  },
  metaInfo () {
    let title = 'Bookmark'
    if (this.hasId && this.model) {
      title = `${this.model.title}`
    }
    if (this.isNew) {
      title = `New ${title}`
    }
    return { title }
  },
}
</script>

<style lang="scss" module>
.bookmarkView {
  @apply leading-normal;
}

.emptyResult {
  @apply p-4;
  @apply text-xs;
}

.panelButtons {
  @apply flex;
  @apply justify-end;
  @apply text-sm;

  > button {
    @apply m-1;
  }
}
</style>
