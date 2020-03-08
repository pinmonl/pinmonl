<template>
  <div :class="$style.container">
    <slot name="before"></slot>
    <div :class="$style.tagContainer">
      <Tag :tag="tag" :class="$style.tag" />
    </div>
    <Anchor :to="{ name: 'tag.children', params: {parentName} }" underline :class="$style.relContainer">
      <IconLabel name="tagMultiple">
        {{ childrenCount }}
      </IconLabel>
    </Anchor>
    <Anchor :to="{ name: 'bookmark', query: bookmarkQuery }" underline :class="$style.relContainer">
      <IconLabel name="bookmark">
        {{ bookmarksCount }}
      </IconLabel>
    </Anchor>
  </div>
</template>

<script>
import IconLabel from '@/components/icon/IconLabel.vue'
import Tag from './Tag.vue'

export default {
  components: {
    IconLabel,
    Tag,
  },
  props: {
    tag: {
      type: Object,
      required: true,
    },
    previsouParentName: {
      type: Array,
      default: () => ([]),
    },
  },
  computed: {
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
      return [ ...this.previsouParentName, this.tag.name ]
    },
  },
}
</script>

<style lang="scss" module>
.container {
  @apply flex;
  @apply text-sm;
  @apply py-2;
  @apply px-4;
  @apply items-center;
  @apply relative;

  &:hover {
    @apply bg-bg;
  }
}

.tagContainer {
  @apply flex-grow;
}

.tag {
  @apply text-sm;
}

.relContainer {
  @apply px-2;
  @apply whitespace-no-wrap;
  @apply inline-block;
  @apply relative;
  @apply z-10;
}
</style>
