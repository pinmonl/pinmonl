<template>
  <div>
    <tag-chip
      v-for="tag in tags"
      :key="tag"
      :tag="tag"
    ></tag-chip>
  </div>
</template>

<script>
import TagChip from './TagChip'

export default {
  name: 'tag-chips',
  components: {
    TagChip,
  },
  props: ['tags'],
}
</script>

<style lang="scss" module>
</style>
