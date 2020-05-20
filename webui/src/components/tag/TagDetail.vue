<template>
  <div :class="$style.container">
    <template v-if="loading">
      <Placeholder>
        <div :class="[$style.phRow, $style.phLine]" />
        <div :class="[$style.phRow, $style.phLine, $style.phLineShort]" />
      </Placeholder>
    </template>

    <template v-else-if="editable">
      <InputGroup>
        <Label>Name</Label>
        <Input ref="name" v-model="model.name" :error="$v.model.name.$error" />
        <template #errors v-if="$v.model.name.$error">
          <p v-if="!$v.model.name.required">Name cannot be empty.</p>
        </template>
      </InputGroup>
      <InputGroup>
        <Label>Parent</Label>
        <TagInput ref="parent" v-model="parent" />
      </InputGroup>
    </template>

    <template v-else>
      <div :class="$style.name">
        <Tag :tag="model" :class="$style.tag" lg />
        <div :class="$style.breadcrumbs" v-if="breadcrumbs">
          <TagList :tags="parents" />
        </div>
      </div>
      <div :class="$style.additionInfo">
        <Anchor :class="$style.info" :to="{ name: 'tag.children', params: {parentName} }" underline v-if="parents">
          <IconLabel name="tagMultiple">{{ childrenCount }}</IconLabel>
        </Anchor>
        <Anchor :class="$style.info" :to="{ name: 'bookmark.list', query: bookmarkQuery }" underline>
          <IconLabel name="bookmark">{{ bookmarksCount }}</IconLabel>
        </Anchor>
      </div>
    </template>

    <slot
      name="controls"
      :submit="handleSubmit"
      :cancel="handleCancel"
    />
  </div>
</template>

<script>
import formMixin from '@/mixins/form'
import modelMixin from '@/mixins/model'
import placeholderMixin from '@/mixins/placeholder'
import keybindingMixin from '@/mixins/keybinding'
import IconLabel from '@/components/icon/IconLabel.vue'
import Tag from './Tag.vue'
import TagInput from './TagInput.vue'
import TagList from './TagList.vue'
import { validationMixin } from 'vuelidate'
import { required } from 'vuelidate/lib/validators'

export default {
  mixins: [formMixin, modelMixin({ prop: 'tag' }), placeholderMixin, validationMixin, keybindingMixin()],
  components: {
    IconLabel,
    Tag,
    TagInput,
    TagList,
  },
  props: {
    tag: {
      type: Object,
      default: null,
    },
    editable: {
      type: Boolean,
      default: false,
    },
    loading: {
      type: Boolean,
      default: false,
    },
  },
  data () {
    return {
      parents: null,
    }
  },
  created () {
    this.init()
  },
  mounted () {
    this.editable && this.autoFocusName()
  },
  computed: {
    model () {
      const gt = this.getter
      const st = this.setter
      return {
        $data: this.tag,
        get name () { return gt('name') },
        set name (v) { st('name', v) },
        get parentId () { return gt('parentId') },
        set parentId (v) { st('parentId', v) },
        get sort () { return gt('sort') },
        set sort (v) { st('sort', v) },
      }
    },
    name: {
      get () {
        return '#' + this.model.name
      },
      set (name) {
        this.model.name = name.slice(1)
      },
    },
    parent: {
      get () {
        const parentId = this.model.parentId
        if (!parentId) {
          return null
        }
        const tags = this.$store.getters['tag/tags']
        return [this.$store.getters['tag/find'](tags, parentId)]
      },
      set (value) {
        if (!value.length) {
          this.model.parentId = ''
        } else {
          const parent = value[value.length - 1]
          this.model.parentId = parent.id
        }
      },
    },
    breadcrumbs () {
      if (!this.parents) {
        return null
      }
      return '/' + this.parents.map(p => p.name).join('/')
    },
    children () {
      const tags = this.$store.getters['tag/tags']
      return this.$store.getters['tag/getByParent'](tags, this.tag.id)
    },
    childrenCount () {
      return this.children.length
    },
    bookmarks () {
      const pinls = this.$store.getters['pinl/pinls']
      return this.$store.getters['pinl/getByTag'](pinls, [this.tag.name])
    },
    bookmarksCount () {
      return this.bookmarks.length
    },
    bookmarkQuery () {
      const search = this.$store.getters['pinl/composeSearch']({ input: '', tags: [this.tag] })
      return { q: search }
    },
    parentName () {
      if (!this.parents) {
        return null
      }
      return [ ...this.parents.map(p => p.name), this.tag.name ]
    },
  },
  methods: {
    async init () {
      if (!this.tag) {
        return
      }
      this.parents = await this.$store.dispatch('tag/getParents', this.tag)
    },
    handleSubmit () {
      this.$v.$touch()
      if (this.$v.$error) {
        return
      }
      this.syncModel()
    },
    handleCancel () {
      this.$v.$reset()
      this.revertModel()
      this.$emit('update:editable', false)
      this.$emit('cancel')
    },
    autoFocusName () {
      this.$refs.name.focus()
    },
    handleKeyPress (e) {
      if (this.editable) {
        return
      }

      if (e.key == 'e') {
        this.$emit('update:editable', true)
        return
      }
    },
  },
  validations () {
    if (!this.editable) {
      return {}
    }
    return {
      model: {
        name: { required },
      },
    }
  },
  watch: {
    'tag' () {
      this.init()
    },
    'editable' (val) {
      if (val) {
        this.$nextTick(() => {
          this.autoFocusName()
        })
      }
    },
  },
}
</script>

<style lang="scss" module>
.name,
.parent {
  @apply mb-2;
}

.phRow {
  @apply mb-2;
}

.phLine {
  height: 1.2rem;
}

.phLineShort {
  width: 30%;
}

.tag {
  @apply font-bold;
}

.breadcrumbs {
  @apply inline-block;
  @apply ml-4;
}

.additionInfo {
  @apply flex;
  @apply flex-wrap;
  @apply -mx-1;

}

.info {
  @apply px-2;
}
</style>
