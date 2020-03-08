<template>
  <div :class="wrapperClass">
    <Backdrop
      v-if="isFocusing"
      @click="handleBlur"
    />
    <Autocomplete
      :options="tags"
      :show="isTagMode"
      value-by="id"
      label-by="name"
      multiple
      v-model="selectedTags"
      @input="resetQuery"
      ref="autocomplete"
    >
      <template #input="slotProps">
        <div :class="[$style.icon, {[$style.tagSign_active]: isTagMode}]">
          <Icon v-if="!isTagMode" name="magnify" />
          <Icon v-if="isTagMode" name="pound" />
        </div>
        <div :class="$style.inputWrapper" @click="focus">
          <div v-for="tag in selectedTags" :key="tag.id" :class="$style.tagContainer">
            <Tag :tag="tag" />
          </div>
          <div :class="$style.inputContainer">
            <Input
              :class="$style.input"
              ref="input"
              noStyle
              @focus="handleFocus"
              v-model="query"
              @keydown.down.prevent="slotProps.next"
              @keydown.up.prevent="slotProps.prev"
              @keydown.enter.prevent="handleEnter"
              @keydown.backspace="handleBackspace"
            />
          </div>
        </div>
      </template>
    </Autocomplete>
  </div>
</template>

<script>
import Autocomplete from '@/components/form/Autocomplete.vue'
import Tag from '@/components/tag/Tag.vue'
import Input from '@/components/form/Input.vue'

export default {
  components: {
    Autocomplete,
    Tag,
    Input,
  },
  model: {
    prop: 'search',
    event: 'input',
  },
  props: {
    search: {
      type: String,
      default: '',
    },
    tagSrc: {
      type: Array,
      default: null,
    },
  },
  data () {
    return {
      input: '',
      isFocusing: false,
      selectedTags: [],
    }
  },
  created () {
    this.initSearch()
  },
  computed: {
    query: {
      get () {
        return this.isTagMode
          ? this.input.replace(/^.*?#/, '')
          : this.input
      },
      set (newValue) {
        let query = newValue
        if (this.isTagMode) {
          query = this.input.replace(/#.*$/, '#') + newValue
        }
        query = query.replace(/#(.*)#$/, '')
        this.input = query
      },
    },
    wrapperClass () {
      return [this.$style.wrapper, {
        [this.$style.wrapper_focus]: this.isFocusing,
      }]
    },
    isTagMode () {
      return this.input.includes('#') && this.isFocusing
    },
    tags () {
      if (!this.isTagMode) {
        return []
      }
      const tags = (this.tagSrc == null) ? this.$store.getters['tag/tags'] : this.tagSrc
      const query = this.query.replace(/^.*#/, '')
      return this.$store.getters['tag/search'](tags, query)
    },
  },
  methods: {
    focus () {
      this.$refs.input.focus()
    },
    blur () {
      this.$refs.input.blur()
    },
    handleFocus () {
      this.isFocusing = true
    },
    handleBlur () {
      this.isFocusing = false
      this.input = this.input.replace(/#.*$/, '')
    },
    resetQuery () {
      if (this.isTagMode) {
        this.input = this.input.replace(/#.*$/, '#')
      }
      this.focus()
    },
    handleBackspace () {
      if (this.query != '') {
        return true
      }
      this.selectedTags = this.selectedTags.slice(0, this.selectedTags.length - 1)
      return false
    },
    handleEnter () {
      if (this.isTagMode) {
        this.$refs.autocomplete.autoToggle()
        return true
      }

      const composed = this.composeSearch()
      this.storedSearch = composed
      this.$emit('input', composed)
      return false
    },
    composeSearch () {
      const search = { input: this.query, tags: this.selectedTags }
      const composed = this.$store.getters['pinl/composeSearch'](search)
      return composed
    },
    parseSearch () {
      const parsed = this.$store.getters['pinl/parseSearch'](this.search)
      return parsed
    },
    initSearch () {
      const { input, tags } = this.parseSearch()
      this.input = input
      this.selectedTags = tags
    },
  },
  watch: {
    search (newValue, oldValue) {
      if (newValue != oldValue) {
        this.initSearch()
      }
    },
    tags (tags) {
      if (tags.length) {
        this.$refs.autocomplete.move(0)
      }
    },
  },
}
</script>

<style lang="scss" module>
.wrapper {
  @apply relative;
  @apply border;
  @apply border-solid;
  @apply border-control;
  @apply rounded;
  min-width: 300px;
  @apply bg-clear;

  @screen sm-down {
    @apply min-w-0;
  }
}

.wrapper_focus {
  @apply border-primary;
}

.icon {
  @apply absolute;
  margin-left: 0.5rem;
  @apply text-text-secondary;
  @apply top-0;
  @apply left-0;
  @apply h-input;
  @apply flex;
  @apply items-center;
}

.tagSign_active {
  @apply text-primary;
}

.inputWrapper {
  @apply pl-10;
  @apply pr-2;
  @apply h-full;
  @apply w-full;
  @apply relative;
  @apply -mr-1;
  @apply flex;
  @apply items-center;
}

.tagContainer {
  @apply inline-block;
  @apply mr-1;
}

.inputContainer {
  @apply flex-1;
}

.input {
  @apply h-input;
  @apply inline-block;
  line-height: theme('height.input');
}
</style>
