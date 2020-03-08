<template>
  <div :class="containerClass">
    <Backdrop v-if="!disabled" absolute @click="focus" />
    <div :class="$style.list">
      <Tag
        v-for="tag in selectedTags"
        :key="tag.id"
        :tag="tag"
        :class="$style.tag"
      />
    </div>
    <Autocomplete
      v-if="!disabled"
      :class="$style.autocomplete"
      :options="tags"
      value-by="id"
      label-by="name"
      multiple
      :show="!!query"
      v-model="selectedTags"
      @input="resetQuery"
      :creatable="creatable"
      @create="handleCreate"
    >
      <template #input="slotProps">
        <Input
          v-model="query"
          noStyle
          ref="input"
          :class="$style.input"
          @keydown.up.prevent="slotProps.prev"
          @keydown.down.prevent="slotProps.next"
          @keydown.enter.prevent="slotProps.toggle"
          @keydown.backspace="handleBackspace"
          :placeholder="placeholder"
          @focus="focused = true"
          @blur="focused = false"
        />
      </template>
      <template #create="slotProps">
        <div :class="slotProps.class" @mouseover="slotProps.move" @mousedown="slotProps.toggle">
          Create "{{ query }}"
        </div>
      </template>
    </Autocomplete>
  </div>
</template>

<script>
import formMixin from '@/mixins/form'
import Tag from './Tag.vue'

export default {
  mixins: [formMixin],
  components: {
    Tag,
  },
  props: {
    value: {
      type: Array,
      default: null,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    noStyle: {
      type: Boolean,
      default: false,
    },
    placeholder: {
      type: String,
    },
  },
  data () {
    return {
      input: '',
      focused: false,
    }
  },
  computed: {
    containerClass () {
      return [this.$style.container, {
        [this.$style.container_border]: !this.noStyle,
        [this.$style.container_borderFocus]: this.focused,
      }]
    },
    query: {
      get () {
        return this.input
      },
      set (query) {
        this.input = query
      },
    },
    selectedTags: {
      get () {
        return this.value || []
      },
      set (tags) {
        this.$emit('input', tags)
      },
    },
    creatable () {
      return !this.hasExactName
    },
    tags () {
      const tags = this.$store.getters['tag/tags']
      return this.$store.getters['tag/search'](tags, this.query)
    },
    hasExactName () {
      const tags = this.$store.getters['tag/tags']
      return !!this.$store.getters['tag/findByName'](tags, this.query)
    },
  },
  methods: {
    resetQuery () {
      this.query = ''
    },
    async handleCreate () {
      const tag = await this.$store.dispatch('tag/create', { name: this.query })
      this.selectedTags = [ ...this.selectedTags, tag ]
      this.resetQuery()
    },
    handleBackspace () {
      if (this.query != '') {
        return true
      }
      this.selectedTags = this.selectedTags.slice(0, this.selectedTags.length - 1)
      return false
    },
    focus () {
      this.$refs.input.focus()
    },
    blur () {
      this.$refs.input.blur()
    },
  },
}
</script>

<style lang="scss" module>
.container {
  @apply flex;
  @apply flex-wrap;
  @apply -mx-1;
  @apply relative;
}

.container_border {
  @apply border-solid;
  @apply border;
  @apply border-control;
  @apply mx-0;
  @apply px-2;
  @apply py-1;
  min-height: theme('height.input');
  @apply rounded;
  transition: all .3s;

  &:hover {
    @apply border-primary;
  }
}

.container_borderFocus {
  @apply border-primary;
}

.list {
  @apply relative;
}

.tag {
  @apply m-1;
}

.autocomplete {
}

.input {
  @apply m-1;
  line-height: 18px;
}
</style>
