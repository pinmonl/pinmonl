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
        <Input ref="name" v-model="model.name" />
      </InputGroup>
      <InputGroup>
        <Label>Parent</Label>
        <TagInput ref="parent" v-model="parent" />
      </InputGroup>
    </template>

    <template v-else>
      <div :class="$style.name" v-text="model.name" />
      <TagInput :class="$style.parent" :value="parent" noStyle disabled />
    </template>
  </div>
</template>

<script>
import formMixin from '@/mixins/form'
import modelMixin from '@/mixins/model'
import placeholderMixin from '@/mixins/placeholder'
import TagInput from './TagInput.vue'

export default {
  mixins: [formMixin, modelMixin({ prop: 'tag' }), placeholderMixin],
  components: {
    TagInput,
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
  },
}
</script>

<style lang="scss" module>
.container {
  @apply p-4;
}

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
</style>
