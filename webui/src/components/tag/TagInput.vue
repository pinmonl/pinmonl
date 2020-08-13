<template>
  <div class="py-1">
    <text-field
      ref="input"
      v-bind="$attrs"
      :label-active="isLabelActive"
      @focus="onFocus"
      v-model="inputValue"
      :input-classes="$style.input"
      @keydown.up.prevent.stop="onUp"
      @keydown.down.prevent.stop="onDown"
      @keydown.enter.prevent.stop="onEnter"
      @keydown.backspace="onBackspace"
      @keydown.tab.prevent.stop="onTab"
      @keydown.esc.prevent.stop="blur"
      :hide-input="!isFocused"
      :loading="loading"
    >
      <template #before>
        <tag-chip
          v-for="tag in selectedTags"
          :tag="tag"
          :key="tag"
        ></tag-chip>
      </template>
    </text-field>

    <v-menu
      :value="isMenuActive && menuCanShow"
      :close-on-click="false"
      :close-on-content-click="false"
      :open-on-click="false"
      :activator="$el"
      max-height="300"
      ref="menu"
      disable-keys
      offset-y
      offset-overflow
    >
      <v-list class="v-select-list" tabindex="-1" dense>
        <v-list-item
          v-for="(tag, n) in tags"
          :key="tag"
          @click="onItemClick(tag, n)"
          :class="getTagItemClasses(tag, n)"
        >
          <v-list-item-content>
            <v-list-item-title>
              <tag-chip :tag="tag"></tag-chip>
            </v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item v-if="isCreatable">
          <v-list-item-content>
            <v-list-item-title>
              Create <tag-chip :tag="stripInputValue"></tag-chip>
            </v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-menu>
  </div>
</template>

<script>
import { mapState } from 'vuex'
import debounce from 'lodash.debounce'
import TagChip from './TagChip'
import TextField from '@/components/form/TextField'

export default {
  name: 'tag-input',
  inheritAttrs: false,
  components: {
    TagChip,
    TextField,
  },
  props: {
    value: {
      type: [String, Array, Object],
    },
    multiple: {
      type: Boolean,
      default: false,
    },
  },
  data () {
    return {
      loading: false,
      isFocused: false,
      isMenuActive: false,
      selectIndex: -1,
      tags: [],
      internalSelectedTags: [],
      inputValue: '',
    }
  },
  mounted () {
    document.addEventListener('mousedown', this.onClickOutside)
  },
  beforeDestroy () {
    document.removeEventListener('mousedown', this.onClickOutside)
  },
  computed: {
    ...mapState(['client']),
    selectedTags: {
      get () {
        if (this.multiple) {
          return this.value || []
        }

        return this.value ? [this.value] : []
      },
      set (val) {
        if (this.multiple) {
          this.$emit('input', val)
          return
        }

        const vlen = val.length
        if (vlen) {
          this.$emit('input', val[vlen - 1])
        } else {
          this.$emit('input', null)
        }
      },
    },
    isEmpty () {
      if (Array.isArray(this.value)) {
        return this.value.length == 0
      }
      return !this.value
    },
    isEmptyList () {
      return this.tags.length == 0
    },
    isLabelActive () {
      return this.isFocused || !this.isEmpty
    },
    menuCanShow () {
      return this.tags.length > 0 || this.isCreatable
    },
    isDirty () {
      return this.inputValue.length > 0
    },
    hasExactMatch () {
      return this.tags.filter(tag => tag == this.stripInputValue).length > 0
    },
    isCreatable () {
      return this.isDirty && !this.hasExactMatch
    },
    stripInputValue () {
      return this.inputValue.replace(/^\/+/, '')
    }
  },
  methods: {
    focus () {
      this.isFocused = true
      this.isMenuActive = true
      this.$refs.input.focus()
      this.$emit('focus')
    },
    blur () {
      this.isFocused = false
      this.isMenuActive = false
      this.inputValue = ''
      this.$refs.input.blur()
    },
    async fetchTags () {
      try {
        this.loading = true

        const query = { q: this.inputValue }
        const res = await this.client.listTags(query)

        this.tags = res.data.map(t => t.name)
      } catch (e) {
        //
      } finally {
        this.loading = false
      }
    },
    debounceFetchData: debounce(async function () {
      await this.fetchTags()
    }, 200),
    selectNextTag () {
      const tag = this.tags[this.selectIndex + 1]

      if (!tag) {
        if (!this.tags.length) {
          return
        }

        this.selectIndex = -1
        this.selectNextTag()
        return
      }

      this.selectIndex++
    },
    selectPrevTag () {
      const tag = this.tags[this.selectIndex - 1]

      if (!tag) {
        if (!this.tags.length) {
          return
        }

        this.selectIndex = this.tags.length
        this.selectPrevTag()
        return
      }

      this.selectIndex--
    },
    includes (tag) {
      return this.selectedTags.includes(tag)
    },
    add (tag) {
      this.selectedTags = [ ...this.selectedTags, tag ]
    },
    remove (tag) {
      const index = this.selectedTags.findIndex(selected => tag == selected)
      this.selectedTags = [
        ...this.selectedTags.slice(0, index),
        ...this.selectedTags.slice(index + 1),
      ]
    },
    toggle (tag) {
      if (this.includes(tag)) {
        this.remove(tag)
      } else {
        this.add(tag)
      }
    },
    onFocus () {
      this.focus()
    },
    onClickOutside (e) {
      // Stop if not focused.
      if (!this.isFocused) {
        return
      }
      // Stop if inside $el.
      if (this.$el && this.$el.contains(e.target)) {
        return
      }
      // Stop if inside menu.
      if (this.$refs.menu && this.$refs.menu.$refs.content && this.$refs.menu.$refs.content.contains(e.target)) {
        return
      }
      this.blur()
    },
    onItemClick (tag) {
      this.toggle(tag)
      this.$refs.input.focus()
      this.updateMenuDimensions()
    },
    onUp () {
      this.selectPrevTag()
      this.updateMenuScrollPosition()
    },
    onDown () {
      this.selectNextTag()
      this.updateMenuScrollPosition()
    },
    onEnter () {
      // Toggle tag.
      if (!this.isEmptyList && !this.isCreatable) {
        let tag = this.tags[this.selectIndex]
        if (!tag) {
          tag = this.tags[0]
        }

        this.toggle(tag)
        this.inputValue = ''

      // Create new tag.
      } else {
        let val = this.inputValue.replace(/^\/*/, '')

        this.toggle(val)
        this.inputValue = ''
      }
    },
    onBackspace (e) {
      if (this.isEmpty) {
        // noop
      } else if (this.inputValue) {
        // noop
      } else {
        e.preventDefault()
        e.stopPropagation()

        const len = this.selectedTags.length
        const last = this.selectedTags[len - 1]
        this.remove(last)
        this.inputValue = `/${last}`
      }
    },
    onTab () {
      if (this.isEmptyList) {
        // noop
      } else {
        let tag = this.tags[this.selectIndex]
        if (!tag) {
          tag = this.tags[0]
        }

        this.inputValue = `/${tag}`
      }
    },
    getTagItemClasses (tag, n) {
      const selected = this.selectedTags.includes(tag)

      return {
        'v-list-item--link': true,
        'v-list-item--active': selected,
        'v-list-item--highlighted': n == this.selectIndex,
        'primary--text': selected,
      }
    },
    updateMenuScrollPosition () {
      this.$nextTick(() => {
        if (!this.$refs.menu) {
          return
        }

        const $el = this.$refs.menu.$refs.content
        const activeTile = $el.querySelector('.v-list-item--highlighted')
        const maxScrollTop = $el.scrollHeight - $el.offsetHeight

        this.$refs.menu.$refs.content.scrollTop = activeTile
          ? Math.min(maxScrollTop, Math.max(0, activeTile.offsetTop - $el.offsetHeight / 2 + activeTile.offsetHeight / 2))
          : $el.scrollTop
      })
    },
    updateMenuDimensions () {
      this.$refs.menu && this.$refs.menu.updateDimensions()
    },
  },
  watch: {
    inputValue () {
      if (!this.loading) {
        this.debounceFetchData()
      }
      this.selectIndex = -1
    },
    tags () {
      this.updateMenuDimensions()
    },
    isFocused (val) {
      if (!val) {
        return
      }
      if (this.isEmptyList) {
        this.fetchTags()
      }
    },
  },
}
</script>

<style lang="scss" module>
.input {
  font-size: 14px;
}
</style>
